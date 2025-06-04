# SpinDB Installation Guide

This guide covers all the different ways to install SpinDB on your system.

## 🚀 Quick Installation (Recommended)

### One-Command Install

The fastest way to install SpinDB is using our installation script:

```bash
curl -sSL https://raw.githubusercontent.com/awade12/spindb/main/install.sh | bash
```

**What this does:**
- Automatically detects your operating system and architecture
- Downloads the appropriate SpinDB binary from GitHub releases
- Installs Docker if not already present
- Installs SpinDB to `/usr/local/bin` (globally accessible)
- Verifies the installation and shows available commands

**Supported Platforms:**
- **Linux**: amd64, arm64 (Ubuntu, Debian, CentOS, Fedora, and others)
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64 (via WSL, Git Bash, or MSYS2)

---

## 📦 Installation Methods

### Method 1: Install Script (Recommended)

#### Basic Installation
```bash
curl -sSL https://raw.githubusercontent.com/awade12/spindb/main/install.sh | bash
```

#### Download and Inspect First
If you prefer to review the script before running:
```bash
curl -O https://raw.githubusercontent.com/awade12/spindb/main/install.sh
chmod +x install.sh
./install.sh
```

#### What Gets Installed
- **SpinDB Binary**: Latest release from GitHub
- **Docker**: Automatically installed if not present
- **Dependencies**: All required system dependencies

---

### Method 2: Building from Source

#### Prerequisites
- **Go 1.21+** ([Install Go](https://golang.org/doc/install))
- **Git**
- **Make** (optional, but recommended)
- **Docker** (for PostgreSQL/MySQL support)

#### Step-by-Step Build
```bash
# Clone the repository
git clone https://github.com/awade12/spindb.git
cd spindb

# Download dependencies
go mod download

# Build SpinDB
make build

# Install globally (optional)
sudo make install
```

#### Manual Build (without Make)
```bash
# Build for current platform
go build -ldflags "-s -w" -o build/spindb .

# Install to system PATH
sudo cp build/spindb /usr/local/bin/
```

#### Cross-Platform Building
```bash
# Build for all platforms
make build-all

# This creates:
# - build/spindb-linux-amd64
# - build/spindb-darwin-amd64
# - build/spindb-darwin-arm64
# - build/spindb-windows-amd64.exe
```

---

### Method 3: GitHub Releases (Manual)

#### Download Pre-built Binaries
1. Visit [SpinDB Releases](https://github.com/awade12/spindb/releases)
2. Download the appropriate binary for your platform:
   - `spindb-linux-amd64` (Linux 64-bit)
   - `spindb-darwin-amd64` (macOS Intel)
   - `spindb-darwin-arm64` (macOS Apple Silicon)
   - `spindb-windows-amd64.exe` (Windows 64-bit)

#### Install Downloaded Binary
```bash
# Make executable (Linux/macOS)
chmod +x spindb-*

# Move to PATH
sudo mv spindb-* /usr/local/bin/spindb

# Verify installation
spindb version
```

---

### Method 4: Package Managers

#### Homebrew (macOS/Linux)
```bash
# Coming soon...
brew install awade12/tap/spindb
```

#### APT (Ubuntu/Debian)
```bash
# Coming soon...
curl -fsSL https://packages.spindb.io/gpg | sudo apt-key add -
echo "deb https://packages.spindb.io/apt stable main" | sudo tee /etc/apt/sources.list.d/spindb.list
sudo apt update
sudo apt install spindb
```

#### YUM/DNF (RHEL/CentOS/Fedora)
```bash
# Coming soon...
sudo yum-config-manager --add-repo https://packages.spindb.io/yum/spindb.repo
sudo yum install spindb
```

---

## 🐳 Docker Requirements

SpinDB uses Docker to manage PostgreSQL and MySQL databases. The install script automatically handles Docker installation, but you can also install it manually:

### Docker Installation

#### Linux
```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# CentOS/RHEL/Fedora
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
```

#### macOS
```bash
# Download Docker Desktop for Mac
# Intel: https://desktop.docker.com/mac/main/amd64/Docker.dmg
# Apple Silicon: https://desktop.docker.com/mac/main/arm64/Docker.dmg

# Or with Homebrew
brew install --cask docker
```

#### Windows
Download Docker Desktop from: https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe

### Verify Docker Installation
```bash
docker --version
docker ps
```

---

## 🔧 Platform-Specific Instructions

### Linux

#### Ubuntu/Debian
```bash
# Update package list
sudo apt update

# Install dependencies
sudo apt install -y curl git

# Install SpinDB
curl -sSL https://raw.githubusercontent.com/awade12/spindb/main/install.sh | bash
```

#### CentOS/RHEL/Fedora
```bash
# Install dependencies
sudo yum install -y curl git

# Install SpinDB
curl -sSL https://raw.githubusercontent.com/awade12/spindb/main/install.sh | bash
```

#### Arch Linux
```bash
# Install dependencies
sudo pacman -S curl git

# Install SpinDB
curl -sSL https://raw.githubusercontent.com/awade12/spindb/main/install.sh | bash
```

### macOS

#### Prerequisites
- **Xcode Command Line Tools**: `xcode-select --install`
- **Homebrew** (recommended): `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`

#### Installation
```bash
# Install SpinDB
curl -sSL https://raw.githubusercontent.com/awade12/spindb/main/install.sh | bash

# Or build from source
brew install go
git clone https://github.com/awade12/spindb.git
cd spindb && make build && sudo make install
```

### Windows

#### WSL2 (Recommended)
```bash
# In WSL2 terminal
curl -sSL https://raw.githubusercontent.com/awade12/spindb/main/install.sh | bash
```

#### Git Bash/MSYS2
```bash
# In Git Bash
curl -sSL https://raw.githubusercontent.com/awade12/spindb/main/install.sh | bash
```

#### PowerShell (Manual)
```powershell
# Download binary from GitHub releases
Invoke-WebRequest -Uri "https://github.com/awade12/spindb/releases/latest/download/spindb-windows-amd64.exe" -OutFile "spindb.exe"

# Move to PATH location
Move-Item spindb.exe C:\Windows\System32\
```

---

## 🛠️ Development Installation

### For Contributors and Developers

#### Setup Development Environment
```bash
# Clone repository
git clone https://github.com/awade12/spindb.git
cd spindb

# Install development dependencies
make dev-setup

# This installs:
# - golangci-lint (for linting)
# - Other development tools
```

#### Available Make Targets
```bash
make build           # Build for current platform
make build-all       # Build for all platforms
make test           # Run tests
make lint           # Run linter
make clean          # Clean build artifacts
make install        # Install to system
make dev-setup      # Setup development environment
```

#### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific package tests
go test -v ./internal/db
```

---

## ✅ Verification and Testing

### Verify Installation
```bash
# Check version
spindb version

# Check help
spindb --help

# List available commands
spindb help
```

### Test Docker Integration
```bash
# Check Docker status
docker --version
docker ps

# Test SpinDB with Docker
spindb create postgres --name test-db --user testuser --password testpass
spindb list
spindb delete --name test-db --force
```

### Test SQLite Support
```bash
# Create SQLite database
spindb create sqlite --file ./test.db

# List databases
spindb list

# Connect to database
spindb connect --name test.db

# Clean up
rm test.db
```

---

## 🚨 Troubleshooting

### Common Issues

#### "command not found: spindb"
```bash
# Check if binary is in PATH
which spindb
ls -la /usr/local/bin/spindb

# Add to PATH if needed
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc
```

#### "Docker daemon not running"
```bash
# Start Docker (Linux)
sudo systemctl start docker

# Start Docker (macOS)
open /Applications/Docker.app

# Check Docker status
docker ps
```

#### "Permission denied" during installation
```bash
# Install to user directory instead
mkdir -p ~/bin
curl -L -o ~/bin/spindb https://github.com/awade12/spindb/releases/latest/download/spindb-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)
chmod +x ~/bin/spindb
echo 'export PATH=$PATH:~/bin' >> ~/.bashrc
```

#### "Failed to pull Docker image"
```bash
# Check internet connection
ping google.com

# Check Docker Hub access
docker pull hello-world

# Try with different registry (if needed)
docker pull postgres:15
```

### Platform-Specific Issues

#### Linux: "snap docker" conflicts
```bash
# Remove snap docker
sudo snap remove docker

# Install Docker CE
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

#### macOS: "Docker Desktop not starting"
```bash
# Reset Docker Desktop
rm -rf ~/Library/Group\ Containers/group.com.docker
rm -rf ~/Library/Containers/com.docker.docker
open /Applications/Docker.app
```

#### Windows: WSL Docker issues
```bash
# In PowerShell (as Administrator)
wsl --update
wsl --shutdown

# Start Docker Desktop
# Ensure "Use WSL 2 based engine" is enabled
```

---

## 🔄 Updating SpinDB

### Update via Install Script
```bash
# Re-run install script to get latest version
curl -sSL https://raw.githubusercontent.com/awade12/spindb/main/install.sh | bash
```

### Manual Update
```bash
# Download latest binary
curl -L -o /tmp/spindb https://github.com/awade12/spindb/releases/latest/download/spindb-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)

# Replace existing binary
sudo mv /tmp/spindb /usr/local/bin/spindb
sudo chmod +x /usr/local/bin/spindb

# Verify update
spindb version
```

### Build Latest from Source
```bash
cd spindb
git pull origin main
make build
sudo make install
```

---

## 🗑️ Uninstallation

### Remove SpinDB Binary
```bash
# Remove binary
sudo rm /usr/local/bin/spindb

# Remove configuration (optional)
rm -rf ~/.spindb
```

### Remove Docker (if installed by SpinDB)
```bash
# Linux
sudo apt remove docker-ce docker-ce-cli containerd.io
# or
sudo yum remove docker-ce docker-ce-cli containerd.io

# macOS
rm -rf /Applications/Docker.app
```

### Clean Up SpinDB Data
```bash
# Remove all SpinDB databases and configs
rm -rf ~/.spindb

# Remove Docker containers created by SpinDB
docker container prune -f
docker image prune -f
```

---

## 📞 Support and Community

### Get Help
- **Documentation**: [README.md](../README.md)
- **GitHub Issues**: [Report bugs or request features](https://github.com/awade12/spindb/issues)
- **Discussions**: [Community support](https://github.com/awade12/spindb/discussions)

### Contributing
- **Roadmap**: [docs/roadmap.md](roadmap.md)
- **Contributing Guide**: [CONTRIBUTING.md](../CONTRIBUTING.md)
- **Development**: See [Development Installation](#development-installation)

---

*Last updated: January 2025*  
*For the latest installation instructions, visit: https://github.com/awade12/spindb* 