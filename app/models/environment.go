package models

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type EnvironmentModel struct{}

func NewEnvironmentModel() *EnvironmentModel {
	return &EnvironmentModel{}
}

type ToolInfo struct {
	Name    string
	Path    string
	Version string
	Status  string
}

func (m *EnvironmentModel) CheckAll() (map[string]ToolInfo, string) {
	toolNames := []string{"adb", "fastboot", "lz4", "java", "python3", "python", "unzip", "tar", "7z", "git"}

	results := make(map[string]ToolInfo)
	for _, name := range toolNames {
		results[name] = m.checkTool(name)
	}

	adbServer := m.checkADBServer()
	osInfo := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	var sb strings.Builder
	sb.WriteString("=== SYSTEM INFORMATION ===\n")
	sb.WriteString(fmt.Sprintf("  OS/Arch    : %s\n", osInfo))
	sb.WriteString(fmt.Sprintf("  ADB Server : %s\n", adbServer))

	return results, sb.String()
}

func (m *EnvironmentModel) checkADBServer() string {
	cmd := exec.Command("adb", "start-server")
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "Failed to start"
	}
	out := stderr.String()
	if strings.Contains(out, "daemon started successfully") || strings.Contains(out, "already running") || out == "" {
		return "Running"
	}
	return "Unknown"
}

func (m *EnvironmentModel) checkTool(name string) ToolInfo {
	path, err := exec.LookPath(name)
	if err != nil {
		return ToolInfo{Name: name, Status: "Not Found"}
	}

	version := m.getVersion(name)
	return ToolInfo{Name: name, Path: path, Version: version, Status: "OK"}
}

func (m *EnvironmentModel) getVersion(name string) string {
	var args []string
	switch name {
	case "java":
		args = []string{"-version"}
	case "python3", "python":
		args = []string{"--version"}
	case "7z":
		args = []string{}
	case "adb":
		args = []string{"--version"}
	case "fastboot":
		args = []string{"--version"}
	case "lz4":
		args = []string{"--version"}
	case "unzip":
		args = []string{"-v"}
	case "tar":
		args = []string{"--version"}
	case "git":
		args = []string{"--version"}
	default:
		args = []string{"--version"}
	}

	cmd := exec.Command(name, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "Unknown"
	}

	out := strings.TrimSpace(stdout.String() + " " + stderr.String())
	lines := strings.Split(out, "\n")
	if len(lines) > 0 {
		first := strings.TrimSpace(lines[0])
		if len(first) > 80 {
			first = first[:80] + "..."
		}
		return first
	}
	return "Unknown"
}

func (m *EnvironmentModel) ToolInfoToString(results map[string]ToolInfo) string {
	var sb strings.Builder
	sb.WriteString("=== INSTALLED TOOLS ===\n")
	for _, name := range []string{"adb", "fastboot", "lz4", "java", "python3", "python", "unzip", "tar", "7z", "git"} {
		info := results[name]
		if info.Status == "OK" {
			sb.WriteString(fmt.Sprintf("  %-10s ✓  %s\n", info.Name, info.Path))
			sb.WriteString(fmt.Sprintf("  %-10s     %s\n", "", info.Version))
		} else {
			sb.WriteString(fmt.Sprintf("  %-10s ✗  Not Found\n", name))
		}
	}
	return sb.String()
}
