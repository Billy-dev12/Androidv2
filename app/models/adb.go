package models

import (
	"bytes"
	"os/exec"
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
