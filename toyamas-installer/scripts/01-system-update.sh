#!/usr/bin/env bash
# ==============================================================================
# ToyamasOS Installer Script: 01-system-update.sh
# Purpose: Perform non-interactive system update & package upgrade on Debian 13
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../install.sh" --source-libs 2>/dev/null || true

log_info "Step 1: Updating Debian package repository and upgrading system..."

export DEBIAN_FRONTEND=noninteractive

# Update repository lists
apt-get update -y

# Upgrade installed packages to latest versions
apt-get dist-upgrade -y --no-install-recommends

# Clean up stale cache
apt-get autoremove -y
apt-get autoclean -y

log_success "System packages updated successfully."
