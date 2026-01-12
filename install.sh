#!/bin/bash
# AgentSync Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/AlfonsSkills/AgentSync/main/install.sh | bash

set -e

REPO="AlfonsSkills/AgentSync"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="agentsync"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        darwin) OS="darwin" ;;
        linux) OS="linux" ;;
        mingw*|msys*|cygwin*) OS="windows" ;;
        *) error "Unsupported OS: $OS" ;;
    esac

    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac

    PLATFORM="${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        PLATFORM="${PLATFORM}.exe"
    fi
}

# Get latest release version
get_latest_version() {
    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        error "Failed to get latest version"
    fi
    info "Latest version: $VERSION"
}

# Download and install
install() {
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}-${PLATFORM}"
    TMP_FILE=$(mktemp)

    info "Downloading ${BINARY_NAME}-${PLATFORM}..."
    if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"; then
        error "Failed to download from $DOWNLOAD_URL"
    fi

    chmod +x "$TMP_FILE"

    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        warn "Need sudo to install to $INSTALL_DIR"
        sudo mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    info "Installed to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Verify installation
verify() {
    if command -v "$BINARY_NAME" &> /dev/null; then
        info "Installation successful!"
        echo ""
        $BINARY_NAME --version
    else
        warn "Installation complete, but '${BINARY_NAME}' not in PATH"
        warn "Add ${INSTALL_DIR} to your PATH, or run: ${INSTALL_DIR}/${BINARY_NAME}"
    fi
}

main() {
    echo ""
    echo "🚀 AgentSync Installer"
    echo ""

    detect_platform
    info "Detected platform: $PLATFORM"

    get_latest_version
    install
    verify

    echo ""
    info "Run 'agentsync --help' to get started"
}

main
