#!/usr/bin/env bash
# ==============================================================================
# ToyamasOS Installer Script: 07-optimization.sh
# Purpose: Disable unnecessary background services & apply sysctl kernel tuning
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../install.sh" --source-libs 2>/dev/null || true

log_info "Step 7: Applying ToyamasOS system optimizations..."

# List of desktop/unused server services to stop, disable, and mask
SERVICES_TO_DISABLE=(
    bluetooth.service
    cups.service
    cups-browsed.service
    avahi-daemon.service
    avahi-daemon.socket
    ModemManager.service
)

log_info "Disabling bloatware services for minimal headless server..."

for svc in "${SERVICES_TO_DISABLE[@]}"; do
    if systemctl list-unit-files | grep -q "^${svc}"; then
        log_info "Disabling and masking service: ${svc}"
        systemctl stop "${svc}" 2>/dev/null || true
        systemctl disable "${svc}" 2>/dev/null || true
        systemctl mask "${svc}" 2>/dev/null || true
    else
        log_info "Service ${svc} is not installed (skipping)."
    fi
done

# Apply custom sysctl kernel profile
SYSCTL_SRC="${SCRIPT_DIR}/../configs/sysctl-toyamas.conf"
SYSCTL_DEST="/etc/sysctl.d/99-toyamas.conf"

if [ -f "$SYSCTL_SRC" ]; then
    log_info "Installing ToyamasOS sysctl kernel parameters to ${SYSCTL_DEST}..."
    cp "$SYSCTL_SRC" "$SYSCTL_DEST"
    sysctl --system >/dev/null 2>&1 || sysctl -p "$SYSCTL_DEST"
    log_success "Kernel parameters applied successfully."
else
    log_warn "sysctl-toyamas.conf not found. Applying inline kernel tuning..."
    cat <<'EOF' > "$SYSCTL_DEST"
vm.swappiness = 100
vm.vfs_cache_pressure = 50
vm.dirty_background_ratio = 5
vm.dirty_ratio = 10
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 2048
net.ipv4.tcp_fin_timeout = 15
net.ipv4.tcp_keepalive_time = 300
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_syncookies = 1
EOF
    sysctl --system >/dev/null 2>&1 || sysctl -p "$SYSCTL_DEST"
fi

log_success "System bloat services disabled and kernel tuning completed."
