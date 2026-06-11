package routes

import (
	"fmt"
	"strings"

	"android-tool-mvc/app/controllers"
	"android-tool-mvc/app/models"
	"android-tool-mvc/resources/views"
)

// Router handles CLI argument parsing and controller routing.
type Router struct {
	deviceController      *controllers.DeviceController
	appController         *controllers.AppController
	fileController        *controllers.FileController
	firmwareController    *controllers.FirmwareController
	diagnosticsController *controllers.DiagnosticsController
	systemController      *controllers.SystemController
	view                  *views.ConsoleView
	executor              *models.ADBExecutor
	history               *models.HistoryModel
	config                *models.ConfigModel
}

// NewRouter creates and initializes a Router.
func NewRouter(
	deviceController *controllers.DeviceController,
	appController *controllers.AppController,
	fileController *controllers.FileController,
	firmwareController *controllers.FirmwareController,
	diagnosticsController *controllers.DiagnosticsController,
	systemController *controllers.SystemController,
	view *views.ConsoleView,
	executor *models.ADBExecutor,
	history *models.HistoryModel,
	config *models.ConfigModel,
) *Router {
	_ = config.Load()
	return &Router{
		deviceController:      deviceController,
		appController:         appController,
		fileController:        fileController,
		firmwareController:    firmwareController,
		diagnosticsController: diagnosticsController,
		systemController:      systemController,
		view:                  view,
		executor:              executor,
		history:               history,
		config:                config,
	}
}

// Route directs the CLI arguments to the correct controller actions.
func (r *Router) Route(args []string) {
	// Check if ADB is installed before routing any operation
	if !r.executor.IsInstalled() {
		r.view.RenderError(fmt.Errorf("ADB (Android Debug Bridge) is not installed or not found in system PATH"))
		return
	}

	if len(args) < 2 {
		r.enterInteractiveMode()
		return
	}

	command := args[1]

	switch command {
	case "devices", "list":
		r.history.Append("devices")
		r.deviceController.Index()

	case "info":
		var deviceID string
		if len(args) > 2 {
			deviceID = args[2]
		}
		r.history.Append("info", deviceID)
		r.deviceController.ShowDeviceInfo(deviceID)

	case "reboot":
		var deviceID string
		if len(args) > 2 {
			deviceID = args[2]
		}
		r.history.Append("reboot", deviceID)
		r.deviceController.Reboot(deviceID)

	case "install":
		if len(args) < 3 {
			r.view.RenderError(fmt.Errorf("missing argument: install command requires an APK path\nUsage: android-tool install <apk-path> [device-id]"))
			return
		}
		apkPath := args[2]
		var deviceID string
		if len(args) > 3 {
			deviceID = args[3]
		}
		r.history.Append("install", apkPath, deviceID)
		r.appController.Install(apkPath, deviceID)

	case "force-install":
		if len(args) < 3 {
			r.view.RenderError(fmt.Errorf("missing argument: force-install command requires an APK path\nUsage: android-tool force-install <apk-path> [device-id]"))
			return
		}
		apkPath := args[2]
		var deviceID string
		if len(args) > 3 {
			deviceID = args[3]
		}
		r.history.Append("force-install", apkPath, deviceID)
		r.appController.InstallForce(apkPath, deviceID)

	case "uninstall":
		if len(args) < 3 {
			r.view.RenderError(fmt.Errorf("missing argument: uninstall command requires a package name\nUsage: android-tool uninstall <package-name> [device-id]"))
			return
		}
		packageName := args[2]
		var deviceID string
		if len(args) > 3 {
			deviceID = args[3]
		}
		r.history.Append("uninstall", packageName, deviceID)
		r.appController.Uninstall(packageName, deviceID)

	case "push":
		if len(args) < 4 {
			r.view.RenderError(fmt.Errorf("missing arguments: push command requires local-path and remote-path\nUsage: android-tool push <local-path> <remote-path> [device-id]"))
			return
		}
		localPath := args[2]
		remotePath := args[3]
		var deviceID string
		if len(args) > 4 {
			deviceID = args[4]
		}
		r.history.Append("push", localPath, remotePath, deviceID)
		r.fileController.Push(localPath, remotePath, deviceID)

	case "pull":
		if len(args) < 4 {
			r.view.RenderError(fmt.Errorf("missing arguments: pull command requires remote-path and local-path\nUsage: android-tool pull <remote-path> <local-path> [device-id]"))
			return
		}
		remotePath := args[2]
		localPath := args[3]
		var deviceID string
		if len(args) > 4 {
			deviceID = args[4]
		}
		r.history.Append("pull", remotePath, localPath, deviceID)
		r.fileController.Pull(remotePath, localPath, deviceID)

	case "screenshot":
		var outputPath string
		var deviceID string
		if len(args) > 2 {
			outputPath = args[2]
		}
		if len(args) > 3 {
			deviceID = args[3]
		}
		r.history.Append("screenshot", outputPath, deviceID)
		r.deviceController.Screenshot(deviceID, outputPath)

	case "diagnostics":
		var deviceID string
		if len(args) > 2 {
			deviceID = args[2]
		}
		r.history.Append("diagnostics", deviceID)
		r.diagnosticsController.ShowAll(deviceID)

	case "env":
		r.history.Append("env")
		r.systemController.ShowEnv()

	case "config":
		r.history.Append("config", args[2:]...)
		if len(args) == 2 {
			r.systemController.ShowConfig()
		} else if len(args) >= 4 && args[2] == "set" {
			r.systemController.SetConfig(args[3], strings.Join(args[4:], " "))
		} else {
			r.view.RenderError(fmt.Errorf("usage: config [show|set <key> <value>]"))
		}

	case "history":
		r.systemController.ShowHistory()

	case "firmware":
		r.history.Append("firmware", args[2:]...)
		if len(args) < 3 {
			r.view.RenderError(fmt.Errorf("missing subcommand\nUsage: android-tool firmware partitions <folder-path>"))
			return
		}
		subcommand := args[2]
		var folderPath string
		if len(args) > 3 {
			folderPath = args[3]
		}
		switch subcommand {
		case "partitions":
			r.firmwareController.ShowPartitionInfo(folderPath)
		default:
			r.view.RenderError(fmt.Errorf("unknown firmware subcommand: %s\nAvailable: partitions", subcommand))
		}

	case "help", "-h", "--help":
		r.view.RenderHelp()

	default:
		r.view.RenderError(fmt.Errorf("unknown command: %s", command))
		r.view.RenderHelp()
	}
}

// enterInteractiveMode enters the interactive TUI menu loop.
func (r *Router) enterInteractiveMode() {
	r.view.SetRawMode(true)
	defer r.view.SetRawMode(false)

	mainOptions := []string{
		"ADB Management Menu",
		"Fastboot Management Menu (Soon)",
		"Firmware Extractor",
		"Device Diagnostics",
		"System Tools",
		"MediaTek Port Monitor (BROM) (Soon)",
		"Exit",
	}

	adbOptions := []string{
		"List Devices",
		"Device Info",
		"Reboot Device",
		"Install APK",
		"Force Install APK",
		"Uninstall Package",
		"Push File (Local -> Device)",
		"Pull File (Device -> Local)",
		"Screenshot",
		"Back to Main Menu",
	}

	firmwareOptions := []string{
		"Extract Outer Archive (.zip, .tgz, .tar.gz, .tar, .tar.md5)",
		"Samsung Firmware Extractor (Inner)",
		"Partition Info Report (scan extracted folder)",
		"Back to Main Menu",
	}

	diagnosticsOptions := []string{
		"Full Diagnostics Report",
		"Memory Info",
		"CPU Info",
		"Storage Info",
		"Display Info",
		"Network & Signal Info",
		"Sensor Info",
		"Back to Main Menu",
	}

	systemOptions := []string{
		"Environment Check",
		"Configuration",
		"Command History",
		"Back to Main Menu",
	}

	currentMenu := "main"
	activeIndex := 0

	for {
		var title string
		var options []string

		if currentMenu == "main" {
			title = "ANDROID V2 CORE ENGINE"
			options = mainOptions
		} else if currentMenu == "adb" {
			title = "ADB MANAGEMENT"
			options = adbOptions
		} else if currentMenu == "firmware" {
			title = "FIRMWARE EXTRACTOR"
			options = firmwareOptions
		} else if currentMenu == "diagnostics" {
			title = "DEVICE DIAGNOSTICS"
			options = diagnosticsOptions
		} else if currentMenu == "system" {
			title = "SYSTEM TOOLS"
			options = systemOptions
		}

		r.view.RenderInteractiveMenu(title, options, activeIndex)

		key, err := r.view.ReadInputChar()
		if err != nil {
			break
		}

		if key == "q" {
			if currentMenu == "adb" || currentMenu == "firmware" || currentMenu == "diagnostics" || currentMenu == "system" {
				currentMenu = "main"
				activeIndex = 0
				continue
			} else {
				break
			}
		}

		switch key {
		case "up":
			activeIndex--
			if activeIndex < 0 {
				activeIndex = len(options) - 1
			}
		case "down":
			activeIndex++
			if activeIndex >= len(options) {
				activeIndex = 0
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			idx := int(key[0] - '1')
			if idx < len(options) {
				activeIndex = idx
				if currentMenu == "main" {
					if activeIndex == 0 {
						currentMenu = "adb"
						activeIndex = 0
					} else if activeIndex == 2 {
						currentMenu = "firmware"
						activeIndex = 0
					} else if activeIndex == 3 {
						currentMenu = "diagnostics"
						activeIndex = 0
					} else if activeIndex == 4 {
						currentMenu = "system"
						activeIndex = 0
					} else if activeIndex == 6 {
						r.view.SetRawMode(false)
						fmt.Println("\nGoodbye!")
						return
					} else {
						r.view.SetRawMode(false)
						fmt.Printf("\nFeature '%s' is not implemented yet.\nPress Enter to return...", options[activeIndex])
						var dummy string
						fmt.Scanln(&dummy)
						r.view.SetRawMode(true)
					}
				} else if currentMenu == "adb" {
					r.executeADBAction(activeIndex, &currentMenu)
					if currentMenu == "main" {
						activeIndex = 0
					}
				} else if currentMenu == "firmware" {
					r.executeFirmwareAction(activeIndex, &currentMenu)
					if currentMenu == "main" {
						activeIndex = 0
					}
				} else if currentMenu == "diagnostics" {
					r.executeDiagnosticsAction(activeIndex, &currentMenu)
					if currentMenu == "main" {
						activeIndex = 0
					}
				} else if currentMenu == "system" {
					r.executeSystemAction(activeIndex, &currentMenu)
					if currentMenu == "main" {
						activeIndex = 0
					}
				}
			}
		case "enter":
			if currentMenu == "main" {
				if activeIndex == 0 {
					currentMenu = "adb"
					activeIndex = 0
				} else if activeIndex == 2 {
					currentMenu = "firmware"
					activeIndex = 0
				} else if activeIndex == 3 {
					currentMenu = "diagnostics"
					activeIndex = 0
				} else if activeIndex == 4 {
					currentMenu = "system"
					activeIndex = 0
				} else if activeIndex == 6 {
					r.view.SetRawMode(false)
					fmt.Println("\nGoodbye!")
					return
				} else {
					r.view.SetRawMode(false)
					fmt.Printf("\nFeature '%s' is not implemented yet.\nPress Enter to return...", options[activeIndex])
					var dummy string
					fmt.Scanln(&dummy)
					r.view.SetRawMode(true)
				}
			} else if currentMenu == "adb" {
				r.executeADBAction(activeIndex, &currentMenu)
				if currentMenu == "main" {
					activeIndex = 0
				}
			} else if currentMenu == "firmware" {
				r.executeFirmwareAction(activeIndex, &currentMenu)
				if currentMenu == "main" {
					activeIndex = 0
				}
			} else if currentMenu == "diagnostics" {
				r.executeDiagnosticsAction(activeIndex, &currentMenu)
				if currentMenu == "main" {
					activeIndex = 0
				}
			} else if currentMenu == "system" {
				r.executeSystemAction(activeIndex, &currentMenu)
				if currentMenu == "main" {
					activeIndex = 0
				}
			}
		}
	}
}

// executeADBAction runs the selected ADB command.
func (r *Router) executeADBAction(index int, currentMenu *string) {
	r.view.SetRawMode(false)
	fmt.Println()

	switch index {
	case 0: // List Devices
		r.history.Append("devices")
		r.deviceController.Index()
	case 1: // Device Info
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("info", deviceID)
		r.deviceController.ShowDeviceInfo(deviceID)
	case 2: // Reboot Device
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("reboot", deviceID)
		r.deviceController.Reboot(deviceID)
	case 3: // Install APK
		apkPath := r.view.PromptInput("Enter APK Path: ")
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("install", apkPath, deviceID)
		r.appController.Install(apkPath, deviceID)
	case 4: // Force Install APK
		apkPath := r.view.PromptInput("Enter APK Path: ")
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("force-install", apkPath, deviceID)
		r.appController.InstallForce(apkPath, deviceID)
	case 5: // Uninstall Package
		pkgName := r.view.PromptInput("Enter Package Name: ")
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("uninstall", pkgName, deviceID)
		r.appController.Uninstall(pkgName, deviceID)
	case 6: // Push File
		fmt.Println("\033[90m💡 [Petunjuk Jalur Berkas / Path Helper]\033[0m")
		fmt.Println("   • \033[36mLocal Path\033[0m  (PC Anda): File/folder yang ingin diunggah.")
		fmt.Println("     Contoh: \033[32m./file.txt\033[0m (di folder saat ini) atau \033[32m/home/billy/Downloads/app.apk\033[0m")
		fmt.Println("   • \033[36mRemote Path\033[0m (Android): Lokasi tujuan di perangkat Android.")
		fmt.Println("     Contoh: \033[32m/sdcard/Download/\033[0m (Penyimpanan Internal) atau \033[32m/data/local/tmp/\033[0m (Folder Temp)")
		fmt.Println()
		localPath := r.view.PromptInput("Masukkan Local File Path: ")
		remotePath := r.view.PromptInput("Masukkan Remote Destination Path: ")
		deviceID := r.view.PromptInput("Masukkan Device ID (kosongkan untuk default): ")
		r.history.Append("push", localPath, remotePath, deviceID)
		r.fileController.Push(localPath, remotePath, deviceID)
	case 7: // Pull File
		fmt.Println("\033[90m💡 [Petunjuk Jalur Berkas / Path Helper]\033[0m")
		fmt.Println("   • \033[36mRemote Path\033[0m (Android): File/folder yang ingin diambil dari HP.")
		fmt.Println("     Contoh: \033[32m/sdcard/Download/foto.jpg\033[0m atau \033[32m/sdcard/DCIM/Camera/\033[0m")
		fmt.Println("   • \033[36mLocal Path\033[0m  (PC Anda): Folder tujuan penyimpanan di PC.")
		fmt.Println("     Contoh: \033[32m./\033[0m (simpan di folder saat ini) atau \033[32m/home/billy/Desktop/\033[0m")
		fmt.Println()
		remotePath := r.view.PromptInput("Masukkan Remote File Path: ")
		localPath := r.view.PromptInput("Masukkan Local Destination Path: ")
		deviceID := r.view.PromptInput("Masukkan Device ID (kosongkan untuk default): ")
		r.history.Append("pull", remotePath, localPath, deviceID)
		r.fileController.Pull(remotePath, localPath, deviceID)
	case 8: // Screenshot
		fmt.Println("\033[90m📸 [Screenshot Capture]\033[0m")
		fmt.Println("   Akan menyimpan screenshot ke file PNG di folder saat ini.")
		fmt.Println()
		outputPath := r.view.PromptInput("Masukkan Output Path (kosongkan untuk default): ")
		deviceID := r.view.PromptInput("Masukkan Device ID (kosongkan untuk default): ")
		r.history.Append("screenshot", outputPath, deviceID)
		r.deviceController.Screenshot(deviceID, outputPath)
	case 9: // Back
		*currentMenu = "main"
		r.view.SetRawMode(true)
		return
	}

	fmt.Printf("\nPress Enter to return to menu...")
	var dummy string
	fmt.Scanln(&dummy)
	r.view.SetRawMode(true)
}

// executeSystemAction runs the selected system tools action.
func (r *Router) executeSystemAction(index int, currentMenu *string) {
	r.view.SetRawMode(false)
	fmt.Println()

	switch index {
	case 0: // Environment Check
		r.history.Append("env")
		r.systemController.ShowEnv()
	case 1: // Configuration
		r.systemController.ShowConfig()
	case 2: // Command History
		r.systemController.ShowHistory()
	case 3: // Back
		*currentMenu = "main"
		r.view.SetRawMode(true)
		return
	}

	fmt.Printf("\nPress Enter to return to menu...")
	var dummy string
	fmt.Scanln(&dummy)
	r.view.SetRawMode(true)
}

// executeDiagnosticsAction runs the selected diagnostics action.
func (r *Router) executeDiagnosticsAction(index int, currentMenu *string) {
	r.view.SetRawMode(false)
	fmt.Println()

	switch index {
	case 0: // Full Diagnostics Report
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("diagnostics", deviceID)
		r.diagnosticsController.ShowAll(deviceID)
	case 1: // Memory Info
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("diagnostics memory", deviceID)
		r.diagnosticsController.ShowMemory(deviceID)
	case 2: // CPU Info
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("diagnostics cpu", deviceID)
		r.diagnosticsController.ShowCPU(deviceID)
	case 3: // Storage Info
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("diagnostics storage", deviceID)
		r.diagnosticsController.ShowStorage(deviceID)
	case 4: // Display Info
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("diagnostics display", deviceID)
		r.diagnosticsController.ShowDisplay(deviceID)
	case 5: // Network & Signal Info
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("diagnostics network", deviceID)
		r.diagnosticsController.ShowNetwork(deviceID)
	case 6: // Sensor Info
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.history.Append("diagnostics sensor", deviceID)
		r.diagnosticsController.ShowSensors(deviceID)
	case 7: // Back
		*currentMenu = "main"
		r.view.SetRawMode(true)
		return
	}

	fmt.Printf("\nPress Enter to return to menu...")
	var dummy string
	fmt.Scanln(&dummy)
	r.view.SetRawMode(true)
}

// executeFirmwareAction runs the selected firmware extractor action.
func (r *Router) executeFirmwareAction(index int, currentMenu *string) {
	r.view.SetRawMode(false)
	fmt.Println()

	switch index {
	case 0: // Outer Extractor
		filePath := r.view.PromptInput("Masukkan Path File Firmware (.zip/.tgz/.tar.gz/.tar/.tar.md5): ")
		outputDir := r.view.PromptInput("Masukkan Folder Output (kosongkan untuk default): ")
		r.firmwareController.ExtractOuterArchive(filePath, outputDir)
	case 1: // Samsung Inner Extractor
		r.firmwareController.ExtractSamsungInner()
	case 2: // Partition Info Report
		folderPath := r.view.PromptInput("Masukkan Path Folder hasil extract firmware: ")
		r.firmwareController.ShowPartitionInfo(folderPath)
	case 3: // Back
		*currentMenu = "main"
		r.view.SetRawMode(true)
		return
	}

	fmt.Printf("\nPress Enter to return to menu...")
	var dummy string
	fmt.Scanln(&dummy)
	r.view.SetRawMode(true)
}
