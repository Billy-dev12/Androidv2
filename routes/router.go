package routes

import (
	"fmt"

	"android-tool-mvc/app/controllers"
	"android-tool-mvc/app/models"
	"android-tool-mvc/resources/views"
)

// Router handles CLI argument parsing and controller routing.
type Router struct {
	deviceController *controllers.DeviceController
	appController    *controllers.AppController
	fileController   *controllers.FileController
	view             *views.ConsoleView
	executor         *models.ADBExecutor
}

// NewRouter creates and initializes a Router.
func NewRouter(
	deviceController *controllers.DeviceController,
	appController *controllers.AppController,
	fileController *controllers.FileController,
	view *views.ConsoleView,
	executor *models.ADBExecutor,
) *Router {
	return &Router{
		deviceController: deviceController,
		appController:    appController,
		fileController:   fileController,
		view:             view,
		executor:         executor,
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
		r.deviceController.Index()

	case "info":
		var deviceID string
		if len(args) > 2 {
			deviceID = args[2]
		}
		r.deviceController.ShowDeviceInfo(deviceID)

	case "reboot":
		var deviceID string
		if len(args) > 2 {
			deviceID = args[2]
		}
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
		r.appController.Install(apkPath, deviceID)

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
		r.fileController.Pull(remotePath, localPath, deviceID)

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
		"Xiaomi Firmware Extractor (Soon)",
		"MediaTek Port Monitor (BROM) (Soon)",
		"Exit",
	}

	adbOptions := []string{
		"List Devices",
		"Device Info",
		"Reboot Device",
		"Install APK",
		"Uninstall Package",
		"Push File (Local -> Device)",
		"Pull File (Device -> Local)",
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
		} else {
			title = "ADB MANAGEMENT"
			options = adbOptions
		}

		r.view.RenderInteractiveMenu(title, options, activeIndex)

		key, err := r.view.ReadInputChar()
		if err != nil {
			break
		}

		if key == "q" {
			if currentMenu == "adb" {
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
		case "1", "2", "3", "4", "5", "6", "7", "8":
			idx := int(key[0] - '1')
			if idx < len(options) {
				activeIndex = idx
				if currentMenu == "main" {
					if activeIndex == 0 {
						currentMenu = "adb"
						activeIndex = 0
					} else if activeIndex == 4 {
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
				} else {
					r.executeADBAction(activeIndex, &currentMenu)
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
				} else if activeIndex == 4 {
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
			} else {
				r.executeADBAction(activeIndex, &currentMenu)
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
		r.deviceController.Index()
	case 1: // Device Info
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.deviceController.ShowDeviceInfo(deviceID)
	case 2: // Reboot Device
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.deviceController.Reboot(deviceID)
	case 3: // Install APK
		apkPath := r.view.PromptInput("Enter APK Path: ")
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.appController.Install(apkPath, deviceID)
	case 4: // Uninstall Package
		pkgName := r.view.PromptInput("Enter Package Name: ")
		deviceID := r.view.PromptInput("Enter Device ID (leave empty for default): ")
		r.appController.Uninstall(pkgName, deviceID)
	case 5: // Push File
		fmt.Println("\033[90m💡 [Petunjuk Jalur Berkas / Path Helper]\033[0m")
		fmt.Println("   • \033[36mLocal Path\033[0m  (PC Anda): File/folder yang ingin diunggah.")
		fmt.Println("     Contoh: \033[32m./file.txt\033[0m (di folder saat ini) atau \033[32m/home/billy/Downloads/app.apk\033[0m")
		fmt.Println("   • \033[36mRemote Path\033[0m (Android): Lokasi tujuan di perangkat Android.")
		fmt.Println("     Contoh: \033[32m/sdcard/Download/\033[0m (Penyimpanan Internal) atau \033[32m/data/local/tmp/\033[0m (Folder Temp)")
		fmt.Println()
		localPath := r.view.PromptInput("Masukkan Local File Path: ")
		remotePath := r.view.PromptInput("Masukkan Remote Destination Path: ")
		deviceID := r.view.PromptInput("Masukkan Device ID (kosongkan untuk default): ")
		r.fileController.Push(localPath, remotePath, deviceID)
	case 6: // Pull File
		fmt.Println("\033[90m💡 [Petunjuk Jalur Berkas / Path Helper]\033[0m")
		fmt.Println("   • \033[36mRemote Path\033[0m (Android): File/folder yang ingin diambil dari HP.")
		fmt.Println("     Contoh: \033[32m/sdcard/Download/foto.jpg\033[0m atau \033[32m/sdcard/DCIM/Camera/\033[0m")
		fmt.Println("   • \033[36mLocal Path\033[0m  (PC Anda): Folder tujuan penyimpanan di PC.")
		fmt.Println("     Contoh: \033[32m./\033[0m (simpan di folder saat ini) atau \033[32m/home/billy/Desktop/\033[0m")
		fmt.Println()
		remotePath := r.view.PromptInput("Masukkan Remote File Path: ")
		localPath := r.view.PromptInput("Masukkan Local Destination Path: ")
		deviceID := r.view.PromptInput("Masukkan Device ID (kosongkan untuk default): ")
		r.fileController.Pull(remotePath, localPath, deviceID)
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
