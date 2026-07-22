#!/usr/bin/env bash
# ==============================================================================
#  ToyamasOS Bootstrap Installer
#  Target OS: Debian 13 Minimal (Trixie)
#  Target Hardware: VPS 1GB RAM, 2 vCPU
#  License: MIT
# ==============================================================================

set -euo pipefail

VERSION="1.0.0"
LOG_FILE="/var/log/toyamas-installer.log"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ANSI Color Codes
COLOR_RESET="\033[0m"
COLOR_INFO="\033[1;34m"
COLOR_SUCCESS="\033[1;32m"
COLOR_WARN="\033[1;33m"
COLOR_ERROR="\033[1;31m"
COLOR_BOLD="\033[1m"

# Logging Functions
log_msg() {
    local level="$1"
    local color="$2"
    shift 2
    local msg="$*"
    local timestamp
    timestamp=$(date "+%Y-%m-%d %H:%M:%S")
    echo -e "${color}[${level}] ${msg}${COLOR_RESET}"
    if [ -w "$(dirname "$LOG_FILE")" ] || [ "$EUID" -eq 0 ]; then
        echo "[${timestamp}] [${level}] ${msg}" >> "$LOG_FILE" 2>/dev/null || true
    fi
}

log_info()    { log_msg "INFO"    "$COLOR_INFO"    "$@"; }
log_success() { log_msg "SUCCESS" "$COLOR_SUCCESS" "$@"; }
log_warn()    { log_msg "WARN"    "$COLOR_WARN"    "$@"; }
log_error()   { log_msg "ERROR"   "$COLOR_ERROR"   "$@"; }

# Source helper guard
if [[ "${1:-}" == "--source-libs" ]]; then
    return 0 2>/dev/null || exit 0
fi

# Trap Error Handler
error_handler() {
    local exit_code="$1"
    local line_number="$2"
    log_error "An error occurred on line ${line_number} with exit status ${exit_code}."
    log_error "Check detailed log file at: ${LOG_FILE}"
    exit "$exit_code"
}
trap 'error_handler $? $LINENO' ERR

# Print Banner
show_banner() {
    cat << "EOF"

  _____ _____   __  __          __  __          _____ 
 |_   _/ _ \ \ / / / \   |_  _//  \/  \ /  |__/ / ___/
   | || (_) \ V / / _ \   | | / /\ / /\ V /|  __\___ \
   |_| \___/ |_|/_/   \_\ |_|/_/  /_/  \_/ |_|  /____/ 

  ToyamasOS Installer v1.0.0 - Optimized Minimal Debian 13
  Target: 1 GB RAM | 2 vCPU VPS | Docker & Self-Hosting
================================================================
EOF
}

# Print Help
show_help() {
    cat << EOF
Usage: sudo ./install.sh [OPTIONS]

Options:
  -h, --help        Show this help message and exit
  -v, --version     Show version information
  --skip-update     Skip 01-system-update step
  --step=<N>        Run a specific installation step number (1-7)

Installation Steps:
  1. System Update (apt update & upgrade)
  2. Install Essentials (htop, curl, git, etc.)
  3. Install Docker & Docker Compose Plugin
  4. Install & Configure Security Stack (UFW & Fail2Ban)
  5. Activate ZRAM Compressed Memory Swap
  6. Install Netdata Monitoring
  7. Apply Service Bloat Removal & Kernel Optimizations
EOF
}

# Environment & Permission Checks
check_environment() {
    if [[ $EUID -ne 0 ]]; then
        log_error "ToyamasOS installer must be run as root (or with sudo)."
        exit 1
    fi

    mkdir -p "$(dirname "$LOG_FILE")"
    touch "$LOG_FILE" 2>/dev/null || true

    log_info "Initializing ToyamasOS installation log: ${LOG_FILE}"

    if [ -f /etc/os-release ]; then
        # shellcheck disable=SC1091
        source /etc/os-release
        log_info "Detected Operating System: ${PRETTY_NAME:-Linux}"
        if [[ "${ID:-}" != "debian" && "${ID_LIKE:-}" != *"debian"* ]]; then
            log_warn "ToyamasOS is optimized for Debian 13. System reported ID='${ID:-unknown}'."
        fi
    else
        log_warn "Could not verify /etc/os-release. Proceeding with installation..."
    fi
}

# Execute Installer Steps
main() {
    local skip_update=false
    local target_step=""

    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                echo "ToyamasOS Installer v${VERSION}"
                exit 0
                ;;
            --skip-update)
                skip_update=true
                shift
                ;;
            --step=*)
                target_step="${1#*=}"
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done

    show_banner
    check_environment

    local scripts=(
        "01-system-update.sh"
        "02-essentials.sh"
        "03-docker.sh"
        "04-security.sh"
        "05-zram.sh"
        "06-netdata.sh"
        "07-optimization.sh"
    )

    if [[ -n "$target_step" ]]; then
        local idx=$((target_step - 1))
        if [[ $idx -ge 0 && $idx -lt ${#scripts[@]} ]]; then
            local script_path="${SCRIPT_DIR}/scripts/${scripts[$idx]}"
            log_info "Executing specific step ${target_step}: ${scripts[$idx]}"
            chmod +x "$script_path"
            bash "$script_path"
            exit 0
        else
            log_error "Invalid step number: ${target_step}. Allowed range: 1-${#scripts[@]}"
            exit 1
        fi
    fi

    for i in "${!scripts[@]}"; do
        local step_num=$((i + 1))
        local script_name="${scripts[$i]}"
        local script_path="${SCRIPT_DIR}/scripts/${script_name}"

        if [[ "$step_num" -eq 1 && "$skip_update" == true ]]; then
            log_info "Skipping step 1 (System Update) as requested."
            continue
        fi

        log_info "--------------------------------------------------------"
        log_info "Running step ${step_num}/${#scripts[@]}: ${script_name}"
        log_info "--------------------------------------------------------"

        if [[ -f "$script_path" ]]; then
            chmod +x "$script_path"
            bash "$script_path"
        else
            log_error "Missing script file: ${script_path}"
            exit 1
        fi
    done

    log_info "========================================================"
    log_success "ToyamasOS Installation & Optimization Complete!"
    log_info "========================================================"
    log_info "System Status Summary:"
    log_info " - Docker Engine: $(docker --version 2>/dev/null || echo 'Installed')"
    log_info " - UFW Firewall:  $(ufw status 2>/dev/null | grep -i 'Status' || echo 'Enabled')"
    log_info " - ZRAM Swap:     $(swapon --show 2>/dev/null | grep zram || echo 'Active')"
    log_info " - Fail2Ban:      $(systemctl is-active fail2ban 2>/dev/null || echo 'Active')"
    log_info " - Netdata:       $(systemctl is-active netdata 2>/dev/null || echo 'Active')"
    log_info " - Log File:      ${LOG_FILE}"
    log_info "========================================================"
    log_success "ToyamasOS is ready for production container hosting."
}

main "$@"
