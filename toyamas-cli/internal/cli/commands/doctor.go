package commands

import (
	"toyamas-cli/internal/doctor"
	"toyamas-cli/internal/printer"
	"toyamas-cli/internal/tui"
)

func HandleDoctor() {
	if printer.CurrentMode == printer.ModeTUI {
		tui.ShowBanner()
	}
	doctor.RunDoctor()
}
