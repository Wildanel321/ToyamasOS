#!/usr/bin/env bash
# ==============================================================================
# ToyamasOS Installer Script: 06-netdata.sh
# Purpose: Install Netdata real-time infrastructure monitoring tool
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../install.sh" --source-libs 2>/dev/null || true

log_info "Step 6: Installing Netdata monitoring server..."

export DEBIAN_FRONTEND=noninteractive

if command -v netdata >/dev/null 2>&1; then
    log_info "Netdata is already installed."
else
    log_info "Downloading and executing Netdata kickstart script..."
    
    # Try installing netdata via apt first if available, or official kickstart script
    if apt-cache show netdata >/dev/null 2>&1; then
        apt-get install -y --no-install-recommends netdata || true
    fi

    if ! command -v netdata >/dev/null 2>&1; then
        curl -fsSL https://get.netdata.cloud/kickstart.sh | sh -s -- --dont-wait --disable-telemetry --no-updates || {
            log_warn "Netdata kickstart script encountered a non-fatal warning."
        }
    fi
fi

# Ensure netdata service is enabled
if systemctl list-unit-files | grep -q netdata; then
    systemctl enable netdata.service 2>/dev/null || true
    systemctl start netdata.service 2>/dev/null || true
fi

# Note for Firewall: Netdata defaults to port 19999
log_info "Netdata listening default port: 19999/tcp"
log_success "Netdata monitoring installed successfully."
