package commands

import (
	"bufio"
	"os"

	"toyamas-cli/internal/printer"
	"toyamas-cli/internal/tui"
)

func HandleLogs() {
	logFile := "/var/log/toyamas-installer.log"

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		printer.LogWarn("Toyamas log file not found at '%s'. Displaying mock installer logs.", logFile)
		if printer.CurrentMode == printer.ModeJSON {
			printer.PrintJSON(map[string]interface{}{
				"logs": []string{
					"[2026-07-22 16:50:00] [INFO] Step 1: System update completed.",
					"[2026-07-22 16:50:10] [SUCCESS] Docker & Compose installed.",
					"[2026-07-22 16:50:20] [SUCCESS] ZRAM compressed swap enabled.",
				},
			})
			return
		}

		tui.ShowBanner()
		printer.LogInfo("=== Toyamas System Logs ===")
		printer.PrintTUI("[2026-07-22 16:50:00] [INFO] Step 1: System update completed.")
		printer.PrintTUI("[2026-07-22 16:50:10] [SUCCESS] Docker Engine & Compose plugin active.")
		printer.PrintTUI("[2026-07-22 16:50:20] [SUCCESS] ZRAM compressed swap enabled (50%% RAM cap).")
		printer.PrintTUI("[2026-07-22 16:50:30] [SUCCESS] Security stack (UFW & Fail2Ban) active.")
		return
	}

	file, err := os.Open(logFile)
	if err != nil {
		printer.LogError("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if printer.CurrentMode == printer.ModeJSON {
		printer.PrintJSON(map[string]interface{}{"log_file": logFile, "lines": lines})
	} else {
		tui.ShowBanner()
		printer.LogInfo("=== Log File: %s ===", logFile)
		for _, line := range lines {
			printer.PrintTUI("%s", line)
		}
	}
}
