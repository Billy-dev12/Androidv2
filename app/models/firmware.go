package models

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
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

// FindSamsungFiles scans the given directory for Samsung AP, BL, CP, CSC, HOME_CSC files.
func (fe *FirmwareExtractor) FindSamsungFiles(dir string) map[string]string {
	files := make(map[string]string)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return files
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToUpper(entry.Name())
		path := filepath.Join(dir, entry.Name())

		// Check prefix and extension (.tar or .tar.md5)
		if strings.HasSuffix(name, ".TAR") || strings.HasSuffix(name, ".TAR.MD5") || strings.HasSuffix(name, ".MD5") {
			if strings.HasPrefix(name, "AP_") {
				files["AP"] = path
			} else if strings.HasPrefix(name, "BL_") {
				files["BL"] = path
			} else if strings.HasPrefix(name, "CP_") {
				files["CP"] = path
			} else if strings.HasPrefix(name, "HOME_CSC_") {
				files["HOME_CSC"] = path
			} else if strings.HasPrefix(name, "CSC_") {
				files["CSC"] = path
			}
		}
	}
	return files
}

// DecompressLZ4 decompresses an LZ4 file using the system lz4 tool.
func (fe *FirmwareExtractor) DecompressLZ4(src, dest string) error {
	cmd := exec.Command("lz4", "-d", "-f", src, dest)
	return cmd.Run()
}

// DecompressFolderLZ4 scans the directory for any .lz4 files, decompresses them, and deletes the originals.
func (fe *FirmwareExtractor) DecompressFolderLZ4(dir string, onProgress func(string)) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".lz4") {
			dest := strings.TrimSuffix(path, filepath.Ext(path))
			if onProgress != nil {
				onProgress(info.Name())
			}
			if err := fe.DecompressLZ4(path, dest); err != nil {
				return err
			}
			// Delete the original lz4 file to clean up
			os.Remove(path)
		}
		return nil
	})
}

// FindSamsungComponentsInZip scans the ZIP and maps component names to actual filenames in ZIP.
func (fe *FirmwareExtractor) FindSamsungComponentsInZip(zipPath string) (map[string]string, error) {
	components := make(map[string]string)
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		name := strings.ToUpper(filepath.Base(f.Name))
		if strings.HasSuffix(name, ".TAR") || strings.HasSuffix(name, ".TAR.MD5") || strings.HasSuffix(name, ".MD5") {
			if strings.HasPrefix(name, "AP_") {
				components["AP"] = f.Name
			} else if strings.HasPrefix(name, "BL_") {
				components["BL"] = f.Name
			} else if strings.HasPrefix(name, "CP_") {
				components["CP"] = f.Name
			} else if strings.HasPrefix(name, "HOME_CSC_") {
				components["HOME_CSC"] = f.Name
			} else if strings.HasPrefix(name, "CSC_") {
				components["CSC"] = f.Name
			}
		}
	}
	return components, nil
}

// PartitionInfo holds metadata about a single partition image file.
type PartitionInfo struct {
	Name      string
	Size      int64
	SizeHuman string
	FileType  string // "img" or "bin"
}

// ScanPartitions walks a directory and returns metadata of all .img and .bin files.
func (fe *FirmwareExtractor) ScanPartitions(dir string) ([]PartitionInfo, error) {
	var partitions []PartitionInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		name := strings.ToLower(info.Name())
		if !strings.HasSuffix(name, ".img") && !strings.HasSuffix(name, ".bin") {
			return nil
		}

		var sizeHuman string
		switch {
		case info.Size() >= 1<<30:
			sizeHuman = fmt.Sprintf("%.2f GB", float64(info.Size())/(1<<30))
		case info.Size() >= 1<<20:
			sizeHuman = fmt.Sprintf("%.2f MB", float64(info.Size())/(1<<20))
		case info.Size() >= 1<<10:
			sizeHuman = fmt.Sprintf("%.2f KB", float64(info.Size())/(1<<10))
		default:
			sizeHuman = fmt.Sprintf("%d B", info.Size())
		}

		ext := strings.TrimPrefix(filepath.Ext(name), ".")
		partitions = append(partitions, PartitionInfo{
			Name:      info.Name(),
			Size:      info.Size(),
			SizeHuman: sizeHuman,
			FileType:  ext,
		})
		return nil
	})

	return partitions, err
}

// ExtractSpecificFilesFromZip extracts only selected files from the ZIP archive.
func (fe *FirmwareExtractor) ExtractSpecificFilesFromZip(zipPath string, filesToExtract map[string]string, dest string, onProgress func(string)) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	// Create a fast lookup map
	lookup := make(map[string]string)
	for comp, zipName := range filesToExtract {
		lookup[zipName] = comp
	}

	for _, f := range r.File {
		_, ok := lookup[f.Name]
		if !ok {
			continue
		}

		if onProgress != nil {
			onProgress(f.Name)
		}

		// Prevent Zip Slip vulnerability
		fpath := filepath.Join(dest, filepath.Base(f.Name))
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
		fmt.Println() // Newline after extracting this component
	}
	return nil
}
