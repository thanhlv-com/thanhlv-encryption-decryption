#!/bin/bash

# thanhlv-encryption-decryption Auto Install Script
# Supports Linux and macOS
# Usage: curl -sSL https://raw.githubusercontent.com/thanhlv-com/thanhlv-encryption-decryption/main/install.sh | bash
# Or with version: curl -sSL https://raw.githubusercontent.com/thanhlv-com/thanhlv-encryption-decryption/main/install.sh | bash -s -- v1.2.0

set -e

# Configuration
GITHUB_REPO="thanhlv-com/thanhlv-encryption-decryption"
APP_NAME="thanhlv-ed"
DEFAULT_VERSION="v1"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case $os in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        *)
            log_error "Unsupported operating system: $os. This script supports Linux and macOS only."
            ;;
    esac

    case $arch in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch. Supported architectures: x86_64/amd64, arm64/aarch64"
            ;;
    esac

    PLATFORM="${OS}-${ARCH}"
    log_info "Detected platform: ${PLATFORM}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."

    if ! command_exists curl; then
        log_error "curl is required but not installed. Please install curl and try again."
    fi

    if ! command_exists tar; then
        log_error "tar is required but not installed. Please install tar and try again."
    fi

    log_success "Prerequisites check passed"
}

# Get latest release version from GitHub API
get_latest_version() {
    if [ -n "$1" ]; then
        VERSION="$1"
        log_info "Using specified version: $VERSION"
    else
        log_info "Fetching latest release version..."
        VERSION=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

        if [ -z "$VERSION" ]; then
            log_warning "Could not fetch latest version, using default: $DEFAULT_VERSION"
            VERSION="$DEFAULT_VERSION"
        else
            log_info "Latest version: $VERSION"
        fi
    fi
}

# Download and install binary
download_and_install() {
    local binary_name="${APP_NAME}-${PLATFORM}"
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${binary_name}"
    local temp_dir=$(mktemp -d)
    local temp_file="${temp_dir}/${binary_name}"

    log_info "Downloading ${APP_NAME} ${VERSION} for ${PLATFORM}..."
    log_info "Download URL: $download_url"

    # Download binary
    if ! curl -L --fail --silent --show-error "$download_url" -o "$temp_file"; then
        log_error "Failed to download binary from $download_url"
    fi

    # Make binary executable
    chmod +x "$temp_file"

    # Check if we need sudo for installation
    if [ ! -w "$INSTALL_DIR" ]; then
        log_info "Installing to $INSTALL_DIR (requires sudo)..."
        sudo mv "$temp_file" "${INSTALL_DIR}/${APP_NAME}"
    else
        log_info "Installing to $INSTALL_DIR..."
        mv "$temp_file" "${INSTALL_DIR}/${APP_NAME}"
    fi

    # Cleanup
    rm -rf "$temp_dir"

    log_success "${APP_NAME} installed successfully to ${INSTALL_DIR}/${APP_NAME}"
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."

    if command_exists "$APP_NAME"; then
        local installed_version
        installed_version=$($APP_NAME --version 2>/dev/null || echo "unknown")
        log_success "Installation verified! Run '$APP_NAME --help' to get started."
        log_info "Installed version: $installed_version"
    else
        log_warning "Binary installed but not found in PATH. You may need to restart your shell or add $INSTALL_DIR to your PATH."
        log_info "You can run the tool directly: ${INSTALL_DIR}/${APP_NAME} --help"
    fi
}

# Show usage information
show_usage() {
    log_info "Usage examples:"
    echo "  # Generate AES key:"
    echo "  $APP_NAME keygen -a aes-256-cbc -b"
    echo ""
    echo "  # Encrypt text:"
    echo "  $APP_NAME encrypt -a aes-256-cbc -k \"MTIzZGY=\" -t \"Hello World!\""
    echo ""
    echo "  # Decrypt text:"
    echo "  $APP_NAME decrypt -a aes-256-cbc -k \"MTIzZGY=\" -t \"<encrypted-text>\""
    echo ""
    echo "  # Get help:"
    echo "  $APP_NAME --help"
}

# Main installation flow
main() {
    echo ""
    log_info "Starting ${APP_NAME} installation..."
    echo ""

    detect_platform
    check_prerequisites
    get_latest_version "$1"
    download_and_install
    verify_installation

    echo ""
    log_success "Installation completed!"
    echo ""
    show_usage
    echo ""
}

# Handle script interruption
trap 'log_error "Installation interrupted"' INT TERM

# Run main function with all arguments
main "$@"