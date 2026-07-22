package commands

import (
	"os/exec"

	"toyamas-cli/internal/autoupdate"
	"toyamas-cli/internal/printer"
	"toyamas-cli/internal/tui"
)

func HandleUpdate(version string) {
	if printer.CurrentMode == printer.ModeTUI {
		tui.ShowBanner()
	}

	spinner := tui.StartSpinner("Updating Debian package repository and system dependencies")

	cmd := exec.Command("apt-get", "update", "-y")
	_ = cmd.Run()

	cmdUpgrade := exec.Command("apt-get", "dist-upgrade", "-y", "--no-install-recommends")
	_ = cmdUpgrade.Run()

	spinner.Stop("System packages updated")

	// Self update CLI
	spinnerCLI := tui.StartSpinner("Checking Toyamas CLI updates from remote release stream")
	err := autoupdate.PerformUpdate(version)
	spinnerCLI.Stop("")

	if err != nil {
		printer.LogWarn("CLI self-update warning: %v", err)
	}

	if printer.CurrentMode == printer.ModeJSON {
		printer.PrintJSON(map[string]interface{}{"status": "success", "updated": true, "version": version})
	} else {
		printer.LogSuccess("ToyamasOS system and CLI update completed.")
	}
}
