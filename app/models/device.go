package models

import (
	"fmt"
	"strings"
)

// Device represents an Android device.
type Device struct {
	ID    string
	State string
}

// DeviceModel coordinates device-level actions.
type DeviceModel struct {
	executor *ADBExecutor
}

// NewDeviceModel creates a new DeviceModel.
func NewDeviceModel(executor *ADBExecutor) *DeviceModel {
	return &DeviceModel{executor: executor}
}

// All lists all connected Android devices by running `adb devices`.
func (m *DeviceModel) All() ([]Device, error) {
	stdout, stderr, err := m.executor.Execute("devices")
	if err != nil {
		return nil, fmt.Errorf("failed to run adb devices: %w (stderr: %s)", err, strings.TrimSpace(stderr))
	}

	var devices []Device
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of devices") {
			continue
		}
		// Typically: <device-id>\t<state>
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			devices = append(devices, Device{
				ID:    fields[0],
				State: fields[1],
			})
		}
	}

	return devices, nil
}

// Reboot restarts the device. If deviceID is empty, it reboots the default device.
func (m *DeviceModel) Reboot(deviceID string) (string, error) {
	args := []string{}
	if deviceID != "" {
		args = append(args, "-s", deviceID)
	}
	args = append(args, "reboot")

	stdout, stderr, err := m.executor.Execute(args...)
	if err != nil {
		return "", fmt.Errorf("failed to reboot: %w (stderr: %s)", err, strings.TrimSpace(stderr))
	}
	return stdout, nil
}

// GetDetailedInfo fetches specific device details via getprop and dumpsys.
func (m *DeviceModel) GetDetailedInfo(deviceID string) (map[string]string, error) {
	info := make(map[string]string)

	// Helper to build args
	getArgs := func(shellCmd ...string) []string {
		args := []string{}
		if deviceID != "" {
			args = append(args, "-s", deviceID)
		}
		args = append(args, "shell")
		args = append(args, shellCmd...)
		return args
	}

	// 1. Brand / Manufacturer
	brand, _, _ := m.executor.Execute(getArgs("getprop", "ro.product.brand")...)
	manufacturer, _, _ := m.executor.Execute(getArgs("getprop", "ro.product.manufacturer")...)
	brandStr := strings.TrimSpace(brand)
	manufStr := strings.TrimSpace(manufacturer)
	if brandStr == "" {
		brandStr = manufStr
	} else if manufStr != "" && !strings.EqualFold(brandStr, manufStr) {
		brandStr = brandStr + " (" + manufStr + ")"
	}
	info["Brand"] = brandStr

	// 2. Model & Codename & Marketing Name
	model, _, _ := m.executor.Execute(getArgs("getprop", "ro.product.model")...)
	codename, _, _ := m.executor.Execute(getArgs("getprop", "ro.product.device")...)
	marketName, _, _ := m.executor.Execute(getArgs("getprop", "ro.product.marketname")...)
	
	info["Model"] = strings.TrimSpace(model)
	info["Marketing Name"] = strings.TrimSpace(marketName)
	info["Device Codename"] = strings.TrimSpace(codename)

	// 3. Release Version
	release, _, _ := m.executor.Execute(getArgs("getprop", "ro.build.version.release")...)
	info["Android Version"] = strings.TrimSpace(release)

	// 4. SDK Version
	sdk, _, _ := m.executor.Execute(getArgs("getprop", "ro.build.version.sdk")...)
	info["SDK Version"] = strings.TrimSpace(sdk)

	// 5. Battery Level
	battery, _, _ := m.executor.Execute(getArgs("dumpsys", "battery")...)
	batteryLevel := "Unknown"
	for _, line := range strings.Split(battery, "\n") {
		if strings.Contains(line, "level:") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				batteryLevel = strings.TrimSpace(parts[1]) + "%"
				break
			}
		}
	}
	info["Battery"] = batteryLevel

	// 6. Chipset / Platform
	platform, _, _ := m.executor.Execute(getArgs("getprop", "ro.board.platform")...)
	hardware, _, _ := m.executor.Execute(getArgs("getprop", "ro.hardware")...)
	platStr := strings.TrimSpace(platform)
	hardStr := strings.TrimSpace(hardware)
	chipset := platStr
	if hardStr != "" && !strings.EqualFold(platStr, hardStr) {
		chipset = fmt.Sprintf("%s (%s)", platStr, hardStr)
	}
	info["Chipset"] = chipset

	// 7. Root Status
	rootCheck, _, _ := m.executor.Execute(getArgs("which", "su")...)
	isRooted := "No (Unrooted)"
	if strings.Contains(rootCheck, "su") {
		isRooted = "Yes (Rooted)"
	} else {
		suTest, _, _ := m.executor.Execute(getArgs("su", "-c", "id")...)
		if strings.Contains(suTest, "uid=0") {
			isRooted = "Yes (Rooted)"
		}
	}
	info["Root Access"] = isRooted

	// 8. Bootloader Status (UBL)
	flashLocked, _, _ := m.executor.Execute(getArgs("getprop", "ro.boot.flash.locked")...)
	verifiedState, _, _ := m.executor.Execute(getArgs("getprop", "ro.boot.verifiedbootstate")...)
	lockStr := strings.TrimSpace(flashLocked)
	stateStr := strings.TrimSpace(verifiedState)

	ublStatus := "Locked"
	if lockStr == "0" || stateStr == "orange" {
		ublStatus = "Unlocked (UBL Yes)"
	} else if lockStr == "1" || stateStr == "green" {
		ublStatus = "Locked (UBL No)"
	} else {
		devLock, _, _ := m.executor.Execute(getArgs("getprop", "ro.secureboot.devicelock")...)
		if strings.TrimSpace(devLock) == "0" {
			ublStatus = "Unlocked (UBL Yes)"
		} else {
			ublStatus = "Locked (or Unknown)"
		}
	}
	info["Bootloader (UBL)"] = ublStatus

	return info, nil
}

