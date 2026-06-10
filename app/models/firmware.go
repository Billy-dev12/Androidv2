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
