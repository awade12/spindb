# SpinDB

A powerful CLI tool to spin up and manage databases (PostgreSQL, MySQL, SQLite) from the terminal with no web UI required.

## Project Structure

```text
spindb/                     # Project root
├── cmd/                    # CLI command definitions
│   ├── root.go             # Root command (spindb)
│   ├── create.go           # spindb create ...
│   ├── list.go             # spindb list
│   ├── connect.go          # spindb connect ...
│   ├── info.go             # spindb info ...
│   ├── delete.go           # spindb delete ...
│   └── version.go          # spindb version
├── internal/               # Internal packages
│   ├── db/                 # Database management logic
│   │   ├── postgres.go     # PostgreSQL operations
│   │   ├── mysql.go        # MySQL operations
│   │   ├── sqlite.go       # SQLite operations
│   │   └── manager.go      # Database manager interface
│   ├── config/             # Configuration management
│   │   ├── config.go       # Config struct and loading
│   │   └── database.go     # Database configuration
│   ├── docker/             # Docker management (for containerized DBs)
│   │   ├── client.go       # Docker client wrapper
│   │   └── compose.go      # Docker compose operations
│   └── utils/              # Utility functions
│       ├── validation.go   # Input validation
│       └── helpers.go      # General helpers
├── configs/                # Configuration files
│   ├── spindb.yaml         # Default configuration
│   └── databases.yaml      # Database templates
├── scripts/                # Build and deployment scripts
│   ├── build.sh           # Build script
│   └── install.sh         # Installation script
├── .gitignore             # Git ignore rules
├── .goreleaser.yaml       # GoReleaser configuration
├── Dockerfile             # Docker image for the CLI
├── docker-compose.yml     # Docker compose for development
├── go.mod                 # Go module file
├── go.sum                 # Go module checksums
├── main.go                # Program entry point
├── Makefile               # Build automation
└── README.md              # Project documentation
```

## Tech Stack

- **CLI Framework**: [Cobra](https://github.com/spf13/cobra) - Standard for Go CLI applications
- **Configuration**: [Viper](https://github.com/spf13/viper) - Configuration management with multiple formats
- **Database Drivers**: 
  - PostgreSQL: `github.com/lib/pq`
  - MySQL: `github.com/go-sql-driver/mysql`
  - SQLite: `github.com/mattn/go-sqlite3`
- **Container Management**: Docker SDK for Go
- **Testing**: Standard Go testing + Testify

## Command Structure

### Database Creation
```bash
# PostgreSQL
spindb create postgres --name myapp-db --user admin --password secret123 --port 5432 --public

# MySQL  
spindb create mysql --name shop-db --user shopuser --password shop123 --port 3306

# SQLite
spindb create sqlite --file ./data/app.db
```

### Database Management
```bash
# List all managed databases
spindb list
spindb list --type postgres  # Filter by database type

# Get database information
spindb info --name myapp-db
spindb info --name shop-db --show-credentials

# Test database connection
spindb connect --name myapp-db
spindb connect --name shop-db --test-only

# Database operations
spindb start --name myapp-db
spindb stop --name myapp-db
spindb restart --name myapp-db

# Cleanup
spindb delete --name myapp-db
spindb delete --name shop-db --force
spindb delete --file ./data/app.db
```

### Configuration and Utilities
```bash
# Configuration management
spindb config init          # Initialize configuration
spindb config show          # Show current configuration
spindb config set key value # Set configuration value

# Utility commands
spindb version              # Show version information
spindb health              # Check system health
spindb cleanup             # Clean up orphaned resources
```

## Database Support

### PostgreSQL
- **Versions**: 12, 13, 14, 15, 16
- **Container**: Official PostgreSQL Docker images
- **Default Port**: 5432
- **Features**: Full PostgreSQL feature set, extensions support

### MySQL
- **Versions**: 5.7, 8.0, 8.1
- **Container**: Official MySQL Docker images  
- **Default Port**: 3306
- **Features**: Full MySQL feature set, configuration presets

### SQLite
- **Version**: Latest (via go-sqlite3)
- **Storage**: File-based, no container required
- **Features**: Embedded database, perfect for development

## Configuration

SpinDB uses a hierarchical configuration system:

1. **Command-line flags** (highest priority)
2. **Environment variables** (prefixed with `SPINDB_`)
3. **Configuration files** (`~/.spindb/config.yaml`, `./spindb.yaml`)
4. **Defaults** (lowest priority)

### Sample Configuration (`~/.spindb/config.yaml`)
```yaml
default:
  postgres:
    version: "15"
    port: 5432
    user: "postgres"
  mysql:
    version: "8.0"
    port: 3306
    user: "root"
  
docker:
  host: "unix:///var/run/docker.sock"
  cleanup_timeout: "30s"
  
storage:
  data_dir: "~/.spindb/data"
  backup_dir: "~/.spindb/backups"
```

## Installation

### From Release (Recommended)
```bash
curl -sSfL https://raw.githubusercontent.com/yourusername/spindb/main/scripts/install.sh | sh
```

### From Source
```bash
git clone https://github.com/yourusername/spindb.git
cd spindb
make build
sudo make install
```

### Using Go
```bash
go install github.com/yourusername/spindb@latest
```

## Development Setup

### Prerequisites
- Go 1.21+
- Docker and Docker Compose
- Make (optional, for build automation)

### Setup
```bash
git clone https://github.com/yourusername/spindb.git
cd spindb
go mod download
make dev-setup  # Sets up development dependencies
```

### Building
```bash
make build           # Build for current platform
make build-all       # Build for all platforms
make test           # Run tests
make lint           # Run linter
```

## Features Roadmap

### Phase 1 (MVP)
- [ ] Basic CLI structure with Cobra
- [ ] PostgreSQL container management
- [ ] MySQL container management  
- [ ] SQLite file management
- [ ] Database connection testing
- [ ] Basic configuration system

### Phase 2 (Enhanced)
- [ ] Database templates and presets
- [ ] Backup and restore functionality
- [ ] Database migrations support
- [ ] Environment-specific configurations
- [ ] Integration with popular ORMs

### Phase 3 (Advanced)
- [ ] Database monitoring and metrics
- [ ] Cluster management (multi-node)
- [ ] Cloud provider integration (AWS RDS, etc.)
- [ ] Database sharing and collaboration
- [ ] Web dashboard (optional)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.