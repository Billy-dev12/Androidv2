package views

import (
	"fmt"
	"strings"

	"android-tool-mvc/app/models"
	"os"
	"os/exec"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	styleBold   = "\033[1m"
)

// ConsoleView handles formatting and outputting data to the terminal.
type ConsoleView struct{}

// NewConsoleView creates a new ConsoleView instance.
func NewConsoleView() *ConsoleView {
	return &ConsoleView{}
}

// RenderSuccess prints a beautiful success message.
func (v *ConsoleView) RenderSuccess(message string) {
	fmt.Printf("\n%s%s✔ SUCCESS:%s %s\n\n", colorGreen, styleBold, colorReset, message)
}

// RenderError prints a beautiful error message.
func (v *ConsoleView) RenderError(err error) {
	fmt.Printf("\n%s%s✘ ERROR:%s %s%s%s\n\n", colorRed, styleBold, colorReset, colorRed, err.Error(), colorReset)
}

// RenderDevicesTable prints a styled table of connected devices.
func (v *ConsoleView) RenderDevicesTable(devices []models.Device) {
	if len(devices) == 0 {
		fmt.Printf("\n%sNo connected devices found.%s\n\n", colorYellow, colorReset)
		return
	}

	// Calculate column widths
	maxIDLen := 9 // "DEVICE ID" length
	maxStateLen := 6 // "STATUS" length

	for _, d := range devices {
		if len(d.ID) > maxIDLen {
			maxIDLen = len(d.ID)
		}
		if len(d.State) > maxStateLen {
			maxStateLen = len(d.State)
		}
	}

	padding := 2
	idColWidth := maxIDLen + padding
	stateColWidth := maxStateLen + padding

	// Drawing helper
	lineSeparator := func(left, middle, right, dash string) string {
		return left + strings.Repeat(dash, idColWidth) + middle + strings.Repeat(dash, stateColWidth) + right
	}

	fmt.Println()
	fmt.Println(colorCyan + lineSeparator("┌", "┬", "┐", "─") + colorReset)
	
	// Header row
	idHeader := fmt.Sprintf(" %-*s", idColWidth-1, "DEVICE ID")
	stateHeader := fmt.Sprintf(" %-*s", stateColWidth-1, "STATUS")
	fmt.Printf("%s│%s%s%s%s│%s%s%s%s│%s\n", colorCyan, colorReset, styleBold, idHeader, colorCyan, colorReset, styleBold, stateHeader, colorCyan, colorReset)
	
	fmt.Println(colorCyan + lineSeparator("├", "┼", "┤", "─") + colorReset)

	// Data rows
	for _, d := range devices {
		idVal := fmt.Sprintf(" %-*s", idColWidth-1, d.ID)
		stateColor := colorGreen
		if d.State != "device" {
			stateColor = colorYellow
		}
		stateVal := fmt.Sprintf(" %s%-*s%s", stateColor, stateColWidth-1, d.State, colorReset)
		fmt.Printf("%s│%s%s%s│%s%s\n", colorCyan, colorReset, idVal, colorCyan, stateVal, colorCyan+"│"+colorReset)
	}

	fmt.Println(colorCyan + lineSeparator("└", "┴", "┘", "─") + colorReset)
	fmt.Println()
}

// RenderHelp prints the usage instructions.
func (v *ConsoleView) RenderHelp() {
	fmt.Println()
	fmt.Printf("%s%sAndroid Debug Bridge CLI Wrapper (MVC Pattern)%s\n", colorCyan, styleBold, colorReset)
	fmt.Printf("%sVersion 1.0.0%s\n", colorGray, colorReset)
	fmt.Println()
	fmt.Printf("%sUSAGE:%s\n", styleBold, colorReset)
	fmt.Println("  android-tool <command> [arguments]")
	fmt.Println()
	fmt.Printf("%sAVAILABLE COMMANDS:%s\n", styleBold, colorReset)
	
	commands := []struct {
		Name string
		Args string
		Desc string
	}{
		{"devices", "", "List all connected Android devices"},
		{"info", "[device-id]", "Get detailed information about the device"},
		{"reboot", "[device-id]", "Reboot the default or specified device"},
		{"install", "<apk-path> [device-id]", "Install an APK file on the target device"},
		{"uninstall", "<package-name> [device-id]", "Uninstall a package from the target device"},
		{"push", "<local-path> <remote-path> [device-id]", "Push file/folder to device"},
		{"pull", "<remote-path> <local-path> [device-id]", "Pull file/folder from device"},
		{"screenshot", "[output-path] [device-id]", "Capture device screen to PNG"},
		{"diagnostics", "[device-id]", "Show device diagnostics (memory, CPU, storage, display, network, sensors)"},
		{"firmware partitions", "<folder>", "List partition images (.img/.bin) in extracted folder"},
		{"firmware buildprop", "<folder>", "Show device info from build.prop in extracted folder"},
		{"env", "", "Check system environment (adb, fastboot, lz4, etc)"},
		{"config", "[show|set <key> <value>]", "View or modify configuration"},
		{"history", "", "Show command history log"},
		{"help", "", "Show this help message"},
	}

	for _, cmd := range commands {
		cmdStr := cmd.Name
		if cmd.Args != "" {
			cmdStr += " " + cmd.Args
		}
		fmt.Printf("  %s%-45s%s %s\n", colorGreen, cmdStr, colorReset, cmd.Desc)
	}
	fmt.Println()
}

// RenderInteractiveMenu prints the menu list of features, highlighting the current selection index.
func (v *ConsoleView) RenderInteractiveMenu(title string, options []string, activeIndex int) {
	fmt.Print("\033[H\033[2J") // Clear screen and home cursor
	fmt.Println()
	fmt.Printf("%s%s=== %s ===%s\n", colorCyan, styleBold, title, colorReset)
	fmt.Printf("%sNavigate with Up/Down Arrow, press Enter to select, or press numbers (1-%d). Press 'q' to exit/back.%s\n\n", colorGray, len(options), colorReset)

	for i, option := range options {
		if i == activeIndex {
			fmt.Printf(" %s%s➔ [%d] %s%s\n", colorGreen, styleBold, i+1, option, colorReset)
		} else {
			fmt.Printf("   %s[%d] %s%s\n", colorGray, i+1, option, colorReset)
		}
	}
	fmt.Println()
}

// SetRawMode configures the terminal to raw/cbreak mode or restores sane settings.
func (v *ConsoleView) SetRawMode(enable bool) {
	if enable {
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	} else {
		exec.Command("stty", "-F", "/dev/tty", "sane").Run()
	}
}

// ReadInputChar reads a single keypress. It handles normal keys and parses arrow escape sequences.
func (v *ConsoleView) ReadInputChar() (string, error) {
	var buf [3]byte
	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		return "", err
	}

	if n == 1 {
		b := buf[0]
		if b == 13 || b == 10 {
			return "enter", nil
		}
		if b >= '1' && b <= '9' {
			return string(b), nil
		}
		if b == 'q' || b == 'Q' {
			return "q", nil
		}
	} else if n == 3 {
		if buf[0] == 27 && buf[1] == 91 {
			switch buf[2] {
			case 'A':
				return "up", nil
			case 'B':
				return "down", nil
			}
		}
	}
	return "other", nil
}

// PromptInput requests a string input from the user (restoring normal stdin temporarily).
func (v *ConsoleView) PromptInput(promptText string) string {
	v.SetRawMode(false) // Temporarily disable raw mode
	defer v.SetRawMode(true) // Re-enable raw mode when returning

	fmt.Printf("%s%s%s %s", colorCyan, styleBold, promptText, colorReset)
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

// RenderSystemInfo prints a styled system/config/history report.
func (v *ConsoleView) RenderSystemInfo(report string) {
	fmt.Println()
	for _, line := range strings.Split(report, "\n") {
		if strings.HasPrefix(line, "===") {
			fmt.Printf("%s%s%s\n", colorCyan, line, colorReset)
		} else if strings.Contains(line, "✓") {
			fmt.Printf("  %s%s%s\n", colorGreen, strings.TrimSpace(line), colorReset)
		} else if strings.Contains(line, "✗") {
			fmt.Printf("  %s%s%s\n", colorRed, strings.TrimSpace(line), colorReset)
		} else if strings.HasPrefix(line, "  ") {
			fmt.Printf("%s%s%s\n", colorGray, line, colorReset)
		} else {
			fmt.Println(line)
		}
	}
	fmt.Println()
}

// RenderDiagnostics prints a styled diagnostics report.
func (v *ConsoleView) RenderDiagnostics(report string) {
	fmt.Println()
	for _, line := range strings.Split(report, "\n") {
		if strings.HasPrefix(line, "===") {
			fmt.Printf("%s%s%s\n", colorCyan, line, colorReset)
		} else if strings.HasPrefix(line, "  ") {
			fmt.Printf("%s%s\n", colorGray, line)
		} else {
			fmt.Println(line)
		}
	}
	fmt.Println()
}

// RenderPartitionInfo prints a table of partition image files.
func (v *ConsoleView) RenderPartitionInfo(partitions []models.PartitionInfo) {
	fmt.Println()
	fmt.Printf("%s%s=== PARTITION INFORMATION ===%s\n", colorCyan, styleBold, colorReset)

	maxNameLen := 4
	maxSizeLen := 4
	maxTypeLen := 4
	for _, p := range partitions {
		if len(p.Name) > maxNameLen {
			maxNameLen = len(p.Name)
		}
		if len(p.SizeHuman) > maxSizeLen {
			maxSizeLen = len(p.SizeHuman)
		}
		if len(p.FileType) > maxTypeLen {
			maxTypeLen = len(p.FileType)
		}
	}
	padding := 2
	nameCol := maxNameLen + padding
	sizeCol := maxSizeLen + padding
	typeCol := maxTypeLen + padding

	sep := func(l, m, r string) string {
		return l + strings.Repeat("─", nameCol) + m + strings.Repeat("─", sizeCol) + m + strings.Repeat("─", typeCol) + r
	}

	fmt.Println(colorCyan + sep("┌", "┬", "┐") + colorReset)
	hName := fmt.Sprintf(" %-*s", nameCol-1, "FILE")
	hSize := fmt.Sprintf(" %-*s", sizeCol-1, "SIZE")
	hType := fmt.Sprintf(" %-*s", typeCol-1, "TYPE")
	fmt.Printf("%s│%s%s%s%s│%s%s%s%s│%s%s%s%s│%s\n", colorCyan, colorReset, styleBold, hName, colorCyan, colorReset, styleBold, hSize, colorCyan, colorReset, styleBold, hType, colorCyan, colorReset)
	fmt.Println(colorCyan + sep("├", "┼", "┤") + colorReset)

	for _, p := range partitions {
		rName := fmt.Sprintf(" %-*s", nameCol-1, p.Name)
		rSize := fmt.Sprintf(" %-*s", sizeCol-1, p.SizeHuman)
		rType := fmt.Sprintf(" %-*s", typeCol-1, p.FileType)
		fmt.Printf("%s│%s%s%s│%s%s%s│%s%s%s│%s\n", colorCyan, colorReset, rName, colorCyan, colorReset, rSize, colorCyan, colorReset, rType, colorCyan, colorReset)
	}

	fmt.Println(colorCyan + sep("└", "┴", "┘") + colorReset)
	fmt.Println()

	// Summary
	var totalSize int64
	for _, p := range partitions {
		totalSize += p.Size
	}
	var totalHuman string
	switch {
	case totalSize >= 1<<30:
		totalHuman = fmt.Sprintf("%.2f GB", float64(totalSize)/(1<<30))
	case totalSize >= 1<<20:
		totalHuman = fmt.Sprintf("%.2f MB", float64(totalSize)/(1<<20))
	default:
		totalHuman = fmt.Sprintf("%d B", totalSize)
	}
	fmt.Printf("  %sTotal:%s %d partition(s), %s\n\n", colorGray, colorReset, len(partitions), totalHuman)
}

// RenderBuildProp prints parsed build.prop information.
func (v *ConsoleView) RenderBuildProp(props map[string]string) {
	interestingKeys := []struct {
		Key     string
		Label   string
	}{
		{"ro.product.model", "Model"},
		{"ro.product.marketname", "Marketing Name"},
		{"ro.product.name", "Product Name"},
		{"ro.product.board", "Board"},
		{"ro.product.cpu.abi", "CPU ABI"},
		{"ro.build.version.release", "Android Version"},
		{"ro.build.version.sdk", "SDK Level"},
		{"ro.build.version.security_patch", "Security Patch"},
		{"ro.build.date", "Build Date"},
		{"ro.build.fingerprint", "Build Fingerprint"},
		{"ro.build.description", "Build Description"},
		{"ro.product.manufacturer", "Manufacturer"},
		{"ro.product.brand", "Brand"},
		{"ro.hardware", "Hardware"},
		{"ro.soc.model", "SoC Model"},
		{"persist.sys.timezone", "Timezone"},
	}

	fmt.Println()
	fmt.Printf("%s%s=== BUILD.PROP INFORMATION ===%s\n", colorCyan, styleBold, colorReset)
	fmt.Printf("  Property source: %s%d properties loaded%s\n\n", colorGray, len(props), colorReset)

	hasAny := false
	for _, item := range interestingKeys {
		val, exists := props[item.Key]
		if !exists || val == "" {
			continue
		}
		hasAny = true
		fmt.Printf("  %s%-18s:%s %s\n", colorGreen, item.Label, colorReset, val)
	}

	if !hasAny {
		fmt.Printf("  %sNo relevant device properties found.%s\n", colorYellow, colorReset)
	}
	fmt.Println()
}

// RenderDeviceInfo prints a styled summary of the device properties.
func (v *ConsoleView) RenderDeviceInfo(info map[string]string) {
	fmt.Println()
	fmt.Printf("%s%s=== DEVICE INFORMATION ===%s\n", colorCyan, styleBold, colorReset)
	keys := []string{"Brand", "Marketing Name", "Model", "Device Codename", "Chipset", "Android Version", "SDK Version", "Root Access", "Bootloader (UBL)", "Battery"}
	for _, key := range keys {
		val, exists := info[key]
		if !exists || val == "" {
			if key == "Marketing Name" || key == "Device Codename" {
				continue
			}
			val = "N/A"
		}
		fmt.Printf("  %s%-18s:%s %s\n", colorGreen, key, colorReset, val)
	}
	fmt.Println()
}

