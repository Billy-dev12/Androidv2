package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"android-tool-mvc/app/models"
	"android-tool-mvc/resources/views"
)

// FirmwareController coordinates the UI prompts and extraction operations.
type FirmwareController struct {
	model *models.FirmwareExtractor
	view  *views.ConsoleView
}

// NewFirmwareController creates a new instance of FirmwareController.
func NewFirmwareController(model *models.FirmwareExtractor, view *views.ConsoleView) *FirmwareController {
	return &FirmwareController{
		model: model,
		view:  view,
	}
}

// ExtractXiaomi handles validation and triggers extraction for Xiaomi firmware.
func (c *FirmwareController) ExtractXiaomi(filePath, outputDir string) {
	if filePath == "" {
		c.view.RenderError(fmt.Errorf("file path cannot be empty"))
		return
	}

	// Resolve output directory if left empty
	if outputDir == "" {
		dir := filepath.Dir(filePath)
		base := filepath.Base(filePath)
		ext := filepath.Ext(filePath)
		var folderName string
		if ext == ".gz" && strings.HasSuffix(strings.ToLower(base), ".tar.gz") {
			folderName = "extracted_" + strings.TrimSuffix(base, ".tar.gz")
		} else {
			folderName = "extracted_" + strings.TrimSuffix(base, ext)
		}
		outputDir = filepath.Join(dir, folderName)
	}

	// Check if file exists
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.view.RenderError(fmt.Errorf("file does not exist: %s", filePath))
		return
	}
	if info.IsDir() {
		c.view.RenderError(fmt.Errorf("path is a directory: %s", filePath))
		return
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	fmt.Printf("\nExtracting %s to %s...\n", filePath, outputDir)

	onProgress := func(fileName string) {
		// Truncate long filenames to look neat
		if len(fileName) > 60 {
			fileName = "..." + fileName[len(fileName)-57:]
		}
		fmt.Printf("\r\033[KExtracting: %s", fileName)
	}

	var extractErr error
	if ext == ".zip" {
		extractErr = c.model.ExtractZip(filePath, outputDir, onProgress)
	} else if ext == ".tgz" || (ext == ".gz" && strings.HasSuffix(strings.ToLower(filePath), ".tar.gz")) {
		extractErr = c.model.ExtractTarGz(filePath, outputDir, onProgress)
	} else if ext == ".tar" {
		extractErr = c.model.ExtractTarRaw(filePath, outputDir, onProgress)
	} else {
		c.view.RenderError(fmt.Errorf("unsupported file format: %s. Supported formats: .zip, .tgz, .tar.gz, .tar", ext))
		return
	}

	fmt.Println() // Clear the progress line
	if extractErr != nil {
		c.view.RenderError(extractErr)
	} else {
		c.view.RenderSuccess(fmt.Sprintf("Successfully extracted firmware to: %s", outputDir))
	}
}
