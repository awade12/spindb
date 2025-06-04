#!/usr/bin/env bash

set -euo pipefail

# Enable debug mode if DEBUG=1 is set
if [ "${DEBUG:-0}" = "1" ]; then
    set -x
fi

OWNER="awade12"
REPO="spindb"
BINARY_NAME="spindb"
INSTALL_DIR="/usr/local/bin"
TEMP_DIR="/tmp/spindb_install"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

detect_os() {
    case "$(uname -s)" in
        Darwin*)
            OS="darwin"
            ;;
        Linux*)
            OS="linux"
            ;;
        MINGW* | MSYS* | CYGWIN*)
            OS="windows"
            ;;
        *)
            error "Unsupported operating system: $(uname -s)"
            ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64 | amd64)
            ARCH="amd64"
            ;;
        arm64 | aarch64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $(uname -m)"
            ;;
    esac
}

check_dependencies() {
    if ! command -v curl >/dev/null 2>&1; then
        error "curl is required but not installed. Please install curl and try again."
    fi
    
    if ! command -v tar >/dev/null 2>&1 && ! command -v unzip >/dev/null 2>&1; then
        error "Either tar or unzip is required but neither is installed."
    fi
}

check_docker() {
    if command -v docker >/dev/null 2>&1; then
        log "Docker is already installed"
        if docker ps >/dev/null 2>&1; then
            log "✓ Docker daemon is running"
            return 0
        else
            warn "Docker is installed but daemon is not running"
            log "Attempting to start Docker daemon..."
            start_docker_daemon
        fi
    else
        log "Docker not found. Installing Docker..."
        install_docker
    fi
}

start_docker_daemon() {
    case "$OS" in
        "linux")
            if command -v systemctl >/dev/null 2>&1; then
                log "Starting Docker with systemctl..."
                sudo systemctl start docker
                sudo systemctl enable docker
            elif command -v service >/dev/null 2>&1; then
                log "Starting Docker with service..."
                sudo service docker start
            else
                warn "Could not start Docker daemon automatically. Please start Docker manually."
            fi
            ;;
        "darwin")
            log "Please start Docker Desktop manually if it's not running"
            log "You can find Docker Desktop in your Applications folder"
            ;;
        "windows")
            log "Please start Docker Desktop manually if it's not running"
            ;;
    esac
}

install_docker() {
    case "$OS" in
        "linux")
            install_docker_linux
            ;;
        "darwin")
            install_docker_macos
            ;;
        "windows")
            install_docker_windows
            ;;
        *)
            error "Unsupported operating system for Docker installation: $OS"
            ;;
    esac
}

install_docker_linux() {
    log "Installing Docker on Linux..."
    
    if command -v apt-get >/dev/null 2>&1; then
        log "Using apt package manager..."
        sudo apt-get update
        sudo apt-get install -y ca-certificates curl gnupg lsb-release
        
        sudo mkdir -p /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
        
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
        
        sudo apt-get update
        sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        
    elif command -v yum >/dev/null 2>&1; then
        log "Using yum package manager..."
        sudo yum install -y yum-utils
        sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
        sudo yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        
    elif command -v dnf >/dev/null 2>&1; then
        log "Using dnf package manager..."
        sudo dnf -y install dnf-plugins-core
        sudo dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
        sudo dnf install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        
    else
        log "Using Docker's convenience script..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        rm get-docker.sh
    fi
    
    sudo usermod -aG docker $USER
    sudo systemctl start docker
    sudo systemctl enable docker
    
    log "✓ Docker installed successfully"
    log "Note: You may need to log out and back in for Docker group permissions to take effect"
}

install_docker_macos() {
    log "Installing Docker Desktop on macOS..."
    
    if [ "$ARCH" = "arm64" ]; then
        DOCKER_URL="https://desktop.docker.com/mac/main/arm64/Docker.dmg"
    else
        DOCKER_URL="https://desktop.docker.com/mac/main/amd64/Docker.dmg"
    fi
    
    log "Downloading Docker Desktop..."
    curl -L -o "$TEMP_DIR/Docker.dmg" "$DOCKER_URL"
    
    log "Mounting Docker Desktop installer..."
    hdiutil attach "$TEMP_DIR/Docker.dmg" -quiet
    
    log "Installing Docker Desktop..."
    sudo cp -R "/Volumes/Docker/Docker.app" "/Applications/"
    
    log "Unmounting installer..."
    hdiutil detach "/Volumes/Docker" -quiet
    
    log "✓ Docker Desktop installed successfully"
    log "Please start Docker Desktop from your Applications folder"
    log "Wait for Docker Desktop to start before using SpinDB"
}

install_docker_windows() {
    log "Installing Docker Desktop on Windows..."
    
    DOCKER_URL="https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe"
    
    log "Downloading Docker Desktop installer..."
    curl -L -o "$TEMP_DIR/DockerDesktopInstaller.exe" "$DOCKER_URL"
    
    log "Running Docker Desktop installer..."
    log "Please follow the installation prompts..."
    "$TEMP_DIR/DockerDesktopInstaller.exe"
    
    log "✓ Docker Desktop installer launched"
    log "Please complete the installation and start Docker Desktop"
    log "Wait for Docker Desktop to start before using SpinDB"
}

get_latest_release() {
    log "Fetching latest release information..."
    
    if command -v jq >/dev/null 2>&1; then
        LATEST_TAG=$(curl -s "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" | jq -r '.tag_name')
    else
        LATEST_TAG=$(curl -s "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi
    
    if [ -z "$LATEST_TAG" ] || [ "$LATEST_TAG" = "null" ]; then
        warn "Could not determine latest release. Using 'latest' as fallback."
        LATEST_TAG="latest"
    fi
    
    log "Latest release: $LATEST_TAG"
}

download_binary() {
    if [ "$OS" = "windows" ]; then
        BINARY_FILENAME="${BINARY_NAME}-${OS}-${ARCH}.exe"
        LOCAL_BINARY_NAME="${BINARY_NAME}.exe"
    else
        BINARY_FILENAME="${BINARY_NAME}-${OS}-${ARCH}"
        LOCAL_BINARY_NAME="${BINARY_NAME}"
    fi
    
    log "Detected platform: ${OS}-${ARCH}"
    log "Downloading ${BINARY_FILENAME}..."
    
    mkdir -p "$TEMP_DIR"
    cd "$TEMP_DIR"
    
    if [ "$LATEST_TAG" = "latest" ]; then
        DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/latest/download/${BINARY_FILENAME}"
    else
        DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/download/${LATEST_TAG}/${BINARY_FILENAME}"
    fi
    
    if curl -L -o "$LOCAL_BINARY_NAME" "$DOWNLOAD_URL" 2>/dev/null; then
        chmod +x "$LOCAL_BINARY_NAME"
        log "Download completed successfully"
    else
        warn "Pre-built binary not available, building from source..."
        build_from_source
    fi
}

build_from_source() {
    log "Building SpinDB from source..."
    
    check_build_dependencies
    
    log "Cloning repository..."
    if ! git clone "https://github.com/${OWNER}/${REPO}.git" "${REPO}"; then
        error "Failed to clone repository"
    fi
    
    cd "${REPO}"
    
    if [ "$LATEST_TAG" != "latest" ]; then
        log "Checking out tag ${LATEST_TAG}..."
        git checkout "${LATEST_TAG}"
    fi
    
    log "Downloading Go dependencies..."
    if ! go mod download; then
        error "Failed to download Go dependencies"
    fi
    
    log "Building binary..."
    if ! go build -ldflags "-s -w" -o "../${LOCAL_BINARY_NAME}" .; then
        error "Failed to build binary"
    fi
    
    cd ..
    chmod +x "$LOCAL_BINARY_NAME"
    log "Build completed successfully"
}

check_build_dependencies() {
    if ! command -v git >/dev/null 2>&1; then
        error "git is required for building from source but not installed"
    fi
    
    if ! command -v go >/dev/null 2>&1; then
        log "Go not found, installing..."
        install_go
    fi
}

install_go() {
    case "$OS" in
        "linux")
            log "Installing Go on Linux..."
            if command -v apt-get >/dev/null 2>&1; then
                sudo apt-get update
                sudo apt-get install -y golang-go
            elif command -v yum >/dev/null 2>&1; then
                sudo yum install -y golang
            elif command -v dnf >/dev/null 2>&1; then
                sudo dnf install -y golang
            else
                install_go_from_source
            fi
            ;;
        "darwin")
            log "Installing Go on macOS..."
            if command -v brew >/dev/null 2>&1; then
                brew install go
            else
                install_go_from_source
            fi
            ;;
        *)
            install_go_from_source
            ;;
    esac
}

install_go_from_source() {
    log "Installing Go from official installer..."
    
    case "$OS" in
        "linux")
            if [ "$ARCH" = "amd64" ]; then
                GO_URL="https://golang.org/dl/go1.21.5.linux-amd64.tar.gz"
            elif [ "$ARCH" = "arm64" ]; then
                GO_URL="https://golang.org/dl/go1.21.5.linux-arm64.tar.gz"
            else
                error "Unsupported architecture for Go installation: $ARCH"
            fi
            ;;
        "darwin")
            if [ "$ARCH" = "amd64" ]; then
                GO_URL="https://golang.org/dl/go1.21.5.darwin-amd64.tar.gz"
            elif [ "$ARCH" = "arm64" ]; then
                GO_URL="https://golang.org/dl/go1.21.5.darwin-arm64.tar.gz"
            else
                error "Unsupported architecture for Go installation: $ARCH"
            fi
            ;;
        *)
            error "Unsupported OS for Go installation: $OS"
            ;;
    esac
    
    curl -L -o go.tar.gz "$GO_URL"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go.tar.gz
    
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    
    if ! command -v go >/dev/null 2>&1; then
        error "Go installation failed"
    fi
    
    log "Go installed successfully"
}

install_binary() {
    log "Installing ${BINARY_NAME} to ${INSTALL_DIR}..."
    
    if [ "$OS" = "windows" ]; then
        INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}.exe"
    else
        INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    if [ -w "$INSTALL_DIR" ]; then
        cp "$LOCAL_BINARY_NAME" "$INSTALL_PATH"
    else
        log "Requesting sudo privileges to install to $INSTALL_DIR"
        sudo cp "$LOCAL_BINARY_NAME" "$INSTALL_PATH"
    fi
    
    log "Installation completed successfully"
}

verify_installation() {
    log "Verifying installation..."
    
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        log "✓ ${BINARY_NAME} is now available in your PATH"
        log "Running version check..."
        "$BINARY_NAME" version 2>/dev/null || "$BINARY_NAME" --version 2>/dev/null || log "Binary installed but version command not available"
    else
        warn "${BINARY_NAME} was installed but is not in your PATH."
        warn "You may need to add ${INSTALL_DIR} to your PATH or restart your terminal."
    fi
}

cleanup() {
    if [ -n "${TEMP_DIR:-}" ] && [ -d "$TEMP_DIR" ]; then
        log "Cleaning up temporary files..."
        rm -rf "$TEMP_DIR"
    fi
}

show_usage() {
    echo -e "${BOLD}SpinDB Installation Complete!${NC}"
    echo ""
    
    if command -v docker >/dev/null 2>&1 && docker ps >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Docker is running and ready${NC}"
    else
        echo -e "${YELLOW}⚠ Make sure Docker is running before using SpinDB${NC}"
    fi
    
    echo ""
    echo "You can now use SpinDB with the following commands:"
    echo ""
    echo -e "  ${BLUE}spindb create${NC}     - Create a new database instance"
    echo -e "  ${BLUE}spindb list${NC}       - List all database instances"
    echo -e "  ${BLUE}spindb start${NC}      - Start a database instance"
    echo -e "  ${BLUE}spindb stop${NC}       - Stop a database instance"
    echo -e "  ${BLUE}spindb connect${NC}    - Connect to a database instance"
    echo -e "  ${BLUE}spindb backup${NC}     - Backup a database instance"
    echo -e "  ${BLUE}spindb env${NC}        - Manage environment variables"
    echo -e "  ${BLUE}spindb template${NC}   - Manage database templates"
    echo -e "  ${BLUE}spindb info${NC}       - Show instance information"
    echo -e "  ${BLUE}spindb delete${NC}     - Delete a database instance"
    echo -e "  ${BLUE}spindb version${NC}    - Show SpinDB version"
    echo ""
    echo "For more help, run: ${BLUE}spindb --help${NC}"
    echo ""
    echo -e "${GREEN}Happy database spinning! 🚀${NC}"
}

main() {
    echo -e "${BOLD}${BLUE}SpinDB Installer${NC}"
    echo "=================="
    echo ""
    
    log "Starting installation process..."
    
    log "Checking dependencies..."
    check_dependencies
    
    log "Detecting OS and architecture..."
    detect_os
    detect_arch
    
    log "Checking Docker..."
    check_docker
    
    log "Getting latest release information..."
    get_latest_release
    
    log "Downloading or building binary..."
    download_binary
    
    log "Installing binary..."
    install_binary
    
    log "Verifying installation..."
    verify_installation
    
    log "Installation process completed"
    
    echo ""
    show_usage
}

trap cleanup EXIT

# Only run main if script is executed directly (not sourced)
if [ "${BASH_SOURCE:-}" = "" ] || [ "${BASH_SOURCE[0]:-}" = "${0}" ]; then
    main "$@"
fi
