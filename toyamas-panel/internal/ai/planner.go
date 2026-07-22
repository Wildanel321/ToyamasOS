package ai

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

type ActionStep struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ActionPlan struct {
	Token       string       `json:"token"`
	ActionType  string       `json:"action_type"` // "deploy_app", "restart_service", "create_backup", "create_proxy", "optimize_ram"
	Target      string       `json:"target"`
	Steps       []ActionStep `json:"steps"`
	Params      map[string]string `json:"params"`
	CreatedAt   time.Time    `json:"created_at"`
}

type AIResponse struct {
	Message    string      `json:"message"`
	ActionPlan *ActionPlan `json:"action_plan,omitempty"`
	RequiresUserConfirmation bool `json:"requires_user_confirmation"`
}

type TokenStore struct {
	mu     sync.Mutex
	tokens map[string]*ActionPlan
}

var globalTokenStore = &TokenStore{
	tokens: make(map[string]*ActionPlan),
}

func generateToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func StoreToken(plan *ActionPlan) {
	globalTokenStore.mu.Lock()
	defer globalTokenStore.mu.Unlock()
	globalTokenStore.tokens[plan.Token] = plan
}

func ValidateAndPopToken(token string) (*ActionPlan, error) {
	globalTokenStore.mu.Lock()
	defer globalTokenStore.mu.Unlock()

	plan, exists := globalTokenStore.tokens[token]
	if !exists {
		return nil, fmt.Errorf("invalid or expired confirmation token")
	}

	delete(globalTokenStore.tokens, token)
	return plan, nil
}

type Planner struct {
	ollama *OllamaClient
}

func NewPlanner(ollama *OllamaClient) *Planner {
	return &Planner{ollama: ollama}
}

func (p *Planner) ProcessPrompt(prompt string, sysCtx *SystemContext) (*AIResponse, error) {
	lowerPrompt := strings.ToLower(prompt)

	// Intent 1: Deploy Laravel
	if strings.Contains(lowerPrompt, "deploy laravel") || strings.Contains(lowerPrompt, "install laravel") {
		token := generateToken()
		plan := &ActionPlan{
			Token:      token,
			ActionType: "deploy_app",
			Target:     "laravel",
			Params: map[string]string{
				"app_id": "laravel",
				"domain": "laravel.local",
			},
			Steps: []ActionStep{
				{Title: "Create Docker Container", Description: "Deploy Laravel PHP 8.2 FPM container with MariaDB database"},
				{Title: "Create Database", Description: "Provision 'laravel' database and user credentials"},
				{Title: "Create Nginx Reverse Proxy", Description: "Configure Nginx routing to port 8000"},
				{Title: "Provision SSL", Description: "Issue SSL certificate for secure HTTPS domain"},
				{Title: "Expose Endpoint", Description: "Return application access URL"},
			},
			CreatedAt: time.Now(),
		}
		StoreToken(plan)

		return &AIResponse{
			Message:    "Saya telah menyiapkan Rencana Aksi untuk men-deploy Laravel Stack. Mohon konfirmasi sebelum saya mengeksekusi aksi berikut pada sistem:",
			ActionPlan: plan,
			RequiresUserConfirmation: true,
		}, nil
	}

	// Intent 2: Create Backup
	if strings.Contains(lowerPrompt, "backup") {
		token := generateToken()
		plan := &ActionPlan{
			Token:      token,
			ActionType: "create_backup",
			Target:     "System & Database Snapshot",
			Params:     map[string]string{},
			Steps: []ActionStep{
				{Title: "Archive System Configurations", Description: "Compress /etc/toyamas, sysctl, and Fail2Ban configs"},
				{Title: "Backup SQLite Database", Description: "Backup Toyamas Panel user accounts and settings"},
				{Title: "Save Tarball Archive", Description: "Save compressed snapshot to /var/backups/toyamas/"},
			},
			CreatedAt: time.Now(),
		}
		StoreToken(plan)

		return &AIResponse{
			Message:    "Saya dapat membuatkan backup snapshot otomatis untuk sistem dan database ToyamasOS. Mohon konfirmasi untuk memulai pembuatan backup:",
			ActionPlan: plan,
			RequiresUserConfirmation: true,
		}, nil
	}

	// Intent 3: Service/RAM Diagnostics & Optimization Recommendation
	if strings.Contains(lowerPrompt, "ram") || strings.Contains(lowerPrompt, "service") || strings.Contains(lowerPrompt, "log") || strings.Contains(lowerPrompt, "diagnosa") {
		msg := fmt.Sprintf("📊 **Analisis Server ToyamasOS**:\n"+
			"- **RAM Usage**: %s\n"+
			"- **Status ZRAM**: %s\n"+
			"- **Kapasitas Disk**: %s\n", sysCtx.RAMUsage, sysCtx.ZRAMStatus, sysCtx.DiskFree)

		if len(sysCtx.DeadServices) > 0 {
			msg += fmt.Sprintf("- ⚠️ **Service Mati Detected**: %s\n", strings.Join(sysCtx.DeadServices, ", "))
		} else {
			msg += "- ✔️ **Status Service**: Semua service utama dalam keadaan normal/active.\n"
		}

		if len(sysCtx.Containers) > 0 {
			msg += fmt.Sprintf("- 🐳 **Docker Containers**: %s\n", strings.Join(sysCtx.Containers, ", "))
		}

		msg += "\n💡 **Rekomendasi Optimasi**:\n" +
			"1. High swappiness (vm.swappiness=100) sudah dikonfigurasi untuk memprioritaskan kompresi ZRAM.\n" +
			"2. Log rotation Docker dikunci maksimal 10MB per container untuk mencegah disk penuh.\n" +
			"3. Service bloatware (Bluetooth, CUPS, Avahi, ModemManager) telah dimatikan."

		return &AIResponse{
			Message:                  msg,
			RequiresUserConfirmation: false,
		}, nil
	}

	// Fallback to Ollama or standard response
	if p.ollama.IsAvailable() {
		sysPrompt := "Anda adalah AI Assistant cerdas untuk ToyamasOS (Debian 13 Minimal, 1GB RAM VPS). Bantu pengguna mengelola server, Docker, dan Nginx."
		respText, err := p.ollama.Generate(sysPrompt, prompt)
		if err == nil && respText != "" {
			return &AIResponse{
				Message:                  respText,
				RequiresUserConfirmation: false,
			}, nil
		}
	}

	return &AIResponse{
		Message: fmt.Sprintf("Hello! Saya Toyamas AI Assistant. Saya siap membantu Anda mengelola server Debian 13 Minimal.\n\n"+
			"Cobalah bertanya:\n"+
			"• 'Deploy Laravel'\n"+
			"• 'Diagnosa RAM dan Service Mati'\n"+
			"• 'Buat Backup Otomatis'"),
		RequiresUserConfirmation: false,
	}, nil
}
