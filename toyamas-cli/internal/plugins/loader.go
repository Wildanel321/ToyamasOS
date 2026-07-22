package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"toyamas-cli/internal/printer"
)

type PluginInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type Manager struct {
	SearchDirs []string
}

func NewManager() *Manager {
	home, _ := os.UserHomeDir()
	dirs := []string{
		"/etc/toyamas/plugins",
		filepath.Join(home, ".toyamas", "plugins"),
		"./plugins",
	}
	return &Manager{SearchDirs: dirs}
}

func (m *Manager) ListPlugins() []PluginInfo {
	var results []PluginInfo
	seen := make(map[string]bool)

	for _, dir := range m.SearchDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			nameNoExt := strings.TrimSuffix(name, filepath.Ext(name))
			if seen[nameNoExt] {
				continue
			}

			fullPath := filepath.Join(dir, name)
			fi, err := entry.Info()
			if err != nil {
				continue
			}

			// Check executable permissions
			if fi.Mode()&0111 != 0 || strings.HasSuffix(name, ".sh") {
				seen[nameNoExt] = true
				results = append(results, PluginInfo{
					Name:        nameNoExt,
					Path:        fullPath,
					Description: fmt.Sprintf("Custom plugin script (%s)", filepath.Base(dir)),
				})
			}
		}
	}
	return results
}

func (m *Manager) ExecutePlugin(name string, args []string) error {
	plugins := m.ListPlugins()
	var target *PluginInfo
	for _, p := range plugins {
		if strings.EqualFold(p.Name, name) {
			target = &p
			break
		}
	}

	if target == nil {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	printer.LogInfo("Executing Toyamas plugin: %s (%s)", target.Name, target.Path)

	cmd := exec.Command(target.Path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("plugin execution failed: %w", err)
	}
	return nil
}
