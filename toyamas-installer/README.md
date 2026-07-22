# ToyamasOS Bootstrap Installer

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Target OS](https://img.shields.io/badge/OS-Debian%2013%20Minimal-red.svg)](https://www.debian.org)
[![Target Hardware](https://img.shields.io/badge/Spec-1GB%20RAM%20%7C%202%20vCPU-green.svg)](#optimizations)

**ToyamasOS** is an open-source, ultra-lightweight Linux distribution setup profile tailored for low-resource Virtual Private Servers (VPS). Built on **Debian 13 (Trixie) Minimal**, ToyamasOS strips out unnecessary desktop bloat, optimizes memory footprint via **ZRAM**, applies server-grade kernel tuning, and installs a production-ready **Docker + Security + Netdata** stack.

---

## 🎯 Features

- 🐧 **Debian 13 Minimal Base**: Headless, bloat-free foundation without any GUI/Desktop environments.
- 🐳 **Docker & Docker Compose Engine**: Official Docker repo integration with container log-rotation pre-configured.
- ⚡ **ZRAM Swap Memory Compression**: ZSTD-compressed RAM swap allocating 50% RAM space for maximum memory density on 1GB VPS.
- 🛡️ **Hardened Firewall & Security**: UFW (default deny input, SSH/HTTP/HTTPS allowed) and Fail2Ban active SSH jail.
- 📊 **Real-time Monitoring**: Netdata telemetry-free performance monitor.
- 🧰 **Essential Tooling**: Pre-installed `htop`, `curl`, `git`, `ca-certificates`, `gnupg`.
- ⚙️ **Kernel Sysctl Tuning**: Tuned `vm.swappiness=100`, socket backlog (`somaxconn`), TCP keepalives, and dirty page writebacks.
- 🚫 **Bloat Elimination**: Explicitly disables and masks `bluetooth`, `cups`, `avahi-daemon`, and `ModemManager`.

---

## 📁 Repository Structure

```
toyamas-installer/
├── install.sh                  # Main CLI installer script
├── configs/                    # Production configuration templates
│   ├── sysctl-toyamas.conf     # Optimized sysctl kernel tuning
│   ├── fail2ban-jail.local     # Fail2Ban SSH protection configuration
│   └── zram.conf               # ZRAM memory allocation config
├── scripts/                    # Modular installation steps
│   ├── 01-system-update.sh     # System apt update & upgrade
│   ├── 02-essentials.sh        # Core tools installation
│   ├── 03-docker.sh            # Official Docker CE setup
│   ├── 04-security.sh          # UFW firewall & Fail2Ban configuration
│   ├── 05-zram.sh              # ZRAM activation script
│   ├── 06-netdata.sh           # Netdata monitoring installation
│   └── 07-optimization.sh      # Service bloat disabling & kernel sysctl
├── services/                   # Systemd service unit files
│   └── toyamas-zram.service    # Systemd unit for ZRAM startup
├── docs/                       # Technical documentation
│   ├── ARCHITECTURE.md         # System design architecture
│   ├── OPTIMIZATIONS.md        # Deep dive into memory, network & kernel tweaks
│   └── USAGE.md                # Quickstart and administration guide
└── README.md                   # Project documentation overview
```

---

## 🚀 Quick Start (One-Line Installer)

Execute directly on a fresh Debian 13 minimal VPS:

```bash
curl -fsSL https://raw.githubusercontent.com/Wildanel321/ToyamasOS/main/toyamas-installer/install.sh | sudo bash
```

Or clone and run locally:

```bash
git clone https://github.com/Wildanel321/ToyamasOS.git
cd ToyamasOS/toyamas-installer
chmod +x install.sh
sudo ./install.sh
```

---

## 🛠️ Advanced CLI Usage

Run a specific step or option:

```bash
# Display help and usage
sudo ./install.sh --help

# Skip system apt update step
sudo ./install.sh --skip-update

# Run only step 5 (ZRAM Setup)
sudo ./install.sh --step=5
```

---

## 📄 Documentation

- [Architecture Overview](docs/ARCHITECTURE.md)
- [Optimizations & Sysctl Specification](docs/OPTIMIZATIONS.md)
- [Usage & Administration Guide](docs/USAGE.md)

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).
