//go:build !windows

package metrics

import "syscall"

func getDiskStats(path string) (DiskStats, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return DiskStats{}, err
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free

	totalGB := float64(total) / (1024 * 1024 * 1024)
	usedGB := float64(used) / (1024 * 1024 * 1024)
	freeGB := float64(free) / (1024 * 1024 * 1024)
	percent := 0.0
	if totalGB > 0 {
		percent = (usedGB / totalGB) * 100.0
	}

	return DiskStats{
		TotalGB: totalGB,
		UsedGB:  usedGB,
		FreeGB:  freeGB,
		Percent: percent,
	}, nil
}
