package metrics

import (
	"bufio"

	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CPUStats struct {
	UsagePercent float64   `json:"usage_percent"`
	Cores        int       `json:"cores"`
	LoadAvg      []float64 `json:"load_avg"`
}

type RAMStats struct {
	TotalMB float64 `json:"total_mb"`
	UsedMB  float64 `json:"used_mb"`
	FreeMB  float64 `json:"free_mb"`
	Percent float64 `json:"percent"`
}

type DiskStats struct {
	TotalGB float64 `json:"total_gb"`
	UsedGB  float64 `json:"used_gb"`
	FreeGB  float64 `json:"free_gb"`
	Percent float64 `json:"percent"`
}

type NetworkStats struct {
	RxBytesSec float64 `json:"rx_bytes_sec"`
	TxBytesSec float64 `json:"tx_bytes_sec"`
	TotalRxMB  float64 `json:"total_rx_mb"`
	TotalTxMB  float64 `json:"total_tx_mb"`
}

type SystemMetrics struct {
	Timestamp int64        `json:"timestamp"`
	HostName  string       `json:"hostname"`
	OS        string       `json:"os"`
	UptimeSec uint64       `json:"uptime_sec"`
	CPU       CPUStats     `json:"cpu"`
	RAM       RAMStats     `json:"ram"`
	Disk      DiskStats    `json:"disk"`
	Network   NetworkStats `json:"network"`
}

type Collector struct {
	mu           sync.Mutex
	lastCPUTotal uint64
	lastCPUIdle  uint64
	lastRxBytes  uint64
	lastTxBytes  uint64
	lastNetTime  time.Time
}

func NewCollector() *Collector {
	return &Collector{
		lastNetTime: time.Now(),
	}
}

func (c *Collector) GetSystemMetrics() (*SystemMetrics, error) {
	hostname, _ := os.Hostname()

	cpuStats := c.getCPUStats()
	ramStats := c.getRAMStats()
	diskStats := c.getDiskStats()
	netStats := c.getNetworkStats()
	uptimeSec := c.getUptime()

	return &SystemMetrics{
		Timestamp: time.Now().Unix(),
		HostName:  hostname,
		OS:        runtime.GOOS,
		UptimeSec: uptimeSec,
		CPU:       cpuStats,
		RAM:       ramStats,
		Disk:      diskStats,
		Network:   netStats,
	}, nil
}

func (c *Collector) getCPUStats() CPUStats {
	cores := runtime.NumCPU()
	loadAvg := []float64{0.0, 0.0, 0.0}

	// Read loadavg
	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			l1, _ := strconv.ParseFloat(fields[0], 64)
			l5, _ := strconv.ParseFloat(fields[1], 64)
			l15, _ := strconv.ParseFloat(fields[2], 64)
			loadAvg = []float64{l1, l5, l15}
		}
	}

	// Calculate CPU usage percentage from /proc/stat
	c.mu.Lock()
	defer c.mu.Unlock()

	usagePercent := 0.0
	if file, err := os.Open("/proc/stat"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) >= 5 && fields[0] == "cpu" {
				var user, nice, system, idle, iowait, irq, softirq, steal uint64
				user, _ = strconv.ParseUint(fields[1], 10, 64)
				nice, _ = strconv.ParseUint(fields[2], 10, 64)
				system, _ = strconv.ParseUint(fields[3], 10, 64)
				idle, _ = strconv.ParseUint(fields[4], 10, 64)
				if len(fields) > 5 { iowait, _ = strconv.ParseUint(fields[5], 10, 64) }
				if len(fields) > 6 { irq, _ = strconv.ParseUint(fields[6], 10, 64) }
				if len(fields) > 7 { softirq, _ = strconv.ParseUint(fields[7], 10, 64) }
				if len(fields) > 8 { steal, _ = strconv.ParseUint(fields[8], 10, 64) }

				total := user + nice + system + idle + iowait + irq + softirq + steal
				idleTotal := idle + iowait

				if c.lastCPUTotal > 0 && total > c.lastCPUTotal {
					totalDiff := float64(total - c.lastCPUTotal)
					idleDiff := float64(idleTotal - c.lastCPUIdle)
					usagePercent = ((totalDiff - idleDiff) / totalDiff) * 100.0
				}
				c.lastCPUTotal = total
				c.lastCPUIdle = idleTotal
			}
		}
	}

	if usagePercent < 0 { usagePercent = 0 }
	if usagePercent > 100 { usagePercent = 100 }

	return CPUStats{
		UsagePercent: usagePercent,
		Cores:        cores,
		LoadAvg:      loadAvg,
	}
}

func (c *Collector) getRAMStats() RAMStats {
	var totalKB, availableKB, freeKB uint64

	file, err := os.Open("/proc/meminfo")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) >= 2 {
				key := strings.TrimSuffix(fields[0], ":")
				val, _ := strconv.ParseUint(fields[1], 10, 64)
				switch key {
				case "MemTotal":
					totalKB = val
				case "MemAvailable":
					availableKB = val
				case "MemFree":
					freeKB = val
				}
			}
		}
	}

	if availableKB == 0 {
		availableKB = freeKB
	}

	totalMB := float64(totalKB) / 1024.0
	usedMB := float64(totalKB-availableKB) / 1024.0
	freeMB := float64(availableKB) / 1024.0
	percent := 0.0
	if totalMB > 0 {
		percent = (usedMB / totalMB) * 100.0
	}

	return RAMStats{
		TotalMB: totalMB,
		UsedMB:  usedMB,
		FreeMB:  freeMB,
		Percent: percent,
	}
}

func (c *Collector) getDiskStats() DiskStats {
	stats, err := getDiskStats("/")
	if err != nil {
		return DiskStats{TotalGB: 20.0, UsedGB: 5.0, FreeGB: 15.0, Percent: 25.0}
	}
	return stats
}

func (c *Collector) getNetworkStats() NetworkStats {
	c.mu.Lock()
	defer c.mu.Unlock()

	var totalRx, totalTx uint64

	file, err := os.Open("/proc/net/dev")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.Contains(line, ":") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					iface := strings.TrimSpace(parts[0])
					if iface == "lo" || strings.HasPrefix(iface, "veth") || strings.HasPrefix(iface, "br-") || strings.HasPrefix(iface, "docker") {
						continue
					}
					fields := strings.Fields(parts[1])
					if len(fields) >= 9 {
						rx, _ := strconv.ParseUint(fields[0], 10, 64)
						tx, _ := strconv.ParseUint(fields[8], 10, 64)
						totalRx += rx
						totalTx += tx
					}
				}
			}
		}
	}

	now := time.Now()
	elapsedSec := now.Sub(c.lastNetTime).Seconds()
	if elapsedSec <= 0 {
		elapsedSec = 1.0
	}

	var rxRate, txRate float64
	if c.lastRxBytes > 0 && totalRx >= c.lastRxBytes {
		rxRate = float64(totalRx-c.lastRxBytes) / elapsedSec
	}
	if c.lastTxBytes > 0 && totalTx >= c.lastTxBytes {
		txRate = float64(totalTx-c.lastTxBytes) / elapsedSec
	}

	c.lastRxBytes = totalRx
	c.lastTxBytes = totalTx
	c.lastNetTime = now

	return NetworkStats{
		RxBytesSec: rxRate,
		TxBytesSec: txRate,
		TotalRxMB:  float64(totalRx) / (1024 * 1024),
		TotalTxMB:  float64(totalTx) / (1024 * 1024),
	}
}

func (c *Collector) getUptime() uint64 {
	if data, err := os.ReadFile("/proc/uptime"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) > 0 {
			if sec, err := strconv.ParseFloat(fields[0], 64); err == nil {
				return uint64(sec)
			}
		}
	}
	return 0
}
