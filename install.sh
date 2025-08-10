#!/bin/bash

# ViberCode CLI Installation Script
# This script installs ViberCode CLI on macOS, Linux, and Windows (WSL)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
REPO="vibercode/cli"
BINARY_NAME="vibercode"
INSTALL_DIR="/usr/local/bin"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${PURPLE}"
    echo "╔════════════════════════════════════════╗"
    echo "║          ViberCode CLI Installer      ║"
    echo "║     Generate Go APIs with superpowers ║"
    echo "╚════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Function to detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $os in
        darwin)
            os="darwin"
            ;;
        linux)
            os="linux"
            ;;
        mingw*|msys*|cygwin*)
            os="windows"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    case $arch in
        x86_64|amd64)
            arch="x86_64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        i386|i686)
            arch="i386"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    echo "${os}_${arch}"
}

# Function to get the latest release version
get_latest_version() {
    print_status "Fetching latest release information..."
    
    local version=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | \
                    grep '"tag_name":' | \
                    sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$version" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi
    
    echo "$version"
}

# Function to download and install the binary
install_binary() {
    local version=$1
    local platform=$2
    local temp_dir=$(mktemp -d)
    
    print_status "Installing ViberCode CLI $version for $platform..."
    
    # Construct download URL
    local download_url="https://github.com/$REPO/releases/download/$version/${BINARY_NAME}_${platform}.tar.gz"
    
    if [[ $platform == *"windows"* ]]; then
        download_url="https://github.com/$REPO/releases/download/$version/${BINARY_NAME}_${platform}.zip"
    fi
    
    print_status "Downloading from: $download_url"
    
    # Download the archive
    if ! curl -L -o "$temp_dir/archive" "$download_url"; then
        print_error "Failed to download ViberCode CLI"
        print_error "URL: $download_url"
        exit 1
    fi
    
    # Extract the archive
    cd "$temp_dir"
    if [[ $platform == *"windows"* ]]; then
        unzip -q archive
    else
        tar -xzf archive
    fi
    
    # Make binary executable
    chmod +x "$BINARY_NAME"
    
    # Install to system directory
    print_status "Installing to $INSTALL_DIR..."
    
    if [ ! -w "$INSTALL_DIR" ]; then
        print_status "Installing with sudo (requires password)..."
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    else
        mv "$BINARY_NAME" "$INSTALL_DIR/"
    fi
    
    # Clean up
    rm -rf "$temp_dir"
    
    print_success "ViberCode CLI installed successfully!"
}

# Function to verify installation
verify_installation() {
    print_status "Verifying installation..."
    
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version=$($BINARY_NAME --version 2>/dev/null || echo "unknown")
        print_success "ViberCode CLI is installed and working!"
        print_status "Version: $version"
        print_status "Location: $(which $BINARY_NAME)"
    else
        print_error "Installation verification failed"
        print_error "Binary not found in PATH"
        exit 1
    fi
}

# Function to show next steps
show_next_steps() {
    echo
    print_success "🎉 Installation Complete!"
    echo
    echo -e "${CYAN}Next steps:${NC}"
    echo "  1. ${GREEN}vibercode --help${NC}          # Show available commands"
    echo "  2. ${GREEN}vibercode vibe${NC}            # Start full development mode"
    echo "  3. ${GREEN}vibercode generate api${NC}    # Generate a new Go API"
    echo
    echo -e "${CYAN}Quick start:${NC}"
    echo "  ${GREEN}vibercode vibe${NC}  # Launch visual editor + AI chat"
    echo
    echo -e "${CYAN}Documentation:${NC}"
    echo "  📖 https://github.com/vibercode/cli"
    echo "  💬 Join our community for support and updates"
    echo
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check for curl
    if ! command -v curl >/dev/null 2>&1; then
        print_error "curl is required but not installed"
        print_status "Please install curl and try again"
        exit 1
    fi
    
    # Check for unzip (Windows)
    local platform=$(detect_platform)
    if [[ $platform == *"windows"* ]] && ! command -v unzip >/dev/null 2>&1; then
        print_error "unzip is required but not installed"
        print_status "Please install unzip and try again"
        exit 1
    fi
    
    # Check for tar (Unix-like)
    if [[ $platform != *"windows"* ]] && ! command -v tar >/dev/null 2>&1; then
        print_error "tar is required but not installed"
        print_status "Please install tar and try again"
        exit 1
    fi
    
    print_success "Prerequisites check passed"
}

# Main installation function
main() {
    print_header
    
    # Check prerequisites
    check_prerequisites
    
    # Detect platform
    local platform=$(detect_platform)
    print_status "Detected platform: $platform"
    
    # Get latest version
    local version=$(get_latest_version)
    print_status "Latest version: $version"
    
    # Check if already installed
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local current_version=$($BINARY_NAME --version 2>/dev/null | grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+' || echo "unknown")
        print_warning "ViberCode CLI is already installed (version: $current_version)"
        
        if [ "$current_version" = "$version" ]; then
            print_status "You already have the latest version installed"
            read -p "Do you want to reinstall? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                print_status "Installation cancelled"
                exit 0
            fi
        fi
    fi
    
    # Install the binary
    install_binary "$version" "$platform"
    
    # Verify installation
    verify_installation
    
    # Show next steps
    show_next_steps
}

# Handle script interruption
trap 'print_error "Installation interrupted"; exit 1' INT TERM

# Run main function
main "$@"