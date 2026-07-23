//go:build windows

package metrics

import "errors"

func getDiskStats(path string) (DiskStats, error) {
	return DiskStats{}, errors.New("disk space metrics not supported on Windows")
}
