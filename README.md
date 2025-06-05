# SpinDB

A powerful CLI tool to spin up and manage databases (PostgreSQL, MySQL, SQLite) from the terminal with no web UI required. SpinDB uses Docker containers for PostgreSQL and MySQL, providing instant database instances with full lifecycle management, advanced template system, comprehensive backup/restore, and multi-environment organization.

## Quick Start

### Prerequisites
- **Go 1.21+** (for building from source)
- **Docker** (for PostgreSQL and MySQL databases)
- **Make** (optional, for development)

### Installation

#### Quick Install (Recommended)
```bash
curl -sSL https://raw.githubusercontent.com/awade12/spindb/main/install.sh | bash
```

#### From Source
```bash
git clone https://github.com/awade12/spindb.git
cd spindb
make build
sudo make install
```

> **📋 For complete installation options and troubleshooting, see [Installation Guide](docs/installation.md)**

### Basic Usage

```bash
# Create a PostgreSQL database (private, auto port)
spindb create postgres --name myapp-db --user admin --password secret123

# Create a PostgreSQL database with public access (externally accessible)
spindb create postgres --name public-api --user admin --password secret123 --public

# Create a MySQL database with specific port
spindb create mysql --name shop-db --user shopuser --password shop123 --port 3307 --public

# Create a SQLite database (file-based)
spindb create sqlite --file ./data/app.db

# Use templates for quick setup
spindb template install postgres-dev my-dev-db --public

# List all databases with status and access level
spindb list

# Get detailed database information
spindb info --name myapp-db

# Connect to database (opens database shell)
spindb connect --name myapp-db

# Start/stop database containers
spindb start --name myapp-db
spindb stop --name myapp-db
spindb restart --name myapp-db

# Create backups
spindb backup create myapp-db --compress

# Delete database and clean up containers
spindb delete --name myapp-db
```

## Features

### ✅ **Core Functionality**
- **CLI interface** with Cobra framework and comprehensive help
- **PostgreSQL databases** with full Docker container management
- **MySQL databases** with full Docker container management  
- **SQLite databases** with file-based creation and management
- **Configuration management** with persistent storage

### ✅ **Docker Integration**
- **Container lifecycle** - Create, start, stop, restart, delete
- **Smart port management** - Automatic port allocation and conflict resolution
- **Access control** - Private (localhost) or public (external) database access
- **Volume mounting** - Persistent data storage
- **Health monitoring** - Container status and database connectivity
- **Resource cleanup** - Proper container and volume cleanup

### ✅ **Database Operations**
- **Connection testing** - Verify database accessibility and health
- **Database shells** - Direct access to psql, mysql, sqlite3
- **Status monitoring** - Real-time container and database status
- **Info display** - Detailed database configuration and connection info
- **Comprehensive listing** - Show all databases with status and details

### ✅ **Security & Port Management**
- **Auto port assignment** - Automatic port allocation when not specified (--port 0)
- **Private by default** - Databases bind to localhost only for security
- **Public access option** - Use --public flag for external accessibility
- **Port conflict resolution** - Automatic detection and assignment of available ports
- **Access level display** - Clear indication of private vs public database access

### ✅ **Database Templates** (New in Phase 3)
- **Predefined templates** - Common database setups (dev, test, prod)
- **Custom templates** - Create and share custom database configurations
- **Template management** - Create, list, show, delete, import, export templates
- **Quick deployment** - Install databases from templates with override options
- **Version presets** - Templates for specific database versions and configurations

### ✅ **Backup & Restore** (New in Phase 3)
- **Multi-database support** - Backup PostgreSQL, MySQL, and SQLite databases
- **Backup options** - Full, schema-only, data-only backup modes
- **Compression** - Optional gzip compression for space efficiency
- **Backup management** - List, restore, and delete backup files
- **Cross-platform** - Compatible backup formats across different systems

### ✅ **Environment Management** (New in Phase 3)
- **Environment profiles** - Create isolated environments (dev, staging, production)
- **Database organization** - Add/remove databases from environments
- **Environment switching** - Quick context switching between environments
- **Bulk operations** - Start, stop, restart all databases in an environment
- **Isolation controls** - Isolate and activate entire environments

### 🔄 **Coming Soon** (Phase 4)
- Cloud integration (AWS RDS, GCP, Azure)
- Performance monitoring and metrics
- Team collaboration features
- Enterprise security and audit logging

## Requirements

### System Dependencies
- **Docker** - Required for PostgreSQL and MySQL databases
  ```bash
  # Install Docker (varies by OS)
  # macOS: brew install --cask docker
  # Ubuntu: sudo apt install docker.io
  # See: https://docs.docker.com/get-docker/
  ```

### Database Client Tools (Optional)
For the `spindb connect` command to open interactive database shells:

- **PostgreSQL**: `psql` client
  ```bash
  # Ubuntu/Debian
  sudo apt install postgresql-client
  
  # CentOS/RHEL
  sudo yum install postgresql
  
  # macOS
  brew install postgresql
  ```

- **MySQL**: `mysql` client
  ```bash
  # Ubuntu/Debian
  sudo apt install mysql-client
  
  # CentOS/RHEL
  sudo yum install mysql
  
  # macOS
  brew install mysql-client
  ```

- **SQLite**: `sqlite3` (usually pre-installed)
  ```bash
  # Ubuntu/Debian
  sudo apt install sqlite3
  
  # CentOS/RHEL
  sudo yum install sqlite
  
  # macOS
  brew install sqlite
  ```

> **📝 Note**: The enhanced install script can automatically install these client tools on Linux systems. If client tools are missing, SpinDB will provide helpful error messages with installation instructions and Docker alternatives.

## Development

### Setup
```bash
git clone https://github.com/awade12/spindb.git
cd spindb
go mod download
make dev-setup
```

### Building
```bash
make build           # Build for current platform
make build-all       # Build for all platforms
make test           # Run tests
make lint           # Run linter
make install        # Install to system
```

### Project Structure
```
spindb/
├── cmd/                 # CLI commands (Cobra)
├── internal/
│   ├── backup/         # Backup and restore management
│   ├── config/         # Configuration and template management
│   ├── db/             # Database management logic
│   ├── docker/         # Docker service integration
│   ├── environment/    # Environment management
│   └── utils/          # Utility functions
├── docs/               # Documentation
├── Makefile           # Build and development tasks
└── main.go            # Application entry point
```

## Examples

### Complete PostgreSQL Workflow
```bash
# Create private PostgreSQL database (localhost only, auto port)
spindb create postgres --name webapp-db --user webuser --password secure123

# Create public PostgreSQL database (externally accessible)
spindb create postgres --name api-db --user apiuser --password secure456 --public

# Check status (shows access level and port)
spindb list

# Get connection details with credentials
spindb info --name webapp-db --show-credentials

# Connect to database
spindb connect --name webapp-db

# Create backup
spindb backup create webapp-db --compress

# Stop database when not needed
spindb stop --name webapp-db

# Start again when needed
spindb start --name webapp-db

# Clean up when done
spindb delete --name webapp-db
```

### Template-Based Development Workflow
```bash
# List available templates
spindb template list

# Create custom template
spindb template create --name my-postgres --type postgres --description "My custom setup" \
  --version 15 --user myuser --password mypass --port 5432 --tags dev,custom

# Install database from template (private by default)
spindb template install my-postgres my-app-db

# Override template settings with public access
spindb template install postgres-dev another-db --password different-pass --public

# Install with auto port assignment
spindb template install postgres-dev auto-port-db --port 0

# Export template for sharing
spindb template export my-postgres my-template.yaml

# Import template
spindb template import shared-template.yaml
```

### Environment Management Workflow
```bash
# Create environments
spindb env create development --description "Development environment"
spindb env create staging --description "Staging environment"

# Install databases
spindb template install postgres-dev dev-api-db
spindb template install mysql-dev dev-analytics-db

# Add databases to development environment
spindb env add development dev-api-db
spindb env add development dev-analytics-db

# Switch to development environment
spindb env switch development

# Check environment status
spindb env show development

# Start all databases in environment
spindb env bulk start development

# Create staging database
spindb template install postgres-test staging-api-db --port 5433
spindb env add staging staging-api-db

# Switch environments
spindb env switch staging

# Isolate development environment
spindb env isolate development

# Clean up
spindb env delete development --force
```

### Backup and Restore Workflow
```bash
# Create various types of backups
spindb backup create my-db                    # Full backup
spindb backup create my-db --compress         # Compressed backup
spindb backup create my-db --schema-only      # Schema only
spindb backup create my-db --data-only        # Data only

# List all backups
spindb backup list

# Restore backup to new database
spindb backup restore my-db_20250604_141922.sql.gz target-db

# Clean up old backups
spindb backup delete old-backup.sql
```

### Development Workflow with All Features
```bash
# Setup development environment
spindb env create development --description "Development environment"

# Create custom template for project
spindb template create --name myproject-db --type postgres --description "My project database" \
  --version 15 --user projectuser --password projectpass --tags project,dev

# Install development databases
spindb template install myproject-db dev-api
spindb template install myproject-db dev-worker --port 5433

# Add to development environment
spindb env add development dev-api
spindb env add development dev-worker

# Switch to development context
spindb env switch development

# Work with databases...
spindb list
spindb info --name dev-api

# Create backup before major changes
spindb backup create dev-api --compress

# Bulk operations
spindb env bulk stop development      # Stop all dev databases
spindb env bulk start development     # Start all dev databases

# Clean up when done
spindb env isolate development
spindb env delete development --force
```

## Command Reference

### Core Commands
- `spindb create {postgres|mysql|sqlite}` - Create database instances
  - `--port 0` for auto port assignment
  - `--public` for external access (private by default)
- `spindb list` - List all managed databases with access levels
- `spindb info --name <db>` - Show database details including access level
- `spindb connect --name <db>` - Connect to database
- `spindb {start|stop|restart} --name <db>` - Control database lifecycle
- `spindb delete --name <db>` - Delete database and cleanup

### Template Commands
- `spindb template list` - List all available templates
- `spindb template create` - Create custom template
- `spindb template show <name>` - Show template details
- `spindb template install <template> <db-name>` - Create database from template
  - `--public` for external access override
  - `--port <number>` for port override
- `spindb template {import|export} <name> <file>` - Share templates
- `spindb template delete <name>` - Delete custom template

### Backup Commands
- `spindb backup create <db>` - Create database backup
- `spindb backup list` - List all backups
- `spindb backup restore <backup> <target-db>` - Restore backup
- `spindb backup delete <backup>` - Delete backup file

### Environment Commands
- `spindb env create <name>` - Create new environment
- `spindb env list` - List all environments
- `spindb env switch <name>` - Switch to environment
- `spindb env show <name>` - Show environment details
- `spindb env {add|remove} <env> <db>` - Manage databases in environment
- `spindb env bulk {start|stop|restart} <env>` - Bulk operations
- `spindb env {isolate|activate} <env>` - Environment isolation
- `spindb env delete <name>` - Delete environment

## Security Best Practices

SpinDB prioritizes security with sensible defaults:

### 🔒 **Access Control**
- **Private by default** - All databases bind to `127.0.0.1` (localhost only)
- **Explicit public access** - Use `--public` flag only when external access is needed
- **Port isolation** - Auto-assigned ports reduce conflicts and exposure

### 🛡️ **Network Security**
```bash
# ✅ Secure (private database)
spindb create postgres --name secure-db --password strongpass123

# ⚠️ Use with caution (public database)
spindb create postgres --name api-db --password strongpass123 --public

# 🔍 Check access levels
spindb list  # Shows "Private" or "Public" for each database
```

### 🔑 **Connection Examples**
```bash
# Private database - only accessible from localhost
psql -h localhost -p 5432 -U user -d mydb

# Public database - accessible from external networks
psql -h your-server-ip -p 5432 -U user -d mydb
```

> **💡 Tip**: Always use strong passwords and consider firewall rules when using `--public` flag in production environments.

## Contributing

We welcome contributions! Please see our [Roadmap](docs/roadmap.md) for current priorities and planned features.

1. **Bug Reports** - Submit issues for bugs or unexpected behavior
2. **Feature Requests** - Suggest new features or improvements
3. **Code Contributions** - Submit PRs for bug fixes or features
4. **Documentation** - Improve docs, examples, and guides

## License

MIT License - see LICENSE file for details. 