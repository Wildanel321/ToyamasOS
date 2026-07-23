package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"toyamas-cli/internal/printer"
	"toyamas-cli/internal/tui"
)

type StatusResult struct {
	HostName string `json:"hostname"`
	OS       string `json:"os"`
	Cores    int    `json:"cores"`
	ZRAM     string `json:"zram"`
	UFW      string `json:"ufw"`
	Fail2Ban string `json:"fail2ban"`
	Docker   string `json:"docker"`
	DiskFree string `json:"disk_free"`
}

func HandleStatus() {
	hostname, _ := os.Hostname()

	zram := "inactive"
	if out, err := exec.Command("swapon", "--show").Output(); err == nil && strings.Contains(string(out), "zram") {
		zram = "active (zstd)"
	}

	ufw := "inactive"
	if out, err := exec.Command("ufw", "status").Output(); err == nil && strings.Contains(strings.ToLower(string(out)), "active") {
		ufw = "active"
	}

	fail2ban := "inactive"
	if err := exec.Command("systemctl", "is-active", "fail2ban").Run(); err == nil {
		fail2ban = "active"
	}

	docker := "inactive"
	if err := exec.Command("docker", "version").Run(); err == nil {
		docker = "active"
	}

	diskFree := "unknown"
	if freeGB, err := getDiskFreeGB("/"); err == nil {
		diskFree = fmt.Sprintf("%.2f GB", freeGB)
	}

	res := StatusResult{
		HostName: hostname,
		OS:       runtime.GOOS + " (Debian 13 Target)",
		Cores:    runtime.NumCPU(),
		ZRAM:     zram,
		UFW:      ufw,
		Fail2Ban: fail2ban,
		Docker:   docker,
		DiskFree: diskFree,
	}

	if printer.CurrentMode == printer.ModeJSON {
		printer.PrintJSON(res)
		return
	}

	tui.ShowBanner()
	headers := []string{"Component", "Current Status", "Details"}
	rows := [][]string{
		{"Hostname", hostname, "System Node Name"},
		{"OS Base", res.OS, fmt.Sprintf("%d vCPUs Detected", res.Cores)},
		{"ZRAM Swap", res.ZRAM, "Compressed RAM Swap"},
		{"UFW Firewall", res.UFW, "Default Deny Inbound"},
		{"Fail2Ban", res.Fail2Ban, "SSH Intrusion Prevention Jail"},
		{"Docker Engine", res.Docker, "Container Daemon API Socket"},
		{"Root Disk Space", res.DiskFree, "Free Storage Capacity"},
	}
	tui.RenderTable(headers, rows)
}
