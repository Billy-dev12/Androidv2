package models

import (
	"fmt"
	"strings"
)

type DiagnosticsModel struct {
	executor *ADBExecutor
}

func NewDiagnosticsModel(executor *ADBExecutor) *DiagnosticsModel {
	return &DiagnosticsModel{executor: executor}
}

func (m *DiagnosticsModel) GetMemory(deviceID string) (string, error) {
	args := m.buildArgs(deviceID, "shell", "dumpsys", "meminfo")
	stdout, stderr, err := m.executor.Execute(args...)
	if err != nil {
		return "", fmt.Errorf("memory info failed: %w (stderr: %s)", err, strings.TrimSpace(stderr))
	}
	var result strings.Builder
	result.WriteString("=== MEMORY INFO ===\n")
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		tr := strings.TrimSpace(line)
		if strings.HasPrefix(tr, "Total RAM:") || strings.HasPrefix(tr, "Free RAM:") ||
			strings.HasPrefix(tr, "Used RAM:") || strings.HasPrefix(tr, "Lost RAM:") ||
			strings.HasPrefix(tr, "Native Heap:") || strings.HasPrefix(tr, "Dalvik Heap:") ||
			strings.HasPrefix(tr, "Total RAM:") || strings.HasPrefix(tr, "Free RAM:") {
			result.WriteString("  " + tr + "\n")
		}
	}
	if result.Len() == 19 {
		for _, line := range lines {
			if strings.Contains(line, "RAM:") || strings.Contains(line, "Heap:") {
				result.WriteString("  " + strings.TrimSpace(line) + "\n")
			}
		}
	}
	return result.String(), nil
}

func (m *DiagnosticsModel) GetStorage(deviceID string) (string, error) {
	args := m.buildArgs(deviceID, "shell", "df")
	stdout, stderr, err := m.executor.Execute(args...)
	if err != nil {
		return "", fmt.Errorf("storage info failed: %w (stderr: %s)", err, strings.TrimSpace(stderr))
	}
	var result strings.Builder
	result.WriteString("=== STORAGE INFO ===\n")
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		tr := strings.TrimSpace(line)
		if tr == "" {
			continue
		}
		fields := strings.Fields(tr)
		if len(fields) < 6 {
			continue
		}
		mount := fields[len(fields)-1]
		if strings.HasPrefix(mount, "/data") || strings.HasPrefix(mount, "/sdcard") ||
			strings.HasPrefix(mount, "/storage") || mount == "/" ||
			strings.Contains(mount, "emulated") {
			result.WriteString(fmt.Sprintf("  %-30s %s used / %s total\n", mount, fields[2], fields[1]))
		}
	}
	if result.Len() == 20 {
		result.WriteString("  " + stdout)
	}
	return result.String(), nil
}

func (m *DiagnosticsModel) GetCPU(deviceID string) (string, error) {
	args := m.buildArgs(deviceID, "shell", "cat", "/proc/cpuinfo")
	stdout, stderr, err := m.executor.Execute(args...)
	if err != nil {
		return "", fmt.Errorf("cpu info failed: %w (stderr: %s)", err, strings.TrimSpace(stderr))
	}
	var result strings.Builder
	result.WriteString("=== CPU INFO ===\n")
	seen := make(map[string]bool)
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		tr := strings.TrimSpace(line)
		if tr == "" {
			continue
		}
		parts := strings.SplitN(tr, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if key == "processor" || key == "Processor" || key == "Hardware" ||
				key == "Features" || key == "CPU architecture" || key == "model name" ||
				key == "BogoMIPS" {
				uniq := key + val
				if !seen[uniq] {
					seen[uniq] = true
					result.WriteString(fmt.Sprintf("  %s: %s\n", key, val))
				}
			}
		}
	}
	if result.Len() == 14 {
		result.WriteString("  " + stdout)
	}
	return result.String(), nil
}

func (m *DiagnosticsModel) GetSensors(deviceID string) (string, error) {
	args := m.buildArgs(deviceID, "shell", "dumpsys", "sensorservice")
	stdout, stderr, err := m.executor.Execute(args...)
	if err != nil {
		return "", fmt.Errorf("sensor info failed: %w (stderr: %s)", err, strings.TrimSpace(stderr))
	}
	var result strings.Builder
	result.WriteString("=== SENSOR INFO ===\n")
	lines := strings.Split(stdout, "\n")
	inList := false
	for _, line := range lines {
		tr := strings.TrimSpace(line)
		if strings.Contains(tr, "Sensor List:") {
			inList = true
			continue
		}
		if inList {
			if tr == "" || strings.HasPrefix(line, "  ") == false {
				if len(result.String()) > 30 {
					break
				}
				continue
			}
			parts := strings.SplitN(tr, ":", 2)
			if len(parts) == 2 {
				result.WriteString(fmt.Sprintf("  %s\n", strings.TrimSpace(parts[0])))
			} else if len(parts) > 0 && parts[0] != "" {
				result.WriteString(fmt.Sprintf("  %s\n", parts[0]))
			}
		}
	}
	if result.Len() == 20 {
		result.WriteString("  (no sensor list found)\n")
	}
	return result.String(), nil
}

func (m *DiagnosticsModel) GetNetwork(deviceID string) (string, error) {
	var result strings.Builder
	result.WriteString("=== NETWORK & SIGNAL INFO ===\n")

	getpropArgs := m.buildArgs(deviceID, "shell", "getprop")
	stdout, _, _ := m.executor.Execute(getpropArgs...)
	props := map[string]string{
		"gsm.sim.operator.alpha":   "Operator",
		"gsm.sim.operator.numeric": "Operator (MCC/MNC)",
		"gsm.operator.alpha":       "Network Operator",
		"gsm.operator.iso-country": "Country",
		"gsm.operator.numeric":     "Network (MCC/MNC)",
	}
	propLines := strings.Split(stdout, "\n")
	for _, line := range propLines {
		for prop, label := range props {
			if strings.Contains(line, "["+prop+"]") {
				parts := strings.Split(line, "]")
				if len(parts) >= 2 {
					val := strings.Trim(parts[1], " []")
					result.WriteString(fmt.Sprintf("  %s: %s\n", label, val))
				}
			}
		}
	}

	telephonyArgs := m.buildArgs(deviceID, "shell", "dumpsys", "telephony")
	telOut, _, _ := m.executor.Execute(telephonyArgs...)
	for _, line := range strings.Split(telOut, "\n") {
		tr := strings.TrimSpace(line)
		if strings.Contains(tr, "Signal Strength") || strings.Contains(tr, "mSignalStrength") {
			result.WriteString(fmt.Sprintf("  %s\n", tr))
			break
		}
	}

	if result.Len() < 30 {
		deviceArgs := m.buildArgs(deviceID, "shell", "getprop", "ro.telephony.default_network")
		netType, _, _ := m.executor.Execute(deviceArgs...)
		if strings.TrimSpace(netType) != "" {
			result.WriteString(fmt.Sprintf("  Preferred Network: %s\n", strings.TrimSpace(netType)))
		}
	}

	return result.String(), nil
}

func (m *DiagnosticsModel) GetDisplay(deviceID string) (string, error) {
	var result strings.Builder
	result.WriteString("=== DISPLAY INFO ===\n")

	sizeArgs := m.buildArgs(deviceID, "shell", "wm", "size")
	sizeOut, _, _ := m.executor.Execute(sizeArgs...)
	for _, line := range strings.Split(sizeOut, "\n") {
		if strings.Contains(line, "Physical size") || strings.Contains(line, "Override size") {
			result.WriteString(fmt.Sprintf("  %s\n", strings.TrimSpace(line)))
		}
	}

	densityArgs := m.buildArgs(deviceID, "shell", "wm", "density")
	densOut, _, _ := m.executor.Execute(densityArgs...)
	for _, line := range strings.Split(densOut, "\n") {
		if strings.Contains(line, "Physical density") || strings.Contains(line, "Override density") {
			result.WriteString(fmt.Sprintf("  %s\n", strings.TrimSpace(line)))
		}
	}

	return result.String(), nil
}

func (m *DiagnosticsModel) GetAll(deviceID string) (string, error) {
	var full strings.Builder
	full.WriteString("\n")
	full.WriteString("╔══════════════════════════════════════════════╗\n")
	full.WriteString("║         DEVICE DIAGNOSTICS REPORT           ║\n")
	full.WriteString("╚══════════════════════════════════════════════╝\n\n")

	if mem, err := m.GetMemory(deviceID); err == nil {
		full.WriteString(mem + "\n")
	} else {
		full.WriteString("=== MEMORY INFO ===\n  " + err.Error() + "\n\n")
	}

	if cpu, err := m.GetCPU(deviceID); err == nil {
		full.WriteString(cpu + "\n")
	} else {
		full.WriteString("=== CPU INFO ===\n  " + err.Error() + "\n\n")
	}

	if storage, err := m.GetStorage(deviceID); err == nil {
		full.WriteString(storage + "\n")
	} else {
		full.WriteString("=== STORAGE INFO ===\n  " + err.Error() + "\n\n")
	}

	if disp, err := m.GetDisplay(deviceID); err == nil {
		full.WriteString(disp + "\n")
	} else {
		full.WriteString("=== DISPLAY INFO ===\n  " + err.Error() + "\n\n")
	}

	if net, err := m.GetNetwork(deviceID); err == nil {
		full.WriteString(net + "\n")
	} else {
		full.WriteString("=== NETWORK INFO ===\n  " + err.Error() + "\n\n")
	}

	if sensor, err := m.GetSensors(deviceID); err == nil {
		full.WriteString(sensor + "\n")
	} else {
		full.WriteString("=== SENSOR INFO ===\n  " + err.Error() + "\n\n")
	}

	return full.String(), nil
}

func (m *DiagnosticsModel) buildArgs(deviceID string, cmd ...string) []string {
	args := []string{}
	if deviceID != "" {
		args = append(args, "-s", deviceID)
	}
	args = append(args, cmd...)
	return args
}
