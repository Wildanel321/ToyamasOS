package ai

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"toyamas-panel/internal/docker"
	"toyamas-panel/internal/metrics"
)

type SystemContext struct {
	RAMUsage      string   `json:"ram_usage"`
	ZRAMStatus    string   `json:"zram_status"`
	DeadServices  []string `json:"dead_services"`
	RecentLogs    []string `json:"recent_logs"`
	Containers    []string `json:"containers"`
	DiskFree      string   `json:"disk_free"`
}

type ContextGatherer struct {
	collector *metrics.Collector
	dockerCli *docker.Client
}

func NewContextGatherer(collector *metrics.Collector, dockerCli *docker.Client) *ContextGatherer {
	return &ContextGatherer{
		collector: collector,
		dockerCli: dockerCli,
	}
}

func (cg *ContextGatherer) GatherContext() *SystemContext {
	ctx := &SystemContext{
		RAMUsage:     "Unknown",
		ZRAMStatus:   "Inactive",
		DeadServices: []string{},
		RecentLogs:   []string{},
		Containers:   []string{},
		DiskFree:     "Unknown",
	}

	// Read Metrics
	if m, err := cg.collector.GetSystemMetrics(); err == nil {
		ctx.RAMUsage = fmt.Sprintf("%.0f MB / %.0f MB (%.1f%%)", m.RAM.UsedMB, m.RAM.TotalMB, m.RAM.Percent)
		ctx.DiskFree = fmt.Sprintf("%.1f GB free", m.Disk.FreeGB)
	}

	// Check ZRAM
	if out, err := exec.Command("swapon", "--show").Output(); err == nil {
		if strings.Contains(string(out), "zram") {
			ctx.ZRAMStatus = "Active (ZSTD compressed swap)"
		}
	}

	// Check Dead/Failed Services
	if out, err := exec.Command("systemctl", "list-units", "--state=failed", "--no-legend").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				fields := strings.Fields(line)
				if len(fields) > 0 {
					ctx.DeadServices = append(ctx.DeadServices, fields[0])
				}
			}
		}
	}

	// Read Docker Containers
	if containers, err := cg.dockerCli.ListContainers(); err == nil {
		for _, c := range containers {
			name := c.ID[:12]
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}
			ctx.Containers = append(ctx.Containers, fmt.Sprintf("%s (%s: %s)", name, c.Image, c.State))
		}
	}

	// Read Logs
	logFile := "/var/log/toyamas-installer.log"
	if data, err := os.ReadFile(logFile); err == nil {
		lines := strings.Split(string(data), "\n")
		start := len(lines) - 10
		if start < 0 {
			start = 0
		}
		for i := start; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) != "" {
				ctx.RecentLogs = append(ctx.RecentLogs, lines[i])
			}
		}
	}

	return ctx
}
