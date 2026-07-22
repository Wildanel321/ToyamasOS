package services

import (
	"fmt"
	"os/exec"
	"strings"
)

type ServiceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Active string `json:"active"`
}

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) ListKeyServices() []ServiceInfo {
	targetServices := []string{
		"docker",
		"ufw",
		"fail2ban",
		"netdata",
		"toyamas-zram",
		"ssh",
		"sshd",
		"nginx",
	}

	var results []ServiceInfo
	for _, name := range targetServices {
		cmd := exec.Command("systemctl", "is-active", name)
		output, err := cmd.Output()
		activeState := strings.TrimSpace(string(output))
		if err != nil && activeState == "" {
			activeState = "inactive"
		}

		// Only include services that exist or are active
		if activeState != "unknown" {
			results = append(results, ServiceInfo{
				Name:   name,
				Status: activeState,
				Active: activeState,
			})
		}
	}
	return results
}

func (m *Manager) RestartService(serviceName string) error {
	// Sanitize service name to prevent command injection
	allowedServices := map[string]bool{
		"docker":       true,
		"ufw":          true,
		"fail2ban":     true,
		"netdata":      true,
		"toyamas-zram": true,
		"ssh":          true,
		"sshd":         true,
		"nginx":        true,
	}

	if !allowedServices[serviceName] {
		return fmt.Errorf("service '%s' is not in the allowed restart list", serviceName)
	}

	cmd := exec.Command("systemctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart service %s: %s (%w)", serviceName, string(output), err)
	}

	return nil
}
