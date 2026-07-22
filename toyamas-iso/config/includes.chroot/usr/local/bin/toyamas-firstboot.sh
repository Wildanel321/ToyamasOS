#!/usr/bin/env bash
# ToyamasOS Firstboot Automated Installer Script
set -euo pipefail

LOG_FILE="/var/log/toyamas-firstboot.log"

echo "[TOYAMAS FIRSTBOOT] Initializing ToyamasOS 1.0 setup..." | tee -a "$LOG_FILE"

# Run ZRAM startup
/usr/local/bin/toyamas-zram-start || true

# Check if Docker is installed, start Docker service
if command -v docker >/dev/null 2>&1; then
    systemctl enable docker.service 2>/dev/null || true
    systemctl start docker.service 2>/dev/null || true
    echo "[TOYAMAS FIRSTBOOT] Docker service active." | tee -a "$LOG_FILE"
fi

# Disable firstboot service after execution
systemctl disable toyamas-firstboot.service 2>/dev/null || true

echo "[TOYAMAS FIRSTBOOT] Firstboot setup complete. Welcome to ToyamasOS!" | tee -a "$LOG_FILE"
