# SpinDB Changelog

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