# Toyamas Panel 🎛️

[![Stack: Golang](https://img.shields.io/badge/Backend-Golang%201.22-blue.svg)](https://go.dev)
[![Database: SQLite](https://img.shields.io/badge/Database-SQLite-lightgrey.svg)](https://sqlite.org)
[![Frontend: TailwindCSS](https://img.shields.io/badge/Frontend-TailwindCSS-teal.svg)](https://tailwindcss.com)
[![App Store](https://img.shields.io/badge/App%20Store-10%20Apps-emerald.svg)](#-modular-one-click-app-store)

**Toyamas Panel** is an ultra-lightweight, real-time web dashboard and **modular App Store** for Linux server administration and Docker container management. Specially engineered for low-resource VPS environments (**1 GB RAM / 2 vCPU** running Debian 13 Minimal).

---

## ⚡ Performance Benchmark

- **Memory Usage**: `< 30 MB RAM` total footprint.
- **CPU Footprint**: Idle < 0.1% CPU.
- **Binary Build**: Multi-stage CGO-Free Go static executable (~15 MB image size).

---

## 🛍️ Modular One-Click App Store

Toyamas Panel features a zero-code modular App Store repository. Adding a new self-hosted app requires only creating a folder with a `manifest.json` and a `docker-compose.yml` template in `apps/`.

### 10 Initial Applications Supported:
1. 🌐 **Nginx**: Web server & reverse proxy.
2. 🪶 **Apache**: HTTP web server engine.
3. 🔴 **Redis**: In-memory key-value cache store with password authentication.
4. 🦭 **MariaDB**: Relational MySQL-compatible database.
5. 🐘 **PostgreSQL**: Advanced open-source object-relational database.
6. 🔴 **Laravel Stack**: PHP 8.2 FPM + Nginx + MariaDB fullstack environment.
7. 💚 **NodeJS Stack**: Node.js 20 LTS JavaScript runtime container.
8. ⛏️ **Minecraft Paper Server**: High-performance PaperMC Java Minecraft server.
9. ☁️ **Nextcloud**: Self-hosted personal cloud storage & file sync platform.
10. ⚡ **Uptime Kuma**: Real-time service monitoring & status pages.

Every application supports:
- ⚡ **One-Click Install**: Auto-generates configuration templates & executes `docker compose up -d`.
- 🔄 **One-Click Update**: Pulls latest images and updates container inline.
- 🗑️ **One-Click Uninstall**: Gracefully tears down containers and cleans resources.

---

## ✨ Features

- 🔐 **Session Authentication**: Secure login/logout system powered by SQLite session tokens and bcrypt password hashing.
- ⚡ **Realtime WebSockets**: Live broadcast streaming host metrics & Docker container updates every 1.5 seconds.
- 📊 **Resource Monitoring**: Live gauges & interactive Chart.js line graphs for CPU, RAM, Disk, and Network I/O.
- 🐳 **Docker Management**: List active/stopped containers with live state indicators, start, stop, and restart controls.
- ⚙️ **Linux Service Control**: Live state monitoring and single-click `systemctl restart` execution for key Linux services (`docker`, `ufw`, `fail2ban`, `netdata`, `toyamas-zram`, `ssh`, `nginx`).
- 🌙 **Dark Mode & Mobile Responsive**: Sleek glassmorphic dark theme styled with TailwindCSS.

---

## 📁 Directory Structure

```
toyamas-panel/
├── apps/                       # Modular App Store Repository
│   ├── nginx/                  # Nginx manifest & compose template
│   ├── apache/                 # Apache manifest & compose template
│   ├── redis/                  # Redis manifest & compose template
│   ├── mariadb/                # MariaDB manifest & compose template
│   ├── postgresql/             # PostgreSQL manifest & compose template
│   ├── laravel/                # Laravel stack manifest & compose template
│   ├── nodejs/                 # NodeJS stack manifest & compose template
│   ├── minecraft-paper/        # Minecraft PaperMC manifest & compose template
│   ├── nextcloud/              # Nextcloud manifest & compose template
│   └── uptime-kuma/            # Uptime Kuma manifest & compose template
├── cmd/
│   └── server/
│       └── main.go             # Main HTTP & WebSocket server entrypoint
├── internal/
│   ├── appstore/               # App Store manifest scanner & compose installer
│   ├── auth/                   # Session & bcrypt password authentication
│   ├── config/                 # Environment configuration loader
│   ├── db/                     # SQLite schema migration & connection
│   ├── docker/                 # Unix Domain Socket Docker API client
│   ├── metrics/                # Host CPU, RAM, Disk, Network reader
│   ├── services/               # Systemd Linux service restart manager
│   └── ws/                     # Realtime WebSocket Hub & broadcaster
├── web/
│   ├── static/
│   │   ├── css/style.css       # Custom scrollbars & styles
│   │   └── js/app.js           # WebSocket client, Chart.js, & App Store logic
│   └── templates/
│       ├── login.html          # Dark mode login page
│       └── dashboard.html      # Responsive dashboard & App Store UI
├── docs/
│   └── DOCKER_DEPLOYMENT.md    # Production deployment documentation
├── Dockerfile                  # Multi-stage lightweight CGO-free build
├── docker-compose.yml          # One-click Docker Compose configuration
├── go.mod                      # Go module definition
└── README.md                   # Toyamas Panel documentation
```

---

## 🚀 Quick Local Run

To run directly on a Go-enabled machine:

```bash
cd toyamas-panel
go run ./cmd/server
```

Open your browser to `http://localhost:8080`. Default login: `admin` / `toyamas123`.
