# ToyamasOS

**ToyamasOS** is an open-source, high-performance minimal Linux bootstrap distribution profile, web management dashboard, and enterprise CLI suite designed for **Debian 13 (Trixie) Minimal**.

Optimized specifically for low-resource cloud VPS instances (**1 GB RAM & 2 vCPUs**), ToyamasOS empowers self-hosters and developers to run Docker containers efficiently with zero system bloat.

---

## 📦 Core Ecosystem Components

### 1. 🚀 [ToyamasOS Bootstrap Installer](./toyamas-installer)
Production-ready Bash installer suite that strips desktop bloat, activates **ZRAM compressed swap**, tunes kernel sysctl parameters, installs official Docker CE, configures UFW & Fail2Ban security, and sets up Netdata monitoring.

```bash
curl -fsSL https://raw.githubusercontent.com/Wildanel321/ToyamasOS/main/toyamas-installer/install.sh | sudo bash
```

---

### 2. 🎛️ [Toyamas Panel](./toyamas-panel)
Ultra-lightweight real-time web dashboard & **One-Click App Store** (**< 30 MB RAM footprint**) built with **Golang**, **SQLite**, **TailwindCSS**, and **WebSockets**. Supports 10 initial applications (Nginx, Apache, Redis, MariaDB, PostgreSQL, Laravel, NodeJS, Minecraft Paper, Nextcloud, Uptime Kuma).

```bash
cd toyamas-panel
docker compose up -d --build
```

---

### 3. 🛠️ [Toyamas CLI](./toyamas-cli)
Enterprise command-line management utility written in Golang featuring a modern TUI interface, spinners, progress bars, machine-readable `--json` output mode, diagnostic doctor, automated backups, auto-updater, and dynamic plugin system.

```bash
sudo toyamas status
sudo toyamas doctor
sudo toyamas install nginx
sudo toyamas backup
```

---

## 📁 Repository Structure

```
ToyamasOS/
├── toyamas-installer/          # Minimal Debian 13 Bootstrap Installer
│   ├── install.sh              # Main CLI installer script
│   ├── configs/                # Sysctl, Fail2Ban, ZRAM configs
│   ├── scripts/                # Modular install steps (01 to 07)
│   ├── services/               # Systemd ZRAM service
│   └── docs/                   # Architecture & optimization docs
├── toyamas-panel/              # Golang + SQLite + WebSockets Web Dashboard & App Store
│   ├── apps/                   # 10 App Store Docker Compose templates
│   ├── cmd/server/             # Main Go application entrypoint
│   ├── internal/               # Auth, DB, Metrics, Docker, Services, AppStore, WS Hub
│   ├── web/                    # Templates & static JS/CSS
│   ├── docs/                   # Docker deployment documentation
│   ├── Dockerfile              # Multi-stage CGO-free build
│   └── docker-compose.yml      # One-click Compose deployment
└── toyamas-cli/                # Enterprise Golang CLI Tool (toyamas)
    ├── cmd/toyamas/            # Main CLI router entrypoint
    ├── internal/               # Commands, TUI, Printer, Doctor, Backup, Plugins, Updater
    ├── plugins/                # Plugin script templates
    └── README.md               # Toyamas CLI documentation
```

---

## 📜 License

Licensed under the [MIT License](LICENSE).
