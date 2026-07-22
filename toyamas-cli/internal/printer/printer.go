package printer

import (
	"encoding/json"
	"fmt"
	"os"
)

type OutputMode int

const (
	ModeTUI OutputMode = iota
	ModeJSON
)

var (
	CurrentMode = ModeTUI
	IsVerbose   = false
)

// SetJSONMode enables raw JSON output mode
func SetJSONMode(enabled bool) {
	if enabled {
		CurrentMode = ModeJSON
	} else {
		CurrentMode = ModeTUI
	}
}

// PrintJSON outputs formatted JSON object and exits clean
func PrintJSON(data interface{}) {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf(`{"error": "failed to marshal json: %v"}`+"\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

// PrintTUI outputs text if in TUI mode
func PrintTUI(format string, args ...interface{}) {
	if CurrentMode == ModeTUI {
		fmt.Printf(format+"\n", args...)
	}
}

// LogInfo outputs info message
func LogInfo(format string, args ...interface{}) {
	if CurrentMode == ModeTUI {
		fmt.Printf("\033[1;34m[INFO]\033[0m "+format+"\n", args...)
	}
}

// LogSuccess outputs success message
func LogSuccess(format string, args ...interface{}) {
	if CurrentMode == ModeTUI {
		fmt.Printf("\033[1;32m[SUCCESS]\033[0m "+format+"\n", args...)
	}
}

// LogWarn outputs warning message
func LogWarn(format string, args ...interface{}) {
	if CurrentMode == ModeTUI {
		fmt.Printf("\033[1;33m[WARN]\033[0m "+format+"\n", args...)
	}
}

// LogError outputs error message
func LogError(format string, args ...interface{}) {
	if CurrentMode == ModeJSON {
		PrintJSON(map[string]string{"error": fmt.Sprintf(format, args...)})
	} else {
		fmt.Printf("\033[1;31m[ERROR]\033[0m "+format+"\n", args...)
	}
}
