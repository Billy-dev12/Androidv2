package models

import (
	"fmt"
	"strings"
)

// ApplicationModel handles app-level package management.
type ApplicationModel struct {
	executor *ADBExecutor
}

// NewApplicationModel creates a new ApplicationModel.
func NewApplicationModel(executor *ADBExecutor) *ApplicationModel {
	return &ApplicationModel{executor: executor}
}

// Install installs an APK file on the target device.
func (m *ApplicationModel) Install(apkPath string, deviceID string) (string, error) {
	args := []string{}
	if deviceID != "" {
		args = append(args, "-s", deviceID)
	}
	args = append(args, "install", apkPath)

	stdout, stderr, err := m.executor.Execute(args...)
	if err != nil {
		return "", fmt.Errorf("install failed: %w (stderr: %s)", err, strings.TrimSpace(stderr))
	}
	return stdout, nil
}

// Uninstall uninstalls a package from the target device.
func (m *ApplicationModel) Uninstall(packageName string, deviceID string) (string, error) {
	args := []string{}
	if deviceID != "" {
		args = append(args, "-s", deviceID)
	}
	args = append(args, "uninstall", packageName)

	stdout, stderr, err := m.executor.Execute(args...)
	if err != nil {
		return "", fmt.Errorf("uninstall failed: %w (stderr: %s)", err, strings.TrimSpace(stderr))
	}
	return stdout, nil
}
