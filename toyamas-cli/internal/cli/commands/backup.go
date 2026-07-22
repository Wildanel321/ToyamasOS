package commands

import (
	"toyamas-cli/internal/backup"
	"toyamas-cli/internal/printer"
	"toyamas-cli/internal/tui"
)

func HandleBackup() {
	if printer.CurrentMode == printer.ModeTUI {
		tui.ShowBanner()
	}

	_, err := backup.CreateBackup()
	if err != nil {
		printer.LogError("Backup failed: %v", err)
	}
}
