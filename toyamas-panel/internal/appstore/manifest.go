package appstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type EnvVarConfig struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Default     string `json:"default"`
	Type        string `json:"type,omitempty"`        // "text", "password", "number"
	Description string `json:"description,omitempty"`
}

type AppManifest struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	Category    string         `json:"category"`
	Icon        string         `json:"icon"`
	Description string         `json:"description"`
	DefaultPort int            `json:"default_port"`
	EnvVars     []EnvVarConfig `json:"env_vars"`
	Status      string         `json:"status"` // "not_installed", "running", "stopped"
	Dir         string         `json:"-"`
}

type Repository struct {
	AppsDir string
}

func NewRepository(appsDir string) *Repository {
	return &Repository{AppsDir: appsDir}
}

func (r *Repository) ListManifests() ([]AppManifest, error) {
	entries, err := os.ReadDir(r.AppsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read apps repository dir '%s': %w", r.AppsDir, err)
	}

	var manifests []AppManifest
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(r.AppsDir, entry.Name(), "manifest.json")
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue
		}

		data, err := os.ReadFile(manifestPath)
		if err != nil {
			continue
		}

		var m AppManifest
		if err := json.Unmarshal(data, &m); err != nil {
			continue
		}

		if m.ID == "" {
			m.ID = entry.Name()
		}
		m.Dir = filepath.Join(r.AppsDir, entry.Name())
		m.Status = "not_installed"

		manifests = append(manifests, m)
	}
	return manifests, nil
}

func (r *Repository) GetManifest(appID string) (*AppManifest, error) {
	manifestPath := filepath.Join(r.AppsDir, appID, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("app '%s' manifest not found: %w", appID, err)
	}

	var m AppManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest for app '%s': %w", appID, err)
	}

	if m.ID == "" {
		m.ID = appID
	}
	m.Dir = filepath.Join(r.AppsDir, appID)
	return &m, nil
}
