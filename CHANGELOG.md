# SpinDB Changelog

## [Unreleased] - 2025-01-06

### 🔧 **Improved Version Management**
- **Dynamic Version Injection**: Replaced hardcoded version string with build-time injection using Git tags and ldflags
- **Git Integration**: Version now automatically determined from Git repository state using `git describe --tags`
- **Build Information**: Enhanced version command to show Git commit hash and build timestamp
- **Professional Versioning**: Adopted industry-standard approach used by most Go projects

### ✅ **Enhanced**
- **Makefile Updates**: Added automatic version detection from Git tags and commits
- **Build Process**: Enhanced build commands with dynamic version, commit hash, and build date injection
- **Version Command**: Now displays comprehensive build information including Git commit and build timestamp
- **Development Workflow**: No more manual version updates required - versions automatically track Git tags

### 🛠️ **Technical Changes**
- **Build-time Variables**: Added `Version`, `GitCommit`, and `BuildDate` variables to `cmd/version.go`
- **Makefile Enhancement**: Added shell commands to extract Git information and inject via ldflags
- **Automatic Versioning**: Version now follows pattern like `v1.2-1-g359cec9-dirty` based on Git state
- **Cross-platform Builds**: All platform builds now include proper version information

### 💡 **Usage Examples**
```bash
# Build with automatic Git-based versioning
make build

# Build with custom version override
VERSION=v0.3.0-beta make build

# Check comprehensive version info
./build/spindb version
# Output:
# SpinDB v1.2-1-g359cec9-dirty
# Git commit: 359cec9
# Built: 2025-06-05_01:53:43
```

### Breaking Changes
- None. Version command output format enhanced but remains backward compatible.

## [v1.2.0] - 2025-01-06

### 🚀 **New Features**
- **Auto Port Assignment**: Databases now auto-assign ports by default (`--port 0`)
- **Public Access Control**: Added `--public` flag for external database access
- **Security by Default**: Databases are now private (localhost only) by default
- **Access Level Display**: List and info commands now show Public/Private access status

### ✅ **Added**
- **MySQL Public Access**: Added `--public` flag support for MySQL databases
- **Smart Port Management**: Automatic port conflict resolution with fallback allocation
- **Template Public Override**: Template installation now supports `--public` flag override
- **Enhanced Database Info**: Show access level (Public/Private) in database information
- **Connection Examples**: Dynamic connection strings based on access level (localhost vs external IP)

### 🔒 **Security Improvements**
- **Private by Default**: All databases bind to `127.0.0.1` (localhost) unless explicitly made public
- **Explicit Public Access**: Requires `--public` flag for external accessibility
- **Port Isolation**: Auto-assigned ports reduce conflicts and accidental exposure
- **Security Documentation**: Added comprehensive security best practices section

### 🛠️ **Enhanced**
- **CLI Flags**: Updated port flags with clear help text (`--port 0` for auto)
- **Output Messages**: Enhanced creation success messages with access level information
- **Command Help**: Improved flag descriptions for better user experience
- **Template Defaults**: Templates now default to auto port assignment for better flexibility

### 📚 **Documentation**
- **README Overhaul**: Comprehensive updates showcasing new security and port features
- **Security Guide**: New section with best practices and connection examples
- **Command Reference**: Updated with new flags and options
- **Examples**: Enhanced workflows demonstrating auto ports and access control

### 🔧 **Technical Changes**
- **Docker Integration**: Modified container port binding based on public/private setting
- **Database Config**: Extended `MySQLConfig` with `Public` field to match PostgreSQL
- **Container Config**: Added `Public` field to `ContainerConfig` struct
- **Manager Logic**: Enhanced database creation with proper access control handling
- **Template System**: Updated template installation to support access control overrides

### 💡 **Usage Examples**
```bash
# Private database (default, localhost only)
spindb create postgres --name mydb --password secret123

# Public database (externally accessible)
spindb create postgres --name api-db --password secret123 --public

# Auto port assignment (default)
spindb create mysql --name shop --password secret123

# Template with public override
spindb template install postgres-dev my-app --public
```

### Breaking Changes
- **Default Ports**: CLI flags now default to `0` (auto) instead of specific ports
- **Network Binding**: Databases now bind to localhost by default (was previously all interfaces)
- **Template Behavior**: Templates now use auto port assignment unless explicitly specified

> **⚠️ Migration Note**: Existing databases created before v1.2.0 may still bind to all interfaces. Recreate databases if localhost-only access is required.

## [Unreleased] - 2025-01-06

### Fixed
- **Database Connect Command**: Fixed issue where `spindb connect` would fail with "executable file not found in $PATH" error when database client tools (psql, mysql, sqlite3) were not installed
- Enhanced error handling to provide helpful installation instructions for missing client tools
- Added alternative Docker-based connection commands when client tools are unavailable

### Enhanced
- **Install Script**: Added automatic detection and optional installation of database client tools during SpinDB installation
- **Error Messages**: Improved error messages with platform-specific installation instructions for PostgreSQL, MySQL, and SQLite clients
- **Documentation**: Updated README with detailed client tool installation instructions and requirements

### Technical Details
- Modified `openConnection` method in `internal/db/manager.go` to check for client tool availability before attempting to execute
- Added `getInstallInstructions` helper function with platform-specific package manager commands
- Enhanced install script with `check_and_offer_client_tools` and `install_database_clients` functions
- Added fallback Docker commands for each database type when native clients are unavailable

### Breaking Changes
- None. All existing functionality remains the same.

## [v1.1.0] - Previous Release
- Full Phase 3 features: Database templates, backup/restore, environment management
- Docker integration with PostgreSQL, MySQL, and SQLite support
- Comprehensive CLI interface with Cobra framework 