package models

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ADBExecutor handles executing ADB commands.
type ADBExecutor struct{}

// NewADBExecutor creates a new ADBExecutor instance.
func NewADBExecutor() *ADBExecutor {
	return &ADBExecutor{}
}

// IsInstalled checks if adb command is available in system PATH.
func (a *ADBExecutor) IsInstalled() bool {
	_, err := exec.LookPath("adb")
	return err == nil
}

// ExecuteToFile runs an adb command and writes stdout directly to a file.
func (a *ADBExecutor) ExecuteToFile(outputPath string, args ...string) error {
	cmd := exec.Command("adb", args...)
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer f.Close()
	cmd.Stdout = f
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("adb command failed: %w (stderr: %s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// Execute runs an adb command with the specified arguments.
// It returns the stdout, stderr, and any execution error.
func (a *ADBExecutor) Execute(args ...string) (string, string, error) {
	cmd := exec.Command("adb", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
