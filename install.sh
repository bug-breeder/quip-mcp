#!/bin/bash

# Quip MCP Server Installation/Update Script
# Usage: 
#   Install: curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash
#   Update:  curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash -s -- --update

set -e

GITHUB_REPO="bug-breeder/quip-mcp"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
UPDATE_MODE=false

# Show help
show_help() {
    echo "Quip MCP Server Installation/Update Script"
    echo
    echo "Usage:"
    echo "  # Install (first time or reinstall)"
    echo "  curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash"
    echo
    echo "  # Update to latest version"
    echo "  curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash -s -- --update"
    echo
    echo "Options:"
    echo "  --update     Update mode - only installs if newer version available"
    echo "  --help, -h   Show this help message"
    echo
    echo "Environment variables:"
    echo "  INSTALL_DIR  Installation directory (default: /usr/local/bin)"
    echo
    echo "Examples:"
    echo "  # Install to custom directory"
    echo "  INSTALL_DIR=~/.local/bin curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash"
    echo
    echo "  # Update existing installation"
    echo "  curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash -s -- --update"
}

# Parse command line arguments
for arg in "$@"; do
    case $arg in
        --update)
            UPDATE_MODE=true
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
    esac
done

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Detect OS and architecture
detect_platform() {
    OS="$(uname -s)"
    ARCH="$(uname -m)"
    
    case $OS in
        Darwin)
            OS="darwin"
            ;;
        Linux)
            OS="linux"
            ;;
        MINGW* | MSYS* | CYGWIN*)
            OS="windows"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
    
    case $ARCH in
        x86_64 | amd64)
            ARCH="amd64"
            ;;
        arm64 | aarch64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    print_status "Getting latest release version..."
    VERSION=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
    if [ -z "$VERSION" ]; then
        print_error "Failed to get latest version"
        exit 1
    fi
    print_status "Latest version: $VERSION"
}

# Compare version strings (returns 0 if v1 == v2, 1 if v1 > v2, 2 if v1 < v2)
compare_versions() {
    local v1="$1"
    local v2="$2"
    
    # Remove 'v' prefix if present
    v1="${v1#v}"
    v2="${v2#v}"
    
    if [ "$v1" = "$v2" ]; then
        return 0
    fi
    
    # Use sort to compare versions
    local sorted=$(printf '%s\n%s\n' "$v1" "$v2" | sort -V | head -n1)
    if [ "$sorted" = "$v1" ]; then
        return 2  # v1 < v2
    else
        return 1  # v1 > v2
    fi
}

# Download and install
install_quip_mcp() {
    local binary_name="quip-mcp"
    if [ "$OS" = "windows" ]; then
        binary_name="quip-mcp.exe"
    fi
    
    local archive_name="quip-mcp_${VERSION#v}_${OS}_${ARCH}"
    local archive_ext=".tar.gz"
    if [ "$OS" = "windows" ]; then
        archive_ext=".zip"
    fi
    
    local download_url="https://github.com/$GITHUB_REPO/releases/download/$VERSION/${archive_name}${archive_ext}"
    
    print_status "Downloading from: $download_url"
    
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    # Download the archive
    if command -v curl >/dev/null 2>&1; then
        curl -sSL "$download_url" -o "archive${archive_ext}"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$download_url" -O "archive${archive_ext}"
    else
        print_error "Neither curl nor wget is available"
        exit 1
    fi
    
    # Extract the archive
    if [ "$OS" = "windows" ]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "archive${archive_ext}"
        else
            print_error "unzip is required for Windows installation"
            exit 1
        fi
    else
        tar -xzf "archive${archive_ext}"
    fi
    
    # Make binary executable
    chmod +x "$binary_name"
    
    # Install binary
    if [ "$INSTALL_DIR" = "/usr/local/bin" ] && [ ! -w "$INSTALL_DIR" ]; then
        print_status "Installing to $INSTALL_DIR (requires sudo)..."
        sudo mv "$binary_name" "$INSTALL_DIR/"
    else
        print_status "Installing to $INSTALL_DIR..."
        mv "$binary_name" "$INSTALL_DIR/"
    fi
    
    # Cleanup
    cd ..
    rm -rf "$TEMP_DIR"
    
    print_success "quip-mcp installed successfully!"
}

# Check if quip-mcp is already installed and handle updates
check_existing() {
    if command -v quip-mcp >/dev/null 2>&1; then
        CURRENT_VERSION=$(quip-mcp --version 2>/dev/null | head -1 | awk '{print $2}' || echo "unknown")
        
        if [ "$CURRENT_VERSION" = "unknown" ]; then
            print_warning "Found existing quip-mcp installation, but cannot determine version"
            CURRENT_VERSION="unknown"
        else
            print_status "Found existing installation: v$CURRENT_VERSION"
        fi
        
        # Compare versions if we can determine current version
        if [ "$CURRENT_VERSION" != "unknown" ]; then
            set +e  # Temporarily disable exit on error
            compare_versions "$CURRENT_VERSION" "${VERSION#v}"
            VERSION_COMPARISON=$?
            set -e  # Re-enable exit on error
            case $VERSION_COMPARISON in
                0)
                    print_success "You already have the latest version (v$CURRENT_VERSION)"
                    if [ "$UPDATE_MODE" = true ]; then
                        print_status "Update completed - no changes needed"
                        exit 0
                    else
                        read -p "Reinstall anyway? [y/N] " -n 1 -r
                        echo
                        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                            print_status "Installation cancelled"
                            exit 0
                        fi
                    fi
                    ;;
                1)
                    print_warning "You have a newer version (v$CURRENT_VERSION) than the latest release (${VERSION})"
                    if [ "$UPDATE_MODE" = true ]; then
                        print_status "No update needed - you have a newer version"
                        exit 0
                    else
                        read -p "Downgrade to ${VERSION}? [y/N] " -n 1 -r
                        echo
                        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                            print_status "Installation cancelled"
                            exit 0
                        fi
                    fi
                    ;;
                2)
                    if [ "$UPDATE_MODE" = true ]; then
                        print_status "Updating from v$CURRENT_VERSION to ${VERSION}..."
                    else
                        print_warning "Updating existing installation from v$CURRENT_VERSION to ${VERSION}"
                        read -p "Continue? [Y/n] " -n 1 -r
                        echo
                        if [[ $REPLY =~ ^[Nn]$ ]]; then
                            print_status "Installation cancelled"
                            exit 0
                        fi
                    fi
                    ;;
            esac
        else
            # Unknown version case
            if [ "$UPDATE_MODE" = true ]; then
                print_warning "Cannot determine current version for update comparison"
                print_status "Proceeding with installation of ${VERSION}..."
            else
                print_warning "This will overwrite the existing installation"
                read -p "Continue? [y/N] " -n 1 -r
                echo
                if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                    print_status "Installation cancelled"
                    exit 0
                fi
            fi
        fi
    elif [ "$UPDATE_MODE" = true ]; then
        print_error "quip-mcp is not installed. Use regular installation instead:"
        print_error "curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash"
        exit 1
    fi
}

# Show next steps
show_next_steps() {
    echo
    if [ "$UPDATE_MODE" = true ]; then
        print_success "Update complete!"
        print_status "quip-mcp updated to ${VERSION}"
    else
        print_success "Installation complete!"
        echo
        echo "Next steps:"
        echo "1. Get your Quip API token from your Quip instance"
        echo "2. Run: quip-mcp --setup"
        echo "3. Add to your MCP client configuration:"
        echo
        echo '   {
     "mcpServers": {
       "quip": {
         "command": "quip-mcp"
       }
     }
   }'
    fi
    echo
    echo "For more information, visit: https://github.com/$GITHUB_REPO"
}

# Show help
show_help() {
    echo "Quip MCP Server Installation/Update Script"
    echo
    echo "Usage:"
    echo "  # Install (first time or reinstall)"
    echo "  curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash"
    echo
    echo "  # Update to latest version"
    echo "  curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash -s -- --update"
    echo
    echo "Options:"
    echo "  --update     Update mode - only installs if newer version available"
    echo "  --help, -h   Show this help message"
    echo
    echo "Environment variables:"
    echo "  INSTALL_DIR  Installation directory (default: /usr/local/bin)"
    echo
    echo "Examples:"
    echo "  # Install to custom directory"
    echo "  INSTALL_DIR=~/.local/bin curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash"
    echo
    echo "  # Update existing installation"
    echo "  curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash -s -- --update"
}

# Main installation flow
main() {
    if [ "$UPDATE_MODE" = true ]; then
        echo "🔄 Quip MCP Server Updater"
    else
        echo "🚀 Quip MCP Server Installer"
    fi
    echo "============================"
    echo
    
    # Check dependencies
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        print_error "Either curl or wget is required"
        exit 1
    fi
    
    detect_platform
    print_status "Detected platform: $OS/$ARCH"
    
    get_latest_version
    check_existing
    install_quip_mcp
    show_next_steps
}

main "$@" 