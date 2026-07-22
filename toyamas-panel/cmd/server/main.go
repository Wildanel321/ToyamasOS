package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"toyamas-panel/internal/ai"
	"toyamas-panel/internal/appstore"
	"toyamas-panel/internal/auth"
	"toyamas-panel/internal/config"
	"toyamas-panel/internal/db"
	"toyamas-panel/internal/docker"
	"toyamas-panel/internal/metrics"
	"toyamas-panel/internal/services"
	"toyamas-panel/internal/ws"
)

type Server struct {
	cfg          *config.Config
	database     *db.DB
	collector    *metrics.Collector
	dockerCli    *docker.Client
	sysManager   *services.Manager
	wsHub        *ws.Hub
	appRepo      *appstore.Repository
	appInstaller *appstore.Installer
	aiOllama     *ai.OllamaClient
	aiGatherer   *ai.ContextGatherer
	aiPlanner    *ai.Planner
	aiExecutor   *ai.Executor
	templates    *template.Template
}

func main() {
	log.Println("[TOYAMAS PANEL] Starting Toyamas Panel Server v1.0.0...")

	cfg := config.LoadConfig()

	database, err := db.InitDB(cfg.DBPath, cfg.DefaultAdmin, cfg.DefaultPass)
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize database: %v", err)
	}
	defer database.Conn.Close()

	collector := metrics.NewCollector()
	dockerCli := docker.NewClient()
	sysManager := services.NewManager()
	wsHub := ws.NewHub(collector, dockerCli)

	appRepo := appstore.NewRepository("apps")
	appInstaller := appstore.NewInstaller(appRepo, "installed_apps", dockerCli)

	// AI Assistant Core Engine
	aiOllama := ai.NewOllamaClient()
	aiGatherer := ai.NewContextGatherer(collector, dockerCli)
	aiPlanner := ai.NewPlanner(aiOllama)
	aiExecutor := ai.NewExecutor(appInstaller)

	go wsHub.Run()

	tmpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		log.Printf("[WARNING] Template Glob parse warning: %v. Will fall back to inline/direct parse.", err)
	}

	srv := &Server{
		cfg:          cfg,
		database:     database,
		collector:    collector,
		dockerCli:    dockerCli,
		sysManager:   sysManager,
		wsHub:        wsHub,
		appRepo:      appRepo,
		appInstaller: appInstaller,
		aiOllama:     aiOllama,
		aiGatherer:   aiGatherer,
		aiPlanner:    aiPlanner,
		aiExecutor:   aiExecutor,
		templates:    tmpl,
	}

	mux := http.NewServeMux()

	// Static Files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Page Handlers
	mux.HandleFunc("/login", srv.handleLoginPage)
	mux.HandleFunc("/", srv.handleDashboardPage)

	// Auth API
	mux.HandleFunc("/api/login", srv.handleAPILogin)
	mux.HandleFunc("/api/logout", srv.handleAPILogout)
	mux.HandleFunc("/api/user", srv.handleAPIUser)

	// Metrics & Container API
	mux.HandleFunc("/api/metrics", srv.requireAuth(srv.handleAPIMetrics))
	mux.HandleFunc("/api/docker/containers", srv.requireAuth(srv.handleAPIDockerContainers))
	mux.HandleFunc("/api/docker/action", srv.requireAuth(srv.handleAPIDockerAction))

	// System Services API
	mux.HandleFunc("/api/services", srv.requireAuth(srv.handleAPIServices))
	mux.HandleFunc("/api/services/restart", srv.requireAuth(srv.handleAPIServiceRestart))

	// App Store API
	mux.HandleFunc("/api/apps", srv.requireAuth(srv.handleAPIAppsList))
	mux.HandleFunc("/api/apps/install", srv.requireAuth(srv.handleAPIAppInstall))
	mux.HandleFunc("/api/apps/uninstall", srv.requireAuth(srv.handleAPIAppUninstall))
	mux.HandleFunc("/api/apps/update", srv.requireAuth(srv.handleAPIAppUpdate))

	// AI Assistant API
	mux.HandleFunc("/api/ai/chat", srv.requireAuth(srv.handleAPIAIChat))
	mux.HandleFunc("/api/ai/confirm", srv.requireAuth(srv.handleAPIAIConfirm))
	mux.HandleFunc("/api/ai/status", srv.requireAuth(srv.handleAPIAIStatus))

	// Realtime WebSocket
	mux.HandleFunc("/ws/metrics", srv.handleWSMetrics)

	serverAddr := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	log.Printf("[TOYAMAS PANEL] Server running on http://%s", serverAddr)
	log.Printf("[TOYAMAS PANEL] Default Login -> User: %s | Password: %s", cfg.DefaultAdmin, cfg.DefaultPass)

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("[FATAL] Server error: %v", err)
	}
}

// Middleware: Require Auth
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("toyamas_session")
		if err != nil || cookie.Value == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		user, err := auth.ValidateSession(s.database, cookie.Value)
		if err != nil || user == nil {
			http.Error(w, `{"error":"session expired"}`, http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// Page: Login
func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseFiles(filepath.Join("web", "templates", "login.html"))
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	_ = t.Execute(w, nil)
}

// Page: Dashboard
func (s *Server) handleDashboardPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	cookie, err := r.Cookie("toyamas_session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := auth.ValidateSession(s.database, cookie.Value)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles(filepath.Join("web", "templates", "dashboard.html"))
	if err != nil {
		http.Error(w, "Dashboard Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = t.Execute(w, map[string]interface{}{
		"Username": user.Username,
	})
}

// API: Login
func (s *Server) handleAPILogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.jsonError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	user, err := auth.Authenticate(s.database, strings.TrimSpace(payload.Username), payload.Password)
	if err != nil {
		s.jsonError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateSession(s.database, user.ID)
	if err != nil {
		s.jsonError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "toyamas_session",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	s.database.LogAction(user.ID, "login", "user_auth")
	s.jsonResponse(w, map[string]interface{}{"success": true, "username": user.Username})
}

// API: Logout
func (s *Server) handleAPILogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("toyamas_session")
	if err == nil && cookie.Value != "" {
		_ = auth.DeleteSession(s.database, cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "toyamas_session",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	s.jsonResponse(w, map[string]interface{}{"success": true})
}

// API: Current User Profile
func (s *Server) handleAPIUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("toyamas_session")
	if err != nil || cookie.Value == "" {
		s.jsonError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := auth.ValidateSession(s.database, cookie.Value)
	if err != nil || user == nil {
		s.jsonError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	s.jsonResponse(w, user)
}

// API: System Metrics
func (s *Server) handleAPIMetrics(w http.ResponseWriter, r *http.Request) {
	m, err := s.collector.GetSystemMetrics()
	if err != nil {
		s.jsonError(w, "Failed to read system metrics", http.StatusInternalServerError)
		return
	}
	s.jsonResponse(w, m)
}

// API: Docker Containers List
func (s *Server) handleAPIDockerContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := s.dockerCli.ListContainers()
	if err != nil {
		s.jsonResponse(w, []docker.ContainerInfo{})
		return
	}
	s.jsonResponse(w, containers)
}

// API: Docker Container Action
func (s *Server) handleAPIDockerAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID     string `json:"id"`
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := s.dockerCli.ActionContainer(req.ID, req.Action); err != nil {
		s.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.jsonResponse(w, map[string]interface{}{"success": true, "id": req.ID, "action": req.Action})
}

// API: Linux System Services
func (s *Server) handleAPIServices(w http.ResponseWriter, r *http.Request) {
	svcs := s.sysManager.ListKeyServices()
	s.jsonResponse(w, svcs)
}

// API: Restart Linux System Service
func (s *Server) handleAPIServiceRestart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := s.sysManager.RestartService(req.Name); err != nil {
		s.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.jsonResponse(w, map[string]interface{}{"success": true, "service": req.Name})
}

// API: App Store Manifests List
func (s *Server) handleAPIAppsList(w http.ResponseWriter, r *http.Request) {
	manifests, err := s.appRepo.ListManifests()
	if err != nil {
		s.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range manifests {
		manifests[i].Status = s.appInstaller.GetAppStatus(manifests[i].ID)
	}

	s.jsonResponse(w, manifests)
}

// API: App Store One-Click Install
func (s *Server) handleAPIAppInstall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		AppID string            `json:"app_id"`
		Env   map[string]string `json:"env"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := s.appInstaller.Install(req.AppID, req.Env); err != nil {
		s.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.jsonResponse(w, map[string]interface{}{"success": true, "app_id": req.AppID, "message": "Installed successfully"})
}

// API: App Store One-Click Uninstall
func (s *Server) handleAPIAppUninstall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		AppID string `json:"app_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := s.appInstaller.Uninstall(req.AppID); err != nil {
		s.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.jsonResponse(w, map[string]interface{}{"success": true, "app_id": req.AppID, "message": "Uninstalled successfully"})
}

// API: App Store One-Click Update
func (s *Server) handleAPIAppUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		AppID string `json:"app_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := s.appInstaller.Update(req.AppID); err != nil {
		s.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.jsonResponse(w, map[string]interface{}{"success": true, "app_id": req.AppID, "message": "Updated successfully"})
}

// API: AI Assistant Chat
func (s *Server) handleAPIAIChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	sysCtx := s.aiGatherer.GatherContext()
	aiRes, err := s.aiPlanner.ProcessPrompt(req.Prompt, sysCtx)
	if err != nil {
		s.jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.jsonResponse(w, aiRes)
}

// API: AI Assistant Action Confirmation
func (s *Server) handleAPIAIConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	res, err := s.aiExecutor.ConfirmAndExecute(req.Token)
	if err != nil {
		s.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.jsonResponse(w, res)
}

// API: AI Assistant Status
func (s *Server) handleAPIAIStatus(w http.ResponseWriter, r *http.Request) {
	available := s.aiOllama.IsAvailable()
	s.jsonResponse(w, map[string]interface{}{
		"ollama_available": available,
		"model":            s.aiOllama.Model,
		"host":             s.aiOllama.BaseURL,
	})
}

// WebSocket Metrics Stream
func (s *Server) handleWSMetrics(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("toyamas_session")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Unauthorized WebSocket request", http.StatusUnauthorized)
		return
	}
	user, err := auth.ValidateSession(s.database, cookie.Value)
	if err != nil || user == nil {
		http.Error(w, "Session expired", http.StatusUnauthorized)
		return
	}

	s.wsHub.ServeWS(w, r)
}

// Helper: JSON Response
func (s *Server) jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

// Helper: JSON Error
func (s *Server) jsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
