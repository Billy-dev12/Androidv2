package models

import (
	"fmt"
	"strings"
)

// FileTransferModel handles file operations between local system and Android device.
type FileTransferModel struct {
	executor *ADBExecutor
}

// NewFileTransferModel creates a new FileTransferModel.
func NewFileTransferModel(executor *ADBExecutor) *FileTransferModel {
	return &FileTransferModel{executor: executor}
}

// Push uploads a local file or directory to a remote path on the device.
func (m *FileTransferModel) Push(localPath string, remotePath string, deviceID string) (string, error) {
	args := []string{}
	if deviceID != "" {
		args = append(args, "-s", deviceID)
	}
	args = append(args, "push", localPath, remotePath)

	stdout, stderr, err := m.executor.Execute(args...)
	if err != nil {
		return "", fmt.Errorf("failed to push file: %w (stderr: %s)", err, strings.TrimSpace(stderr))
	}
	return stdout, nil
}

// Pull downloads a remote file or directory from the device to a local path.
func (m *FileTransferModel) Pull(remotePath string, localPath string, deviceID string) (string, error) {
	args := []string{}
	if deviceID != "" {
		args = append(args, "-s", deviceID)
	}
	args = append(args, "pull", remotePath, localPath)

	stdout, stderr, err := m.executor.Execute(args...)
	if err != nil {
		return "", fmt.Errorf("failed to pull file: %w (stderr: %s)", err, strings.TrimSpace(stderr))
	}
	return stdout, nil
}
