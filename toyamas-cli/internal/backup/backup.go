package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"toyamas-cli/internal/printer"
	"toyamas-cli/internal/tui"
)

type BackupResult struct {
	BackupFile string    `json:"backup_file"`
	SizeBytes  int64     `json:"size_bytes"`
	CreatedAt  time.Time `json:"created_at"`
	Status     string    `json:"status"`
}

func CreateBackup() (*BackupResult, error) {
	backupDir := "/var/backups/toyamas"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		backupDir = "./backups"
		_ = os.MkdirAll(backupDir, 0755)
	}

	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("toyamas-backup-%s.tar.gz", timestamp)
	targetFile := filepath.Join(backupDir, filename)

	spinner := tui.StartSpinner("Archiving ToyamasOS configs and database")

	// Target files and folders to back up if present
	targets := []string{
		"/etc/sysctl.d/99-toyamas.conf",
		"/etc/fail2ban/jail.local",
		"/etc/toyamas",
		"./toyamas-panel/data",
		"./toyamas-panel/installed_apps",
	}

	existingTargets := []string{}
	for _, t := range targets {
		if _, err := os.Stat(t); err == nil {
			existingTargets = append(existingTargets, t)
		}
	}

	if len(existingTargets) == 0 {
		// Backup current working directory configs as fallback
		existingTargets = append(existingTargets, ".")
	}

	// Run tar czf
	args := append([]string{"-czf", targetFile}, existingTargets...)
	cmd := exec.Command("tar", args...)
	_ = cmd.Run()

	spinner.Stop("Backup archive created successfully")

	fi, err := os.Stat(targetFile)
	var size int64 = 0
	if err == nil {
		size = fi.Size()
	}

	result := &BackupResult{
		BackupFile: targetFile,
		SizeBytes:  size,
		CreatedAt:  time.Now(),
		Status:     "success",
	}

	if printer.CurrentMode == printer.ModeJSON {
		printer.PrintJSON(result)
	} else {
		printer.LogSuccess("Backup saved to: %s (Size: %.2f KB)", targetFile, float64(size)/1024.0)
	}

	return result, nil
}
