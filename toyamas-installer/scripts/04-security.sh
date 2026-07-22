#!/usr/bin/env bash
# ==============================================================================
# ToyamasOS Installer Script: 04-security.sh
# Purpose: Configure UFW firewall and Fail2Ban intrusion prevention
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../install.sh" --source-libs 2>/dev/null || true

log_info "Step 4: Installing and configuring UFW firewall & Fail2Ban..."

export DEBIAN_FRONTEND=noninteractive

apt-get install -y --no-install-recommends ufw fail2ban

# Configure UFW rules
log_info "Configuring UFW firewall rules..."
ufw --force reset >/dev/null 2>&1 || true

ufw default deny incoming
ufw default allow outgoing

# Allow custom or default SSH port
SSH_PORT="${SSH_PORT:-22}"
log_info "Allowing SSH access on port ${SSH_PORT}/tcp..."
ufw allow "${SSH_PORT}/tcp" comment "SSH Port"

# Allow Web Server ports
log_info "Allowing HTTP (80) and HTTPS (443) traffic..."
ufw allow 80/tcp comment "HTTP Web Server"
ufw allow 443/tcp comment "HTTPS Web Server"

# Enable Firewall non-interactively
ufw --force enable

# Configure Fail2Ban
log_info "Configuring Fail2Ban SSH jail..."
CONFIG_JAIL="${SCRIPT_DIR}/../configs/fail2ban-jail.local"

if [ -f "$CONFIG_JAIL" ]; then
    cp "$CONFIG_JAIL" /etc/fail2ban/jail.local
else
    log_warn "Custom fail2ban-jail.local not found. Using default SSH jail template."
    cat <<'EOF' > /etc/fail2ban/jail.local
[DEFAULT]
bantime  = 1h
findtime = 10m
maxretry = 5
banaction = ufw

[sshd]
enabled = true
port    = ssh
maxretry = 3
EOF
fi

# Enable and restart Fail2Ban
systemctl enable fail2ban.service
systemctl restart fail2ban.service

log_success "Security stack (UFW & Fail2Ban) successfully active."
