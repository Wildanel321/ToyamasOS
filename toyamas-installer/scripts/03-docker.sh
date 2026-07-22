#!/usr/bin/env bash
# ==============================================================================
# ToyamasOS Installer Script: 03-docker.sh
# Purpose: Install and configure Docker CE and Docker Compose Plugin
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../install.sh" --source-libs 2>/dev/null || true

log_info "Step 3: Installing Docker & Docker Compose..."

export DEBIAN_FRONTEND=noninteractive

# Create keyrings directory
mkdir -p /etc/apt/keyrings
chmod 0755 /etc/apt/keyrings

# Download official Docker GPG key
log_info "Fetching Docker GPG key..."
curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg --overwrite
chmod a+r /etc/apt/keyrings/docker.gpg

# Detect architecture and Debian codename
ARCH=$(dpkg --print-architecture)
CODENAME=$(lsb_release -cs 2>/dev/null || echo "trixie")

# Map Debian 13 testing/trixie to bookworm if trixie repo is not yet released
if [ "$CODENAME" = "trixie" ] || [ "$CODENAME" = "sid" ]; then
    DOCKER_CODENAME="bookworm"
else
    DOCKER_CODENAME="$CODENAME"
fi

# Add Docker repository
echo "deb [arch=${ARCH} signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian ${DOCKER_CODENAME} stable" \
    | tee /etc/apt/sources.list.d/docker.list > /dev/null

# Update repository & install Docker
apt-get update -y
apt-get install -y --no-install-recommends \
    docker-ce \
    docker-ce-cli \
    containerd.io \
    docker-buildx-plugin \
    docker-compose-plugin

# Configure Docker daemon optimizations for 1GB VPS (log rotation)
log_info "Configuring Docker daemon optimizations..."
mkdir -p /etc/docker
cat <<'EOF' > /etc/docker/daemon.json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "live-restore": true
}
EOF

# Enable & restart Docker service
systemctl daemon-reload
systemctl enable docker.service
systemctl restart docker.service

log_success "Docker Engine and Docker Compose installed successfully."
log_info "Docker version: $(docker --version)"
log_info "Docker Compose version: $(docker compose version)"
