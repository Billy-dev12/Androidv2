package models

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FirmwareExtractor handles the archive decompression logic.
type FirmwareExtractor struct{}

// NewFirmwareExtractor creates a new instance of FirmwareExtractor.
func NewFirmwareExtractor() *FirmwareExtractor {
	return &FirmwareExtractor{}
}

// ExtractZip extracts a ZIP archive to the target destination.
func (fe *FirmwareExtractor) ExtractZip(src, dest string, onProgress func(string)) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		if onProgress != nil {
			onProgress(f.Name)
		}

		// Prevent Zip Slip vulnerability
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// ExtractTarGz extracts a TGZ/TAR.GZ archive to the target destination.
func (fe *FirmwareExtractor) ExtractTarGz(src, dest string, onProgress func(string)) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	return fe.ExtractTar(gzr, dest, onProgress)
}

// ExtractTarRaw extracts a raw TAR archive to the target destination.
func (fe *FirmwareExtractor) ExtractTarRaw(src, dest string, onProgress func(string)) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	return fe.ExtractTar(f, dest, onProgress)
}

// ExtractTar reads from an io.Reader and extracts tar headers.
func (fe *FirmwareExtractor) ExtractTar(r io.Reader, dest string, onProgress func(string)) error {
	tarReader := tar.NewReader(r)

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if onProgress != nil {
			onProgress(header.Name)
		}

		// Prevent Directory Traversal vulnerability
		fpath := filepath.Join(dest, header.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in tar: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fpath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}

// DetectFirmwareType scans the output directory and identifies the firmware brand/type.
func (fe *FirmwareExtractor) DetectFirmwareType(dir string) string {
	var detected []string

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		name := strings.ToLower(info.Name())
		if name == "payload.bin" {
			detected = append(detected, "Google Pixel / OnePlus / Xiaomi (payload.bin)")
		} else if name == "update.app" {
			detected = append(detected, "Huawei / Honor (UPDATE.APP)")
		} else if strings.HasSuffix(name, ".ozip") {
			detected = append(detected, "Oppo / Realme (.ozip)")
		} else if name == "update.pkg" {
			detected = append(detected, "Vivo / iQOO (update.pkg)")
		} else if strings.Contains(name, "scatter") && strings.HasSuffix(name, ".txt") {
			detected = append(detected, fmt.Sprintf("MediaTek Scatter (%s)", info.Name()))
		} else if strings.HasSuffix(name, ".tar.md5") || strings.HasSuffix(name, ".lz4") {
			detected = append(detected, fmt.Sprintf("Samsung Firmware (%s)", info.Name()))
		} else if name == "flashfile.xml" || name == "servicefile.xml" {
			detected = append(detected, fmt.Sprintf("Motorola Firmware (%s)", info.Name()))
		}
		return nil
	})

	if len(detected) == 0 {
		return "Generic / Unknown Firmware Structure"
	}

	// Remove duplicates if any, and join
	unique := make(map[string]bool)
	var result []string
	for _, d := range detected {
		if !unique[d] {
			unique[d] = true
			result = append(result, d)
		}
	}
	return strings.Join(result, ", ")
}
