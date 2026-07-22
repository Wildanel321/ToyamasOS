package appstore

import (
	"fmt"

	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"toyamas-panel/internal/docker"
)

type Installer struct {
	repo       *Repository
	installDir string
	dockerCli  *docker.Client
}

func NewInstaller(repo *Repository, installDir string, dockerCli *docker.Client) *Installer {
	if installDir == "" {
		installDir = "./installed_apps"
	}
	_ = os.MkdirAll(installDir, 0755)
	return &Installer{
		repo:       repo,
		installDir: installDir,
		dockerCli:  dockerCli,
	}
}

func (i *Installer) GetAppStatus(appID string) string {
	targetDir := filepath.Join(i.installDir, appID)
	composePath := filepath.Join(targetDir, "docker-compose.yml")
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return "not_installed"
	}

	// Query containers to check if app container is running
	containers, err := i.dockerCli.ListContainers()
	if err == nil {
		for _, c := range containers {
			for _, name := range c.Names {
				if strings.Contains(strings.ToLower(name), appID) {
					if c.State == "running" {
						return "running"
					}
					return "stopped"
				}
			}
		}
	}
	return "installed"
}

func (i *Installer) Install(appID string, envValues map[string]string) error {
	m, err := i.repo.GetManifest(appID)
	if err != nil {
		return err
	}

	srcCompose := filepath.Join(m.Dir, "docker-compose.yml")
	if _, err := os.Stat(srcCompose); os.IsNotExist(err) {
		return fmt.Errorf("missing docker-compose.yml template for app '%s'", appID)
	}

	targetAppDir := filepath.Join(i.installDir, appID)
	if err := os.MkdirAll(targetAppDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for app '%s': %w", appID, err)
	}

	// Copy docker-compose.yml
	composeData, err := os.ReadFile(srcCompose)
	if err != nil {
		return fmt.Errorf("failed to read compose template: %w", err)
	}
	destCompose := filepath.Join(targetAppDir, "docker-compose.yml")
	if err := os.WriteFile(destCompose, composeData, 0644); err != nil {
		return fmt.Errorf("failed to write compose file: %w", err)
	}

	// Write .env file
	envLines := []string{}
	for _, envConfig := range m.EnvVars {
		val := envConfig.Default
		if userVal, exists := envValues[envConfig.Name]; exists && userVal != "" {
			val = userVal
		}
		envLines = append(envLines, fmt.Sprintf("%s=%s", envConfig.Name, val))
	}

	envPath := filepath.Join(targetAppDir, ".env")
	if err := os.WriteFile(envPath, []byte(strings.Join(envLines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	// Run docker compose up -d
	cmd := exec.Command("docker", "compose", "-f", destCompose, "up", "-d")
	cmd.Dir = targetAppDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker compose up failed for %s: %s (%w)", appID, string(output), err)
	}

	return nil
}

func (i *Installer) Uninstall(appID string) error {
	targetAppDir := filepath.Join(i.installDir, appID)
	destCompose := filepath.Join(targetAppDir, "docker-compose.yml")

	if _, err := os.Stat(destCompose); err == nil {
		cmd := exec.Command("docker", "compose", "-f", destCompose, "down", "-v")
		cmd.Dir = targetAppDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Log warning but proceed with directory removal
			fmt.Printf("[APPSTORE WARNING] docker compose down warning for %s: %s\n", appID, string(output))
		}
	}

	if err := os.RemoveAll(targetAppDir); err != nil {
		return fmt.Errorf("failed to remove directory for app '%s': %w", appID, err)
	}

	return nil
}

func (i *Installer) Update(appID string) error {
	targetAppDir := filepath.Join(i.installDir, appID)
	destCompose := filepath.Join(targetAppDir, "docker-compose.yml")

	if _, err := os.Stat(destCompose); os.IsNotExist(err) {
		return fmt.Errorf("app '%s' is not installed", appID)
	}

	// docker compose pull
	cmdPull := exec.Command("docker", "compose", "-f", destCompose, "pull")
	cmdPull.Dir = targetAppDir
	_, _ = cmdPull.CombinedOutput()

	// docker compose up -d
	cmdUp := exec.Command("docker", "compose", "-f", destCompose, "up", "-d")
	cmdUp.Dir = targetAppDir
	output, err := cmdUp.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker compose update failed for %s: %s (%w)", appID, string(output), err)
	}

	return nil
}
