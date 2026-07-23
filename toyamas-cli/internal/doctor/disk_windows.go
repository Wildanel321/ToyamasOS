//go:build windows

package doctor

import "errors"

func getDiskFreeGB(path string) (float64, error) {
	return 0, errors.New("disk space metrics not supported on Windows")
}
