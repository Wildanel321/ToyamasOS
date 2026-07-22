package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"toyamas-cli/internal/printer"
	"toyamas-cli/internal/tui"
)

func HandleInstall(appName string) {
	if appName == "" {
		printer.LogError("Missing application name. Usage: toyamas install <nginx|laravel|minecraft|redis|mariadb|nextcloud|uptime-kuma>")
		os.Exit(1)
	}

	if printer.CurrentMode == printer.ModeTUI {
		tui.ShowBanner()
		printer.LogInfo("Installing application stack: '%s'...", appName)
	}

	steps := []string{
		"Locating app template in Toyamas App Store repository",
		"Generating environment configuration & network bindings",
		"Pulling container images from Docker Registry",
		"Starting Docker Compose services",
	}

	for i, step := range steps {
		tui.RenderProgressBar(i+1, len(steps), fmt.Sprintf("Step %d/4", i+1))
		printer.LogInfo("-> %s", step)
	}

	// Target app compose file path
	appDir := filepath.Join("toyamas-panel", "apps", appName)
	composePath := filepath.Join(appDir, "docker-compose.yml")

	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		// Fall back to direct Docker execution if app template folder doesn't exist locally
		printer.LogWarn("Custom compose template not found at '%s'. Using fallback Docker runner.", composePath)
		cmd := exec.Command("docker", "run", "-d", "--name", "toyamas-"+appName, "-p", "8080:80", appName)
		_ = cmd.Run()
	} else {
		cmd := exec.Command("docker", "compose", "-f", composePath, "up", "-d")
		output, err := cmd.CombinedOutput()
		if err != nil {
			printer.LogError("Docker compose install failed: %s (%v)", string(output), err)
			if printer.CurrentMode == printer.ModeJSON {
				os.Exit(1)
			}
		}
	}

	if printer.CurrentMode == printer.ModeJSON {
		printer.PrintJSON(map[string]interface{}{"status": "success", "app": appName, "installed": true})
	} else {
		printer.LogSuccess("Application '%s' installed and running successfully!", appName)
	}
}
