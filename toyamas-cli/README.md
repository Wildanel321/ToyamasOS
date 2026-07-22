# Toyamas CLI (`toyamas`) 🛠️

[![Stack: Golang](https://img.shields.io/badge/Language-Golang%201.22-blue.svg)](https://go.dev)
[![Feature: TUI](https://img.shields.io/badge/Interface-Modern%20TUI-teal.svg)](#tui-terminal)
[![Feature: JSON Mode](https://img.shields.io/badge/Output-JSON%20Mode-yellow.svg)](#json-mode)
[![Feature: Plugin System](https://img.shields.io/badge/Plugins-Dynamic%20System-emerald.svg)](#plugin-system)

**`toyamas`** is an enterprise-grade command-line tool written in Golang for managing, monitoring, backing up, and diagnosing **ToyamasOS** (Debian 13 Minimal 1GB RAM / 2 vCPU VPS).

---

## ⚡ Key Commands

```bash
# Display system health, ZRAM, firewall, and Docker status
sudo toyamas status

# Output status in machine-readable JSON format
sudo toyamas status --json

# Run system package updates and check CLI self-updates
sudo toyamas update

# One-click install application stack (Nginx, Laravel, Minecraft, etc.)
sudo toyamas install nginx
sudo toyamas install laravel
sudo toyamas install minecraft

# Create compressed backup snapshot archive (/var/backups/toyamas/)
sudo toyamas backup

# Stream system & installer logs
sudo toyamas logs

# Run comprehensive system diagnostic health checks
sudo toyamas doctor

# List registered custom plugins
toyamas plugins
```

---

## 🎨 Features & Architecture

### 1. Modern TUI Terminal
Styled with ANSI colors, Unicode tables, loading spinners, and progress bars.

### 2. JSON Mode (`--json` / `-j`)
Every command supports `--json` returning structured JSON objects for CI/CD pipelines, shell scripts, and remote monitoring APIs.

### 3. Diagnostic Engine (`doctor`)
Performs 6 automated diagnostic checks:
- OS Base (Debian 13 minimal verification)
- Memory & ZRAM (RAM usage & compressed RAM swap state)
- Docker Engine & Socket connection
- Security Stack (UFW firewall status)
- Fail2Ban SSH protection jail
- Root partition free space

### 4. Dynamic Plugin System
Placing an executable script or binary in `/etc/toyamas/plugins` or `~/.toyamas/plugins` automatically registers it as a subcommand. Executing `toyamas custom-cmd` dynamically routes execution to the plugin!

---

## 🚀 Building & Installation

To compile the CLI binary:

```bash
cd toyamas-cli
go build -o toyamas ./cmd/toyamas
sudo mv toyamas /usr/local/bin/
```
