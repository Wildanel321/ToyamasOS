#!/usr/bin/env bash
# ==============================================================================
# ToyamasOS Installer Script: 05-zram.sh
# Purpose: Activate ZRAM compressed swap for low-memory optimization (1GB RAM)
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../install.sh" --source-libs 2>/dev/null || true

log_info "Step 5: Activating ZRAM compressed memory swap..."

export DEBIAN_FRONTEND=noninteractive

# Ensure zram kernel module is loaded
modprobe zram num_devices=1 2>/dev/null || true

# Check if zram kernel module exists
if ! lsmod | grep -q zram; then
    log_warn "zram module is not loaded. Attempting to load..."
    modprobe zram || log_warn "zram kernel module not available on host kernel."
fi

# Write helper script to manage ZRAM device
cat <<'EOF' > /usr/local/bin/toyamas-zram-start
#!/usr/bin/env bash
set -euo pipefail

# Read config or use defaults
CONFIG_FILE="/etc/toyamas/zram.conf"
COMP_ALGORITHM="zstd"
ZRAM_PERCENT="50"
ZRAM_PRIORITY="100"

if [ -f "$CONFIG_FILE" ]; then
    # shellcheck disable=SC1090
    source "$CONFIG_FILE"
fi

# Calculate memory size (50% of total RAM)
TOTAL_KB=$(grep MemTotal /proc/meminfo | awk '{print $2}')
ZRAM_KB=$(( TOTAL_KB * ZRAM_PERCENT / 100 ))
ZRAM_BYTES=$(( ZRAM_KB * 1024 ))

# Reset existing zram0 if active
if swapon --show | grep -q "/dev/zram0"; then
    swapoff /dev/zram0 2>/dev/null || true
fi
if [ -b /dev/zram0 ]; then
    zramctl --reset /dev/zram0 2>/dev/null || true
fi

# Find available algorithm (fallback to lz4 if zstd not supported)
if ! grep -q "$COMP_ALGORITHM" /sys/block/zram0/comp_algorithm 2>/dev/null; then
    COMP_ALGORITHM="lz4"
fi

# Configure ZRAM
zramctl /dev/zram0 --algorithm "$COMP_ALGORITHM" --size "${ZRAM_BYTES}"
mkswap /dev/zram0 >/dev/null
swapon -p "$ZRAM_PRIORITY" /dev/zram0

echo "ToyamasOS ZRAM initialized: /dev/zram0 (${ZRAM_KB} KB, algorithm: ${COMP_ALGORITHM})"
EOF

cat <<'EOF' > /usr/local/bin/toyamas-zram-stop
#!/usr/bin/env bash
set -euo pipefail

if swapon --show | grep -q "/dev/zram0"; then
    swapoff /dev/zram0 2>/dev/null || true
fi
if [ -b /dev/zram0 ]; then
    zramctl --reset /dev/zram0 2>/dev/null || true
fi

echo "ToyamasOS ZRAM stopped."
EOF

chmod +x /usr/local/bin/toyamas-zram-start
chmod +x /usr/local/bin/toyamas-zram-stop

# Store ZRAM configuration file
mkdir -p /etc/toyamas
if [ -f "${SCRIPT_DIR}/../configs/zram.conf" ]; then
    cp "${SCRIPT_DIR}/../configs/zram.conf" /etc/toyamas/zram.conf
fi

# Install systemd service
if [ -f "${SCRIPT_DIR}/../services/toyamas-zram.service" ]; then
    cp "${SCRIPT_DIR}/../services/toyamas-zram.service" /etc/systemd/system/toyamas-zram.service
fi

systemctl daemon-reload
systemctl enable toyamas-zram.service
systemctl restart toyamas-zram.service || /usr/local/bin/toyamas-zram-start

log_success "ZRAM compressed swap enabled successfully."
log_info "Active swap devices:"
swapon --show || true
