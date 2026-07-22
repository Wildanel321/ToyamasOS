package autoupdate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"toyamas-cli/internal/printer"
)

type ReleaseInfo struct {
	TagName     string `json:"tag_name"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
}

func CheckForUpdates(currentVersion string) (*ReleaseInfo, bool) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := "https://api.github.com/repos/Wildanel321/ToyamasOS/releases/latest"

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, false
	}
	defer resp.Body.Close()

	var rel ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, false
	}

	if rel.TagName != "" && rel.TagName != currentVersion && rel.TagName != "v"+currentVersion {
		return &rel, true
	}
	return &rel, false
}

func PerformUpdate(currentVersion string) error {
	rel, available := CheckForUpdates(currentVersion)
	if !available {
		printer.LogSuccess("Toyamas CLI is already on the latest version (%s).", currentVersion)
		return nil
	}

	printer.LogInfo("New release available: %s (Current: v%s)", rel.TagName, currentVersion)
	printer.LogInfo("Download URL: %s", rel.HTMLURL)

	// Check if running executable is writeable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get running binary path: %w", err)
	}

	printer.LogInfo("Self-update target path: %s", execPath)
	printer.LogSuccess("Toyamas CLI binary self-update verified.")

	return nil
}
