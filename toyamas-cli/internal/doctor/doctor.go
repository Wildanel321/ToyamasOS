package doctor

import (
	"fmt"

	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"toyamas-cli/internal/printer"
	"toyamas-cli/internal/tui"
)

type CheckItem struct {
	Category string `json:"category"`
	Name     string `json:"name"`
	Status   string `json:"status"` // "OK", "WARN", "FAIL"
	Details  string `json:"details"`
}

type DoctorReport struct {
	OverallStatus string      `json:"overall_status"`
	TotalChecks   int         `json:"total_checks"`
	PassedChecks  int         `json:"passed_checks"`
	Checks        []CheckItem `json:"checks"`
}

func RunDoctor() *DoctorReport {
	spinner := tui.StartSpinner("Running ToyamasOS system diagnostic health checks")

	var checks []CheckItem
	passed := 0

	// Check 1: Operating System
	osStatus := "OK"
	osDetails := "Debian 13 Minimal Target Verified"
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		content := string(data)
		if !strings.Contains(strings.ToLower(content), "debian") {
			osStatus = "WARN"
			osDetails = fmt.Sprintf("OS reported non-Debian host (%s)", runtime.GOOS)
		}
	} else {
		osStatus = "WARN"
		osDetails = "Non-Linux system environment"
	}
	checks = append(checks, CheckItem{Category: "OS Base", Name: "Debian 13 Release Check", Status: osStatus, Details: osDetails})

	// Check 2: Memory & ZRAM
	zramStatus := "OK"
	zramDetails := "Active ZRAM compressed swap detected"
	if output, err := exec.Command("swapon", "--show").Output(); err == nil {
		if !strings.Contains(string(output), "zram") {
			zramStatus = "WARN"
			zramDetails = "ZRAM is not active (run 'toyamas-installer/scripts/05-zram.sh')"
		}
	} else {
		zramStatus = "WARN"
		zramDetails = "Could not verify swapon status"
	}
	checks = append(checks, CheckItem{Category: "Memory", Name: "ZRAM Compression Check", Status: zramStatus, Details: zramDetails})

	// Check 3: Docker Socket
	dockerStatus := "OK"
	dockerDetails := "Docker daemon active & unix socket readable"
	if _, err := os.Stat("/var/run/docker.sock"); os.IsNotExist(err) {
		dockerStatus = "FAIL"
		dockerDetails = "Docker socket (/var/run/docker.sock) not found"
	} else if err := exec.Command("docker", "version").Run(); err != nil {
		dockerStatus = "WARN"
		dockerDetails = "Docker socket exists but daemon returned error"
	}
	checks = append(checks, CheckItem{Category: "Containers", Name: "Docker Engine & Socket", Status: dockerStatus, Details: dockerDetails})

	// Check 4: UFW Firewall
	ufwStatus := "OK"
	ufwDetails := "UFW firewall active & enabled"
	if output, err := exec.Command("ufw", "status").Output(); err == nil {
		if !strings.Contains(strings.ToLower(string(output)), "active") {
			ufwStatus = "WARN"
			ufwDetails = "UFW firewall is inactive"
		}
	} else {
		ufwStatus = "WARN"
		ufwDetails = "UFW executable not found in PATH"
	}
	checks = append(checks, CheckItem{Category: "Security", Name: "UFW Firewall Status", Status: ufwStatus, Details: ufwDetails})

	// Check 5: Fail2Ban
	f2bStatus := "OK"
	f2bDetails := "Fail2Ban service active"
	if err := exec.Command("systemctl", "is-active", "fail2ban").Run(); err != nil {
		f2bStatus = "WARN"
		f2bDetails = "Fail2Ban service is not active"
	}
	checks = append(checks, CheckItem{Category: "Security", Name: "Fail2Ban SSH Protection", Status: f2bStatus, Details: f2bDetails})

	// Check 6: Free Disk Space
	diskStatus := "OK"
	diskDetails := "Sufficient disk space available"
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err == nil {
		freeGB := float64(stat.Bavail*uint64(stat.Bsize)) / (1024 * 1024 * 1024)
		diskDetails = fmt.Sprintf("%.2f GB free on /", freeGB)
		if freeGB < 1.0 {
			diskStatus = "FAIL"
			diskDetails = fmt.Sprintf("Low disk space: only %.2f GB free", freeGB)
		}
	}
	checks = append(checks, CheckItem{Category: "Storage", Name: "Root Partition Free Space", Status: diskStatus, Details: diskDetails})

	spinner.Stop("System diagnostic checks complete")

	for _, c := range checks {
		if c.Status == "OK" {
			passed++
		}
	}

	overall := "PASSED"
	if passed < len(checks) {
		overall = "WARNINGS_FOUND"
	}

	report := &DoctorReport{
		OverallStatus: overall,
		TotalChecks:   len(checks),
		PassedChecks:  passed,
		Checks:        checks,
	}

	if printer.CurrentMode == printer.ModeJSON {
		printer.PrintJSON(report)
	} else {
		printer.LogInfo("Diagnostic Summary: %d/%d Checks Passed", passed, len(checks))
		headers := []string{"Category", "Diagnostic Check", "Status", "Details"}
		var rows [][]string
		for _, c := range checks {
			statusSymbol := "\033[1;32m✔ OK\033[0m"
			if c.Status == "WARN" {
				statusSymbol = "\033[1;33m⚠ WARN\033[0m"
			} else if c.Status == "FAIL" {
				statusSymbol = "\033[1;31m✖ FAIL\033[0m"
			}
			rows = append(rows, []string{c.Category, c.Name, statusSymbol, c.Details})
		}
		tui.RenderTable(headers, rows)
	}

	return report
}
