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

// ExtractOuterArchive handles validation, triggers extraction, and detects the inner firmware type.
func (c *FirmwareController) ExtractOuterArchive(filePath, outputDir string) {
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
		if ext == ".md5" && strings.HasSuffix(strings.ToLower(base), ".tar.md5") {
			folderName = "extracted_" + strings.TrimSuffix(base, ".tar.md5")
		} else if ext == ".gz" && strings.HasSuffix(strings.ToLower(base), ".tar.gz") {
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
	} else if ext == ".tar" || ext == ".md5" || strings.HasSuffix(strings.ToLower(filePath), ".tar.md5") {
		extractErr = c.model.ExtractTarRaw(filePath, outputDir, onProgress)
	} else {
		c.view.RenderError(fmt.Errorf("unsupported file format: %s. Supported formats: .zip, .tgz, .tar.gz, .tar, .tar.md5", ext))
		return
	}

	fmt.Println() // Clear the progress line
	if extractErr != nil {
		c.view.RenderError(extractErr)
	} else {
		c.view.RenderSuccess(fmt.Sprintf("Successfully extracted firmware to: %s", outputDir))
		
		// Run content validation / auto-detection
		detectedType := c.model.DetectFirmwareType(outputDir)
		fmt.Printf("\033[36m=== FIRMWARE VALIDATION / DETECTION ===\033[0m\n")
		fmt.Printf("  Detected Content: \033[32m%s\033[0m\n\n", detectedType)
	}
}

// ExtractSamsungInner handles selective component extraction for Samsung firmware from a folder containing .tar.md5 files.
func (c *FirmwareController) ExtractSamsungInner() {
	folderPath := c.view.PromptInput("Masukkan Path Folder berisi file Samsung (.tar.md5): ")
	if folderPath == "" {
		c.view.RenderError(fmt.Errorf("folder path cannot be empty"))
		return
	}

	// Check if directory exists
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) || !info.IsDir() {
		c.view.RenderError(fmt.Errorf("directory does not exist or is not a folder: %s", folderPath))
		return
	}

	// Find Samsung files
	samsungFiles := c.model.FindSamsungFiles(folderPath)
	if len(samsungFiles) == 0 {
		c.view.RenderError(fmt.Errorf("no Samsung firmware files (AP, BL, CP, CSC, HOME_CSC) found in folder: %s", folderPath))
		return
	}

	// Print detected files
	fmt.Printf("\n\033[36m=== DETECTED SAMSUNG COMPONENTS ===\033[0m\n")
	for key, path := range samsungFiles {
		fmt.Printf("  - %s: %s\n", key, filepath.Base(path))
	}
	fmt.Println()

	// Ask for selection
	choice := c.view.PromptInput("Pilih file yang ingin diekstrak (pisah koma, misal: AP,BL) atau tekan Enter untuk semua (Auto): ")
	selectedKeys := []string{}
	if strings.TrimSpace(choice) == "" {
		// Auto mode: select all detected files
		for key := range samsungFiles {
			selectedKeys = append(selectedKeys, key)
		}
	} else {
		// Manual mode: parse choices
		parts := strings.Split(choice, ",")
		for _, part := range parts {
			key := strings.TrimSpace(strings.ToUpper(part))
			if _, exists := samsungFiles[key]; exists {
				selectedKeys = append(selectedKeys, key)
			} else {
				fmt.Printf("\033[33mWarning: Component '%s' not found or invalid. Skipping.\033[0m\n", key)
			}
		}
	}

	if len(selectedKeys) == 0 {
		c.view.RenderError(fmt.Errorf("no valid components selected for extraction"))
		return
	}

	// Ask for output directory
	outputDir := c.view.PromptInput("Masukkan Folder Output (kosongkan untuk default): ")
	if outputDir == "" {
		outputDir = filepath.Join(folderPath, "extracted_samsung")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		c.view.RenderError(err)
		return
	}

	onProgress := func(fileName string) {
		if len(fileName) > 60 {
			fileName = "..." + fileName[len(fileName)-57:]
		}
		fmt.Printf("\r\033[KExtracting: %s", fileName)
	}

	// Extract selected files
	for _, key := range selectedKeys {
		filePath := samsungFiles[key]
		// Create specific sub-folder for this component (e.g. outputDir/AP)
		compOutputDir := filepath.Join(outputDir, key)
		fmt.Printf("\nExtracting %s component (%s) to %s...\n", key, filepath.Base(filePath), compOutputDir)
		
		err := c.model.ExtractTarRaw(filePath, compOutputDir, onProgress)
		fmt.Println() // Clear progress line
		if err != nil {
			fmt.Printf("\033[31mError extracting %s: %v\033[0m\n", key, err)
		} else {
			fmt.Printf("Finished processing %s component.\n", key)
		}
	}

	c.view.RenderSuccess(fmt.Sprintf("Samsung firmware extraction completed in: %s", outputDir))
}
