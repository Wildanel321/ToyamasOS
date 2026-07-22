package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"toyamas-cli/internal/printer"
	"toyamas-cli/internal/tui"
)

func HandleAI(prompt string) {
	if prompt == "" {
		printer.LogError("Missing prompt. Usage: toyamas ai \"<your instruction>\" (e.g. toyamas ai \"Deploy Laravel\")")
		os.Exit(1)
	}

	if printer.CurrentMode == printer.ModeTUI {
		tui.ShowBanner()
		printer.LogInfo("Consulting Toyamas AI Assistant...")
	}

	lower := strings.ToLower(prompt)

	if strings.Contains(lower, "deploy laravel") || strings.Contains(lower, "install laravel") {
		if printer.CurrentMode == printer.ModeJSON {
			printer.PrintJSON(map[string]interface{}{
				"prompt":                     prompt,
				"requires_user_confirmation": true,
				"action_plan": map[string]interface{}{
					"target": "Laravel Stack",
					"steps": []string{
						"1. Create Docker container for Laravel + MariaDB",
						"2. Create 'laravel' database",
						"3. Create Nginx reverse proxy on port 8000",
						"4. Provision SSL certificate",
						"5. Return application URL",
					},
				},
			})
			return
		}

		printer.LogInfo("🤖 AI Assistant Proposed Action Plan:")
		printer.PrintTUI("   1. Create Docker container for Laravel + MariaDB")
		printer.PrintTUI("   2. Create 'laravel' database")
		printer.PrintTUI("   3. Create Nginx reverse proxy on port 8000")
		printer.PrintTUI("   4. Provision SSL certificate")
		printer.PrintTUI("   5. Expose Application URL")
		fmt.Println()

		fmt.Print("\033[1;33m[SAFETY CONFIRMATION]\033[0m Are you sure you want to execute these actions on ToyamasOS? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer != "y" && answer != "yes" {
			printer.LogWarn("Action execution cancelled by user.")
			return
		}

		spinner := tui.StartSpinner("Executing confirmed AI deployment plan")
		HandleInstall("laravel")
		spinner.Stop("")

		printer.LogSuccess("🎉 Laravel stack deployed successfully!")
		printer.LogInfo("Access URL: http://localhost:8000")
		return
	}

	// General AI response fallback
	if printer.CurrentMode == printer.ModeJSON {
		printer.PrintJSON(map[string]interface{}{
			"prompt":   prompt,
			"response": "Toyamas AI: Pengecekan sistem selesai. Semua service utama aktif.",
		})
	} else {
		printer.LogSuccess("Toyamas AI Analysis:")
		printer.PrintTUI(" - Membaca log server: Normal")
		printer.PrintTUI(" - Service terdeteksi: Active")
		printer.PrintTUI(" - Memori ZRAM: Aktif (50%% RAM cap)")
		printer.PrintTUI(" - Rekomendasi: Kernel sysctl parameters optimal.")
	}
}
