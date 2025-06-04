#!/usr/bin/env bash

set -e

# Enable debug mode if DEBUG=1 is set
if [ "${DEBUG:-0}" = "1" ]; then
    set -x
fi

OWNER="awade12"
REPO="spindb"
BINARY_NAME="spindb"
INSTALL_DIR="/usr/local/bin"
TEMP_DIR="/tmp/spindb_install"

# Enhanced color palette
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
BOLD='\033[1m'
DIM='\033[2m'
UNDERLINE='\033[4m'
BLINK='\033[5m'
REVERSE='\033[7m'
NC='\033[0m'

# Gradient colors for fancy effects
GRAD1='\033[38;5;93m'   # Purple
GRAD2='\033[38;5;99m'   # Light Purple
GRAD3='\033[38;5;105m'  # Blue Purple
GRAD4='\033[38;5;111m'  # Light Blue
GRAD5='\033[38;5;117m'  # Cyan Blue

# Animation characters
SPINNER_CHARS="⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
PROGRESS_CHARS="▏▎▍▌▋▊▉█"

# Cool ASCII art
show_banner() {
    clear
    echo -e "${GRAD1}"
    cat << 'EOF'
    ███████╗██████╗ ██╗███╗   ██╗██████╗ ██████╗ 
    ██╔════╝██╔══██╗██║████╗  ██║██╔══██╗██╔══██╗
    ███████╗██████╔╝██║██╔██╗ ██║██║  ██║██████╔╝
    ╚════██║██╔═══╝ ██║██║╚██╗██║██║  ██║██╔══██╗
    ███████║██║     ██║██║ ╚████║██████╔╝██████╔╝
    ╚══════╝╚═╝     ╚═╝╚═╝  ╚═══╝╚═════╝ ╚═════╝ 
EOF
    echo -e "${NC}"
    echo -e "${GRAD2}         🚀 Database Management Made Easy 🚀${NC}"
    echo -e "${GRAD3}            ═══════════════════════════════${NC}"
    echo ""
}

# Animated spinner function
spinner() {
    local pid=$1
    local message=$2
    local delay=0.1
    local i=0
    
    while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
        local char=${SPINNER_CHARS:$i:1}
        printf "\r${CYAN}${char}${NC} ${message}${CYAN}...${NC}"
        sleep $delay
        i=$(((i + 1) % ${#SPINNER_CHARS}))
    done
    printf "\r${GREEN}✓${NC} ${message}${GREEN} completed!${NC}\n"
}

# Progress bar function
progress_bar() {
    local current=$1
    local total=$2
    local width=50
    local percentage=$((current * 100 / total))
    local completed=$((current * width / total))
    local remaining=$((width - completed))
    
    printf "\r${BLUE}["
    for ((i=0; i<completed; i++)); do
        printf "${GREEN}█${NC}"
    done
    for ((i=0; i<remaining; i++)); do
        printf "${DIM}░${NC}"
    done
    printf "${BLUE}] ${WHITE}${percentage}%%${NC}"
}

# Fancy logging functions
log() {
    echo -e "${GREEN}[${BOLD}INFO${NC}${GREEN}]${NC} ${WHITE}$1${NC}"
}

warn() {
    echo -e "${YELLOW}[${BOLD}WARN${NC}${YELLOW}]${NC} ${YELLOW}$1${NC}"
}

error() {
    echo -e "${RED}[${BOLD}ERROR${NC}${RED}]${NC} ${RED}$1${NC}"
    exit 1
}

success() {
    echo -e "${GREEN}[${BOLD}SUCCESS${NC}${GREEN}]${NC} ${GREEN}$1${NC}"
}

# Cool section headers
section_header() {
    echo ""
    echo -e "${GRAD4}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GRAD4}║${NC} ${BOLD}${WHITE}$1${NC}${GRAD4} ║${NC}"
    echo -e "${GRAD4}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

# Animated typing effect
type_text() {
    local text="$1"
    local delay=${2:-0.03}
    
    for ((i=0; i<${#text}; i++)); do
        printf "${text:$i:1}"
        sleep $delay
    done
    echo ""
}

# Loading animation
loading_animation() {
    local message="$1"
    local duration=${2:-3}
    
    echo -e "${CYAN}${message}${NC}"
    
    for ((i=0; i<duration*4; i++)); do
        local frame=$((i % 4))
        case $frame in
            0) printf "\r${BLUE}⠋${NC} Processing..." ;;
            1) printf "\r${BLUE}⠙${NC} Processing..." ;;
            2) printf "\r${BLUE}⠹${NC} Processing..." ;;
            3) printf "\r${BLUE}⠸${NC} Processing..." ;;
        esac
        sleep 0.25
    done
    printf "\r${GREEN}✓${NC} Complete!        \n"
}

# System info display
show_system_info() {
    echo -e "${GRAD3}┌─ System Information ─────────────────────────────────────────┐${NC}"
    echo -e "${GRAD3}│${NC} ${BOLD}OS:${NC}           $(uname -s)"
    echo -e "${GRAD3}│${NC} ${BOLD}Architecture:${NC} $(uname -m)"
    echo -e "${GRAD3}│${NC} ${BOLD}Kernel:${NC}       $(uname -r)"
    echo -e "${GRAD3}│${NC} ${BOLD}User:${NC}         $(whoami)"
    echo -e "${GRAD3}└──────────────────────────────────────────────────────────────┘${NC}"
    echo ""
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
    section_header "🔍 Checking Dependencies"
    
    local deps=("curl" "tar")
    local missing=()
    
    for dep in "${deps[@]}"; do
        if command -v "$dep" >/dev/null 2>&1; then
            success "$dep is installed ✓"
        else
            missing+=("$dep")
            error "$dep is missing ✗"
        fi
        sleep 0.2
    done
    
    if [ ${#missing[@]} -gt 0 ]; then
        error "Missing dependencies: ${missing[*]}. Please install them and try again."
    fi
    
    if ! command -v unzip >/dev/null 2>&1 && ! command -v tar >/dev/null 2>&1; then
        error "Either tar or unzip is required but neither is installed."
    fi
}

check_docker() {
    section_header "🐳 Docker Setup"
    
    if command -v docker >/dev/null 2>&1; then
        success "Docker is installed!"
        
        # Animated check for Docker daemon
        printf "${CYAN}🔄 Checking Docker daemon status...${NC}"
        sleep 1
        
        if docker ps >/dev/null 2>&1; then
            printf "\r${GREEN}✅ Docker daemon is running and ready!${NC}          \n"
            return 0
        else
            printf "\r${YELLOW}⚠️  Docker daemon is not running${NC}               \n"
            log "Attempting to start Docker daemon..."
            start_docker_daemon
        fi
    else
        warn "Docker not found. Installing Docker..."
        loading_animation "🚀 Preparing Docker installation" 2
        install_docker
    fi
}

start_docker_daemon() {
    case "$OS" in
        "linux")
            if command -v systemctl >/dev/null 2>&1; then
                log "🔧 Starting Docker with systemctl..."
                sudo systemctl start docker
                sudo systemctl enable docker
            elif command -v service >/dev/null 2>&1; then
                log "🔧 Starting Docker with service..."
                sudo service docker start
            else
                warn "Could not start Docker daemon automatically. Please start Docker manually."
            fi
            ;;
        "darwin")
            log "🍎 Please start Docker Desktop manually if it's not running"
            log "You can find Docker Desktop in your Applications folder"
            ;;
        "windows")
            log "🪟 Please start Docker Desktop manually if it's not running"
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
    log "🐧 Installing Docker on Linux..."
    
    # Progress simulation
    for i in {1..20}; do
        progress_bar $i 20
        sleep 0.1
    done
    echo ""
    
    if command -v apt-get >/dev/null 2>&1; then
        log "📦 Using apt package manager..."
        sudo apt-get update >/dev/null 2>&1 &
        spinner $! "Updating package lists"
        
        sudo apt-get install -y ca-certificates curl gnupg lsb-release >/dev/null 2>&1 &
        spinner $! "Installing prerequisites"
        
        sudo mkdir -p /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg 2>/dev/null &
        spinner $! "Adding Docker GPG key"
        
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
        
        sudo apt-get update >/dev/null 2>&1 &
        spinner $! "Updating package lists with Docker repository"
        
        sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin >/dev/null 2>&1 &
        spinner $! "Installing Docker components"
        
    elif command -v yum >/dev/null 2>&1; then
        log "📦 Using yum package manager..."
        sudo yum install -y yum-utils >/dev/null 2>&1 &
        spinner $! "Installing yum utilities"
        
        sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo >/dev/null 2>&1 &
        spinner $! "Adding Docker repository"
        
        sudo yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin >/dev/null 2>&1 &
        spinner $! "Installing Docker components"
        
    elif command -v dnf >/dev/null 2>&1; then
        log "📦 Using dnf package manager..."
        sudo dnf -y install dnf-plugins-core >/dev/null 2>&1 &
        spinner $! "Installing dnf plugins"
        
        sudo dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo >/dev/null 2>&1 &
        spinner $! "Adding Docker repository"
        
        sudo dnf install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin >/dev/null 2>&1 &
        spinner $! "Installing Docker components"
        
    else
        log "📦 Using Docker's convenience script..."
        curl -fsSL https://get.docker.com -o get-docker.sh 2>/dev/null &
        spinner $! "Downloading Docker installation script"
        
        sudo sh get-docker.sh >/dev/null 2>&1 &
        spinner $! "Running Docker installation script"
        
        rm get-docker.sh
    fi
    
    sudo usermod -aG docker $USER >/dev/null 2>&1
    sudo systemctl start docker >/dev/null 2>&1
    sudo systemctl enable docker >/dev/null 2>&1
    
    success "Docker installed successfully! 🎉"
    log "📝 Note: You may need to log out and back in for Docker group permissions to take effect"
}

install_docker_macos() {
    log "🍎 Installing Docker Desktop on macOS..."
    
    if [ "$ARCH" = "arm64" ]; then
        DOCKER_URL="https://desktop.docker.com/mac/main/arm64/Docker.dmg"
    else
        DOCKER_URL="https://desktop.docker.com/mac/main/amd64/Docker.dmg"
    fi
    
    curl -L -o "$TEMP_DIR/Docker.dmg" "$DOCKER_URL" 2>/dev/null &
    spinner $! "Downloading Docker Desktop"
    
    hdiutil attach "$TEMP_DIR/Docker.dmg" -quiet 2>/dev/null &
    spinner $! "Mounting Docker Desktop installer"
    
    sudo cp -R "/Volumes/Docker/Docker.app" "/Applications/" 2>/dev/null &
    spinner $! "Installing Docker Desktop"
    
    hdiutil detach "/Volumes/Docker" -quiet 2>/dev/null &
    spinner $! "Cleaning up installer"
    
    success "Docker Desktop installed successfully! 🎉"
    log "📱 Please start Docker Desktop from your Applications folder"
    log "⏳ Wait for Docker Desktop to start before using SpinDB"
}

install_docker_windows() {
    log "🪟 Installing Docker Desktop on Windows..."
    
    DOCKER_URL="https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe"
    
    curl -L -o "$TEMP_DIR/DockerDesktopInstaller.exe" "$DOCKER_URL" 2>/dev/null &
    spinner $! "Downloading Docker Desktop installer"
    
    log "🚀 Running Docker Desktop installer..."
    log "👆 Please follow the installation prompts..."
    "$TEMP_DIR/DockerDesktopInstaller.exe"
    
    success "Docker Desktop installer launched! 🎉"
    log "✅ Please complete the installation and start Docker Desktop"
    log "⏳ Wait for Docker Desktop to start before using SpinDB"
}

get_latest_release() {
    section_header "📡 Fetching Release Information"
    
    printf "${CYAN}🔍 Checking GitHub for latest release...${NC}"
    sleep 1
    
    if command -v jq >/dev/null 2>&1; then
        LATEST_TAG=$(curl -s "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" | jq -r '.tag_name')
    else
        LATEST_TAG=$(curl -s "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi
    
    if [ -z "$LATEST_TAG" ] || [ "$LATEST_TAG" = "null" ]; then
        printf "\r${YELLOW}⚠️  Could not determine latest release${NC}            \n"
        warn "Using 'latest' as fallback."
        LATEST_TAG="latest"
    else
        printf "\r${GREEN}✅ Found latest release: ${BOLD}$LATEST_TAG${NC}         \n"
    fi
}

download_binary() {
    section_header "📥 Downloading SpinDB"
    
    if [ "$OS" = "windows" ]; then
        BINARY_FILENAME="${BINARY_NAME}-${OS}-${ARCH}.exe"
        LOCAL_BINARY_NAME="${BINARY_NAME}.exe"
    else
        BINARY_FILENAME="${BINARY_NAME}-${OS}-${ARCH}"
        LOCAL_BINARY_NAME="${BINARY_NAME}"
    fi
    
    log "🖥️  Platform detected: ${BOLD}${OS}-${ARCH}${NC}"
    log "📦 Binary: ${BOLD}${BINARY_FILENAME}${NC}"
    
    mkdir -p "$TEMP_DIR"
    cd "$TEMP_DIR"
    
    if [ "$LATEST_TAG" = "latest" ]; then
        DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/latest/download/${BINARY_FILENAME}"
    else
        DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/download/${LATEST_TAG}/${BINARY_FILENAME}"
    fi
    
    # Fancy download progress
    echo -e "${CYAN}🚀 Starting download...${NC}"
    
    if curl -L -o "$LOCAL_BINARY_NAME" "$DOWNLOAD_URL" 2>/dev/null; then
        # Simulate progress for visual effect
        for i in {1..30}; do
            progress_bar $i 30
            sleep 0.05
        done
        echo ""
        
        chmod +x "$LOCAL_BINARY_NAME"
        success "Download completed successfully! 🎉"
    else
        warn "❌ Pre-built binary not available"
        log "🔨 Building from source instead..."
        build_from_source
    fi
}

build_from_source() {
    section_header "🔨 Building from Source"
    
    log "🛠️  Building SpinDB from source..."
    
    check_build_dependencies
    
    git clone "https://github.com/${OWNER}/${REPO}.git" "${REPO}" 2>/dev/null &
    spinner $! "Cloning repository"
    
    cd "${REPO}"
    
    if [ "$LATEST_TAG" != "latest" ]; then
        git checkout "${LATEST_TAG}" 2>/dev/null &
        spinner $! "Checking out tag ${LATEST_TAG}"
    fi
    
    go mod download 2>/dev/null &
    spinner $! "Downloading Go dependencies"
    
    go build -ldflags "-s -w" -o "../${LOCAL_BINARY_NAME}" . 2>/dev/null &
    spinner $! "Building binary"
    
    cd ..
    chmod +x "$LOCAL_BINARY_NAME"
    success "Build completed successfully! 🎉"
}

check_build_dependencies() {
    log "🔍 Checking build dependencies..."
    
    if ! command -v git >/dev/null 2>&1; then
        error "git is required for building from source but not installed"
    fi
    
    if ! command -v go >/dev/null 2>&1; then
        log "📦 Go not found, installing..."
        install_go
    fi
}

install_go() {
    case "$OS" in
        "linux")
            log "🐧 Installing Go on Linux..."
            if command -v apt-get >/dev/null 2>&1; then
                sudo apt-get update >/dev/null 2>&1
                sudo apt-get install -y golang-go >/dev/null 2>&1 &
                spinner $! "Installing Go via apt"
            elif command -v yum >/dev/null 2>&1; then
                sudo yum install -y golang >/dev/null 2>&1 &
                spinner $! "Installing Go via yum"
            elif command -v dnf >/dev/null 2>&1; then
                sudo dnf install -y golang >/dev/null 2>&1 &
                spinner $! "Installing Go via dnf"
            else
                install_go_from_source
            fi
            ;;
        "darwin")
            log "🍎 Installing Go on macOS..."
            if command -v brew >/dev/null 2>&1; then
                brew install go >/dev/null 2>&1 &
                spinner $! "Installing Go via Homebrew"
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
    log "📥 Installing Go from official installer..."
    
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
    
    curl -L -o go.tar.gz "$GO_URL" 2>/dev/null &
    spinner $! "Downloading Go"
    
    sudo rm -rf /usr/local/go 2>/dev/null
    sudo tar -C /usr/local -xzf go.tar.gz 2>/dev/null &
    spinner $! "Installing Go"
    
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    
    if ! command -v go >/dev/null 2>&1; then
        error "Go installation failed"
    fi
    
    success "Go installed successfully! 🎉"
}

install_binary() {
    section_header "📦 Installing SpinDB"
    
    if [ "$OS" = "windows" ]; then
        INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}.exe"
    else
        INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    printf "${CYAN}📍 Installing to: ${BOLD}${INSTALL_PATH}${NC}\n"
    
    if [ -w "$INSTALL_DIR" ]; then
        cp "$LOCAL_BINARY_NAME" "$INSTALL_PATH" 2>/dev/null &
        spinner $! "Installing binary"
    else
        log "🔐 Requesting sudo privileges..."
        sudo cp "$LOCAL_BINARY_NAME" "$INSTALL_PATH" 2>/dev/null &
        spinner $! "Installing binary (with sudo)"
    fi
    
    success "Installation completed successfully! 🎉"
}

verify_installation() {
    section_header "✅ Verifying Installation"
    
    printf "${CYAN}🔍 Checking installation...${NC}"
    sleep 1
    
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        printf "\r${GREEN}✅ SpinDB is now available in your PATH!${NC}       \n"
        
        echo -e "${CYAN}🔧 Running version check...${NC}"
        sleep 0.5
        
        if version_output=$("$BINARY_NAME" version 2>/dev/null || "$BINARY_NAME" --version 2>/dev/null); then
            echo -e "${GREEN}📌 Version: ${BOLD}${version_output}${NC}"
        else
            log "Binary installed but version command not available"
        fi
    else
        printf "\r${YELLOW}⚠️  SpinDB installed but not in PATH${NC}           \n"
        warn "You may need to add ${INSTALL_DIR} to your PATH or restart your terminal."
    fi
}

cleanup() {
    if [ -n "${TEMP_DIR:-}" ] && [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR" 2>/dev/null
    fi
}

check_and_offer_client_tools() {
    section_header "🔧 Database Client Tools"
    
    missing_clients=""
    
    # Check each client with fancy output
    echo -e "${CYAN}Checking for database clients...${NC}"
    
    if command -v psql >/dev/null 2>&1; then
        echo -e "${GREEN}✅ PostgreSQL client (psql) - Found${NC}"
    else
        echo -e "${YELLOW}❌ PostgreSQL client (psql) - Missing${NC}"
        missing_clients="${missing_clients}postgresql-client "
    fi
    
    if command -v mysql >/dev/null 2>&1; then
        echo -e "${GREEN}✅ MySQL client (mysql) - Found${NC}"
    else
        echo -e "${YELLOW}❌ MySQL client (mysql) - Missing${NC}"
        missing_clients="${missing_clients}mysql-client "
    fi
    
    if command -v sqlite3 >/dev/null 2>&1; then
        echo -e "${GREEN}✅ SQLite client (sqlite3) - Found${NC}"
    else
        echo -e "${YELLOW}❌ SQLite client (sqlite3) - Missing${NC}"
        missing_clients="${missing_clients}sqlite3 "
    fi
    
    if [ -n "$missing_clients" ]; then
        echo ""
        echo -e "${YELLOW}⚠️  Optional database client tools are missing:${NC} ${missing_clients}"
        echo -e "${CYAN}💡 These tools are needed for the 'spindb connect' command to open interactive shells.${NC}"
        echo ""
        
        if [ "$OS" = "linux" ]; then
            echo -e "${BOLD}Would you like to install them now? (y/N)${NC}"
            read -r install_clients
            
            if [ "$install_clients" = "y" ] || [ "$install_clients" = "Y" ]; then
                install_database_clients
            else
                echo ""
                log "💡 You can install them later with:"
                if command -v apt-get >/dev/null 2>&1; then
                    echo -e "${CYAN}  sudo apt update && sudo apt install postgresql-client mysql-client sqlite3${NC}"
                elif command -v yum >/dev/null 2>&1; then
                    echo -e "${CYAN}  sudo yum install postgresql mysql sqlite${NC}"
                elif command -v dnf >/dev/null 2>&1; then
                    echo -e "${CYAN}  sudo dnf install postgresql mysql sqlite${NC}"
                fi
            fi
        else
            echo -e "${CYAN}💡 To install them later:${NC}"
            if [ "$OS" = "darwin" ]; then
                echo -e "${CYAN}  brew install postgresql mysql-client sqlite${NC}"
            fi
        fi
    else
        success "All database client tools are installed! 🎉"
    fi
}

install_database_clients() {
    log "📦 Installing database client tools..."
    
    if command -v apt-get >/dev/null 2>&1; then
        sudo apt-get update >/dev/null 2>&1 &
        spinner $! "Updating package lists"
        
        if ! command -v psql >/dev/null 2>&1; then
            sudo apt-get install -y postgresql-client >/dev/null 2>&1 &
            spinner $! "Installing PostgreSQL client"
        fi
        
        if ! command -v mysql >/dev/null 2>&1; then
            sudo apt-get install -y mysql-client >/dev/null 2>&1 &
            spinner $! "Installing MySQL client"
        fi
        
        if ! command -v sqlite3 >/dev/null 2>&1; then
            sudo apt-get install -y sqlite3 >/dev/null 2>&1 &
            spinner $! "Installing SQLite3"
        fi
        
    elif command -v yum >/dev/null 2>&1; then
        if ! command -v psql >/dev/null 2>&1; then
            sudo yum install -y postgresql >/dev/null 2>&1 &
            spinner $! "Installing PostgreSQL client"
        fi
        
        if ! command -v mysql >/dev/null 2>&1; then
            sudo yum install -y mysql >/dev/null 2>&1 &
            spinner $! "Installing MySQL client"
        fi
        
        if ! command -v sqlite3 >/dev/null 2>&1; then
            sudo yum install -y sqlite >/dev/null 2>&1 &
            spinner $! "Installing SQLite3"
        fi
        
    elif command -v dnf >/dev/null 2>&1; then
        if ! command -v psql >/dev/null 2>&1; then
            sudo dnf install -y postgresql >/dev/null 2>&1 &
            spinner $! "Installing PostgreSQL client"
        fi
        
        if ! command -v mysql >/dev/null 2>&1; then
            sudo dnf install -y mysql >/dev/null 2>&1 &
            spinner $! "Installing MySQL client"
        fi
        
        if ! command -v sqlite3 >/dev/null 2>&1; then
            sudo dnf install -y sqlite >/dev/null 2>&1 &
            spinner $! "Installing SQLite3"
        fi
    fi
    
    success "Database client tools installation completed! 🎉"
}

show_usage() {
    echo ""
    echo -e "${GRAD1}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GRAD1}║${NC} ${BOLD}${WHITE}🎉 SpinDB Installation Complete! 🎉${NC}${GRAD1}                      ║${NC}"
    echo -e "${GRAD1}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    # Status indicators
    if command -v docker >/dev/null 2>&1 && docker ps >/dev/null 2>&1; then
        echo -e "${GREEN}🐳 Docker Status: ${BOLD}Running and Ready${NC}"
    else
        echo -e "${YELLOW}🐳 Docker Status: ${BOLD}Make sure Docker is running${NC}"
    fi
    
    # Check client tools status
    if command -v psql >/dev/null 2>&1 && command -v mysql >/dev/null 2>&1 && command -v sqlite3 >/dev/null 2>&1; then
        echo -e "${GREEN}🔧 Database Clients: ${BOLD}All Installed${NC}"
    else
        echo -e "${YELLOW}🔧 Database Clients: ${BOLD}Some Missing${NC} ${DIM}(needed for 'spindb connect')${NC}"
    fi
    
    echo ""
    echo -e "${GRAD2}┌─ Available Commands ──────────────────────────────────────────┐${NC}"
    echo -e "${GRAD2}│${NC}"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb create${NC}     ${DIM}→${NC} Create a new database instance"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb list${NC}       ${DIM}→${NC} List all database instances"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb start${NC}      ${DIM}→${NC} Start a database instance"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb stop${NC}       ${DIM}→${NC} Stop a database instance"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb connect${NC}    ${DIM}→${NC} Connect to a database instance"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb backup${NC}     ${DIM}→${NC} Backup a database instance"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb env${NC}        ${DIM}→${NC} Manage environment variables"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb template${NC}   ${DIM}→${NC} Manage database templates"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb info${NC}       ${DIM}→${NC} Show instance information"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb delete${NC}     ${DIM}→${NC} Delete a database instance"
    echo -e "${GRAD2}│${NC}  ${BLUE}spindb version${NC}    ${DIM}→${NC} Show SpinDB version"
    echo -e "${GRAD2}│${NC}"
    echo -e "${GRAD2}└───────────────────────────────────────────────────────────────┘${NC}"
    echo ""
    echo -e "${CYAN}💡 For more help, run: ${BOLD}spindb --help${NC}"
    echo ""
    echo -e "${GRAD4}🚀 ${BOLD}Happy database spinning!${NC} ${GRAD5}✨${NC}"
    echo ""
}

main() {
    # Show cool banner
    show_banner
    
    # Animated intro
    type_text "🎯 Initializing SpinDB installation process..." 0.05
    sleep 1
    
    # Show system info
    show_system_info
    
    # Start installation steps
    log "🚀 Starting installation process..."
    
    check_dependencies
    
    log "🔍 Detecting system configuration..."
    detect_os
    detect_arch
    
    check_docker
    
    get_latest_release
    
    download_binary
    
    install_binary
    
    verify_installation
    
    # Check and offer to install database client tools
    check_and_offer_client_tools
    
    # Final success animation
    echo ""
    loading_animation "🎉 Finalizing installation" 2
    
    show_usage
}

trap cleanup EXIT

# Only run main if script is executed directly (not sourced)
if [ "${BASH_SOURCE:-}" = "" ] || [ "${BASH_SOURCE[0]:-}" = "${0}" ]; then
    main "$@"
fi
