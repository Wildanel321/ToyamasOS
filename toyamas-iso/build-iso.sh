#!/usr/bin/env bash
# ==============================================================================
#  ToyamasOS 1.0 Custom ISO Builder Pipeline
#  Target: Generate ToyamasOS-1.0.iso (Debian 13 Minimal Live/Install ISO)
# ==============================================================================

set -euo pipefail

ISO_NAME="ToyamasOS-1.0.iso"
BUILD_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ANSI Colors
COLOR_RESET="\033[0m"
COLOR_INFO="\033[1;34m"
COLOR_SUCCESS="\033[1;32m"
COLOR_WARN="\033[1;33m"
COLOR_ERROR="\033[1;31m"

log_info()    { echo -e "${COLOR_INFO}[INFO] $*${COLOR_RESET}"; }
log_success() { echo -e "${COLOR_SUCCESS}[SUCCESS] $*${COLOR_RESET}"; }
log_warn()    { echo -e "${COLOR_WARN}[WARN] $*${COLOR_RESET}"; }
log_error()   { echo -e "${COLOR_ERROR}[ERROR] $*${COLOR_RESET}"; }

# Root Check
if [[ $EUID -ne 0 ]]; then
    log_error "ToyamasOS ISO builder must be run as root (or with sudo)."
    exit 1
fi

log_info "========================================================"
log_info "  ToyamasOS 1.0 Custom ISO Build Pipeline"
log_info "========================================================"

# Check Build Prerequisites
PREREQS=(lb debootstrap mksquashfs xorriso gpg curl isohybrid)
MISSING=()

for tool in "${PREREQS[@]}"; do
    if ! command -v "$tool" >/dev/null 2>&1; then
        MISSING+=("$tool")
    fi
done

if [[ ${#MISSING[@]} -gt 0 ]]; then
    log_warn "Missing required ISO build tools: ${MISSING[*]}"
    log_info "Installing prerequisites via apt..."
    apt-get update -y
    apt-get install -y --no-install-recommends \
        live-build \
        debootstrap \
        squashfs-tools \
        xorriso \
        grub-pc-bin \
        grub-efi-amd64-bin \
        debian-archive-keyring \
        debian-keyring \
        curl \
        gnupg \
        syslinux-utils
fi

# Ensure Debian archive keyring is up to date with Debian 12 (Bookworm) and Debian 13 (Trixie) signing keys.
# This is required if the build host is running Ubuntu or an older Debian version.
log_info "Ensuring Debian archive keyring is up to date..."
mkdir -p /usr/share/keyrings

if [[ ! -f /usr/share/keyrings/debian-archive-keyring.gpg ]]; then
    touch /usr/share/keyrings/debian-archive-keyring.gpg
fi

# Debian 13 (Trixie) Key ID: 762F67A0B2C39DE4
# Debian 12 (Bookworm) Key ID: 8783D481
DEBIAN_KEYS=(
    "762F67A0B2C39DE4:13"
    "8783D481:12"
)

for key_info in "${DEBIAN_KEYS[@]}"; do
    IFS=":" read -r key_id key_version <<< "$key_info"
    if ! gpg --no-default-keyring --keyring /usr/share/keyrings/debian-archive-keyring.gpg --list-keys "$key_id" >/dev/null 2>&1; then
        log_info "Importing Debian $key_version signing key ($key_id) into keyring..."
        tmp_key=$(mktemp)
        if curl -sSL -o "$tmp_key" "https://ftp-master.debian.org/keys/archive-key-${key_version}.asc" || \
           wget -qO "$tmp_key" "https://ftp-master.debian.org/keys/archive-key-${key_version}.asc"; then
            gpg --no-default-keyring --keyring /usr/share/keyrings/debian-archive-keyring.gpg --import "$tmp_key"
            log_success "Imported Debian $key_version signing key."
        else
            log_error "Failed to download Debian $key_version signing key."
        fi
        rm -f "$tmp_key"
    else
        log_info "Debian $key_version signing key ($key_id) already in keyring."
    fi
done

# Patch live-build's chroot_linux-image script to find Contents-*.gz under the main/ component.
# This fixes a 404 error when building Debian 12/13.
for script_path in /usr/lib/live/build/chroot_linux-image /usr/lib/live/build/lb_chroot_linux-image; do
    if [[ -f "$script_path" ]]; then
        # Revert any previous corrupted patching (e.g. if double main/ or {_AREA}/main/ was injected)
        if grep -q -F 'main/main/Contents-${LB_ARCHITECTURES}' "$script_path"; then
            log_info "Reverting double main/ in $(basename "$script_path")..."
            sed -i 's|main/main/Contents-${LB_ARCHITECTURES}|main/Contents-${LB_ARCHITECTURES}|g' "$script_path"
        fi
        if grep -q -F '${_AREA}/main/Contents-${LB_ARCHITECTURES}' "$script_path"; then
            log_info "Reverting area main/ insertion in $(basename "$script_path")..."
            sed -i 's|\${_AREA}/main/Contents-\${LB_ARCHITECTURES}|\${_AREA}/Contents-\${LB_ARCHITECTURES}|g' "$script_path"
        fi

        # Only patch if the script has the unpatched direct path: dists/${LB_DISTRIBUTION}/Contents-${LB_ARCHITECTURES}
        if grep -q -F 'dists/${LB_DISTRIBUTION}/Contents-${LB_ARCHITECTURES}' "$script_path"; then
            log_info "Patching live-build $(basename "$script_path") script for Debian 13 compatibility..."
            sed -i 's|dists/${LB_DISTRIBUTION}/Contents-${LB_ARCHITECTURES}|dists/${LB_DISTRIBUTION}/main/Contents-${LB_ARCHITECTURES}|g' "$script_path"
        fi
    fi
done

# Clean previous build artifacts
cd "$BUILD_DIR"
log_info "Cleaning previous live-build state and cache..."
lb clean --all >/dev/null 2>&1 || true
rm -rf cache .build config/bootstrap config/chroot config/common config/binary

# Execute live-build config
log_info "Configuring live-build recipe for Debian 13 Minimal..."
lb config \
    --mode debian \
    --system live \
    --distribution trixie \
    --architectures amd64 \
    --archive-areas "main contrib non-free non-free-firmware" \
    --parent-mirror-bootstrap "http://deb.debian.org/debian/" \
    --parent-mirror-chroot "http://deb.debian.org/debian/" \
    --mirror-bootstrap "http://deb.debian.org/debian/" \
    --mirror-chroot "http://deb.debian.org/debian/" \
    --debian-installer false \
    --binary-images iso-hybrid \
    --bootloader syslinux \
    --compression squashfs \
    --iso-application "ToyamasOS Minimal Server" \
    --iso-publisher "ToyamasOS Team <https://github.com/Wildanel321/ToyamasOS>" \
    --iso-volume "TOYAMASOS_1_0" \
    --security false \
    --apt-indices false \
    --initsystem systemd \
    --memtest none

# Inject debootstrap options directly to the config file to bypass command-line parsing bugs in Ubuntu's live-build
echo 'LB_DEBOOTSTRAP_OPTIONS="--include=coreutils,usr-is-merged,systemd --no-check-gpg"' >> config/bootstrap

# Execute live-build ISO compilation
log_info "Building ToyamasOS rootfs, squashfs, and hybrid bootloader..."
lb build

# Find and rename ISO
RAW_ISO=$(ls live-image-amd64.hybrid.iso 2>/dev/null || ls *.iso 2>/dev/null | head -n 1)

if [[ -f "$RAW_ISO" ]]; then
    mv "$RAW_ISO" "$ISO_NAME"
    log_info "Generating SHA256 checksum..."
    sha256sum "$ISO_NAME" | tee "${ISO_NAME}.sha256"

    log_info "========================================================"
    log_success "ToyamasOS ISO Build Complete!"
    log_success "Output Image: ${BUILD_DIR}/${ISO_NAME}"
    log_success "Checksum:     ${BUILD_DIR}/${ISO_NAME}.sha256"
    log_info "========================================================"
    log_info "Flashing instructions for USB drive:"
    log_info "  sudo dd if=${ISO_NAME} of=/dev/sdX bs=4M status=progress"
else
    log_error "ISO compilation failed. No output image generated."
    exit 1
fi
