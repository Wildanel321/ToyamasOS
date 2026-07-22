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
PREREQS=(lb debootstrap mksquashfs xorriso)
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
        grub-efi-amd64-bin
fi

# Clean previous build artifacts
cd "$BUILD_DIR"
log_info "Cleaning previous live-build state..."
lb clean --all >/dev/null 2>&1 || true

# Execute live-build config
log_info "Configuring live-build recipe for Debian 13 Minimal..."
if [ -f "./auto/config" ]; then
    chmod +x ./auto/config
    ./auto/config
else
    lb config --distribution trixie --architectures amd64 --binary-images iso-hybrid
fi

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
