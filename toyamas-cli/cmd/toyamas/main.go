package main

import (
	"fmt"
	"os"
	"strings"

	"toyamas-cli/internal/cli/commands"
	"toyamas-cli/internal/plugins"
	"toyamas-cli/internal/printer"
	"toyamas-cli/internal/tui"
)

const Version = "1.0.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		showHelp()
		os.Exit(0)
	}

	// Parse Global Flags
	cleanArgs := []string{}
	for _, arg := range args {
		if arg == "--json" || arg == "-j" {
			printer.SetJSONMode(true)
		} else {
			cleanArgs = append(cleanArgs, arg)
		}
	}

	if len(cleanArgs) == 0 {
		showHelp()
		os.Exit(0)
	}

	subCmd := cleanArgs[0]

	switch subCmd {
	case "status":
		commands.HandleStatus()

	case "update":
		commands.HandleUpdate(Version)

	case "install":
		appName := ""
		if len(cleanArgs) > 1 {
			appName = cleanArgs[1]
		}
		commands.HandleInstall(appName)

	case "backup":
		commands.HandleBackup()

	case "logs":
		commands.HandleLogs()

	case "doctor":
		commands.HandleDoctor()

	case "ai":
		prompt := ""
		if len(cleanArgs) > 1 {
			prompt = strings.Join(cleanArgs[1:], " ")
		}
		commands.HandleAI(prompt)

	case "plugins":
		mgr := plugins.NewManager()
		list := mgr.ListPlugins()
		if printer.CurrentMode == printer.ModeJSON {
			printer.PrintJSON(list)
		} else {
			tui.ShowBanner()
			printer.LogInfo("=== Toyamas Plugin Registry ===")
			if len(list) == 0 {
				printer.PrintTUI("No custom plugins found in /etc/toyamas/plugins or ./plugins.")
			} else {
				headers := []string{"Plugin Name", "Script Path", "Description"}
				var rows [][]string
				for _, p := range list {
					rows = append(rows, []string{p.Name, p.Path, p.Description})
				}
				tui.RenderTable(headers, rows)
			}
		}

	case "help", "-h", "--help":
		showHelp()

	case "version", "-v", "--version":
		if printer.CurrentMode == printer.ModeJSON {
			printer.PrintJSON(map[string]string{"version": Version})
		} else {
			fmt.Printf("Toyamas CLI v%s\n", Version)
		}

	default:
		// Attempt to execute dynamic plugin
		pluginMgr := plugins.NewManager()
		err := pluginMgr.ExecutePlugin(subCmd, cleanArgs[1:])
		if err != nil {
			printer.LogError("Unknown command or plugin '%s'. Run 'toyamas help' for usage.", subCmd)
			os.Exit(1)
		}
	}
}

func showHelp() {
	if printer.CurrentMode == printer.ModeJSON {
		printer.PrintJSON(map[string]interface{}{
			"usage": "toyamas [command] [options]",
			"commands": []string{
				"status", "update", "install <app>", "backup", "logs", "doctor", "ai <prompt>", "plugins", "version",
			},
		})
		return
	}

	tui.ShowBanner()
	fmt.Println(`
Usage: sudo toyamas [COMMAND] [FLAGS]

Commands:
  status           Show system health, ZRAM, firewall, and container overview
  update           Run system package update and check CLI self-updates
  install <app>    Install application stack (e.g. nginx, laravel, minecraft)
  backup           Create compressed snapshot archive of configs and databases
  logs             Stream installer and system logs
  doctor           Perform diagnostic health checks (OS, RAM, ZRAM, Docker, UFW)
  ai <prompt>      Ask AI Assistant (e.g. toyamas ai "Deploy Laravel")
  plugins          List active dynamic plugins in /etc/toyamas/plugins

Global Flags:
  --json, -j       Output results in structured machine-readable JSON format
  -h, --help       Display help information
  -v, --version    Display version information

Examples:
  sudo toyamas status
  sudo toyamas ai "Deploy Laravel"
  sudo toyamas ai "Diagnosa RAM"
  sudo toyamas doctor
  sudo toyamas backup
`)
}
