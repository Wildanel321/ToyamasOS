#!/usr/bin/env bash
# ==============================================================================
# ToyamasOS Installer Script: 02-essentials.sh
# Purpose: Install essential CLI utilities (htop, curl, git, etc.)
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../install.sh" --source-libs 2>/dev/null || true

log_info "Step 2: Installing core essential utilities..."

export DEBIAN_FRONTEND=noninteractive

PACKAGES=(
    curl
    git
    htop
    ca-certificates
    gnupg
    lsb-release
    apt-transport-https
    procps
    unzip
    tar
    sudo
)

log_info "Installing packages: ${PACKAGES[*]}"

apt-get install -y --no-install-recommends "${PACKAGES[@]}"

log_success "Essential utilities installed successfully."
