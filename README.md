# ToyamasOS

**ToyamasOS** is an open-source, high-performance minimal Linux distribution, web management dashboard, enterprise CLI suite, and custom bootable ISO builder designed for **Debian 13 (Trixie) Minimal**.

Optimized specifically for low-resource instances (**1 GB RAM & 2 CPUs**), ToyamasOS empowers self-hosters and developers to run Docker containers efficiently with zero system bloat.

---

## 🗺️ Complete 6-Tahap Pipeline

```
Tahap 1: Debian 13 Minimal Base
  └──> Tahap 2: Toyamas Installer (ZRAM, Bloatware Removal, Sysctl Tuning)
        └──> Tahap 3: Toyamas Panel (Web Dashboard & 1-Click App Store)
              └──> Tahap 4: Toyamas CLI (toyamas)
                    └──> Tahap 5: Toyamas AI Assistant (Ollama & Qwen3 Integration)
                          └──> Tahap 6: Custom ISO Distro (ToyamasOS-1.0.iso)
```

---

## 📦 Core Ecosystem Components

### 1. 🚀 [ToyamasOS Bootstrap Installer](./toyamas-installer)
Production-ready Bash installer suite that strips desktop bloat, activates **ZRAM compressed swap**, tunes kernel sysctl parameters, installs official Docker CE, configures UFW & Fail2Ban security, and sets up Netdata monitoring.

```bash
curl -fsSL https://raw.githubusercontent.com/Wildanel321/ToyamasOS/main/toyamas-installer/install.sh | sudo bash
```

---

### 2. 🎛️ [Toyamas Panel](./toyamas-panel)
Ultra-lightweight real-time web dashboard & **One-Click App Store** (**< 30 MB RAM footprint**) built with **Golang**, **SQLite**, **TailwindCSS**, and **WebSockets**. Includes 11 initial applications (Nginx, Apache, Redis, MariaDB, PostgreSQL, Laravel, NodeJS, Minecraft Paper, Nextcloud, Uptime Kuma, Ollama).

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

### 4. 🤖 [Toyamas AI Assistant](./toyamas-panel/internal/ai)
Intelligent local server assistant powered by **Ollama** and **Qwen3**. Reads server logs, diagnoses dead services, analyzes RAM usage, and proposes automated container/reverse proxy/backup action plans with **strict user confirmation safety protocol**.

---

### 5. 💿 [ToyamasOS Custom ISO Builder](./toyamas-iso)
Production ISO build pipeline using `live-build`, `debootstrap`, `squashfs-tools`, and `grub` to generate **`ToyamasOS-1.0.iso`** — a bootable Linux distribution installer for bare metal servers and cloud VPS virtual machines.

```bash
cd toyamas-iso
sudo ./build-iso.sh
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
├── toyamas-panel/              # Golang + SQLite + WebSockets Web Dashboard & App Store & AI Engine
│   ├── apps/                   # 11 App Store Docker Compose templates
│   ├── cmd/server/             # Main Go application entrypoint
│   ├── internal/               # AI Engine, AppStore, Auth, DB, Docker, Metrics, Services, WS Hub
│   ├── web/                    # Templates & static JS/CSS
│   ├── docs/                   # Docker deployment documentation
│   ├── Dockerfile              # Multi-stage CGO-free build
│   └── docker-compose.yml      # One-click Compose deployment
├── toyamas-cli/                # Enterprise Golang CLI Tool (toyamas)
│   ├── cmd/toyamas/            # Main CLI router entrypoint
│   ├── internal/               # AI CLI Handler, Commands, TUI, Printer, Doctor, Backup, Plugins
│   ├── plugins/                # Plugin script templates
│   └── README.md               # Toyamas CLI documentation
└── toyamas-iso/                # Custom ISO Distribution Builder
    ├── auto/config             # live-build recipe
    ├── config/                 # Bootloaders (GRUB/Isolinux), package manifests, chroot hooks
    ├── build-iso.sh            # Automated ISO compilation pipeline
    └── README.md               # ISO build documentation
```

---

## 📜 License

Licensed under the [MIT License](LICENSE).
