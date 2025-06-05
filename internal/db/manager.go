package db

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/awade12/spindb/internal/config"
	"github.com/awade12/spindb/internal/docker"
)

type DatabaseManager interface {
	CreatePostgres(cfg *PostgresConfig) error
	CreateMySQL(cfg *MySQLConfig) error
	CreateSQLite(cfg *SQLiteConfig) error
	ListDatabases(dbType string) error
	Connect(name string, testOnly bool) error
	ShowInfo(name string, showCredentials bool) error
	Delete(name, file string, force bool) error
	Start(name string) error
	Stop(name string) error
	Restart(name string) error
}

type Manager struct {
	config        *config.Config
	dockerService *docker.Service
	store         *config.DatabaseStore
	connTester    *ConnectionTester
}

func NewManager() *Manager {
	cfg := config.Load()
	dockerSvc, _ := docker.NewService()
	store := config.NewDatabaseStore()
	connTester := NewConnectionTester()

	return &Manager{
		config:        cfg,
		dockerService: dockerSvc,
		store:         store,
		connTester:    connTester,
	}
}

func (m *Manager) CreatePostgres(cfg *PostgresConfig) error {
	if m.dockerService == nil {
		return fmt.Errorf("docker service not available")
	}

	if err := m.dockerService.IsDockerRunning(); err != nil {
		return fmt.Errorf("docker is not running: %w", err)
	}

	port := cfg.Port
	if port == 0 {
		port = m.config.Default.Postgres.Port
	}

	availablePort, err := m.dockerService.FindAvailablePort(port)
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}

	containerName := fmt.Sprintf("spindb-postgres-%s", cfg.Name)
	image := fmt.Sprintf("postgres:%s", cfg.Version)

	ctx := context.Background()

	fmt.Printf("Pulling PostgreSQL image %s...\n", image)
	if err := m.dockerService.PullImage(ctx, image); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	dataDir := filepath.Join(m.config.Storage.DataDir, "postgres", cfg.Name)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	containerConfig := &docker.ContainerConfig{
		Name:  containerName,
		Image: image,
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", cfg.Name),
			fmt.Sprintf("POSTGRES_USER=%s", cfg.User),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.Password),
		},
		Ports: map[string]string{
			"5432": strconv.Itoa(availablePort),
		},
		Volumes: []string{
			docker.CreateVolumeMount(dataDir, "/var/lib/postgresql/data"),
		},
		Public: cfg.Public,
	}

	fmt.Printf("Creating PostgreSQL container %s...\n", containerName)
	containerID, err := m.dockerService.CreateContainer(ctx, containerConfig)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	fmt.Printf("Starting PostgreSQL container...\n")
	if err := m.dockerService.StartContainer(ctx, containerID); err != nil {
		m.dockerService.RemoveContainer(ctx, containerID, true)
		return fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("Waiting for PostgreSQL to be ready...\n")
	dsn := fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable",
		availablePort, cfg.User, cfg.Password, cfg.Name)
	if err := m.connTester.WaitForDatabase("postgres", dsn, 60*time.Second); err != nil {
		return fmt.Errorf("PostgreSQL failed to start: %w", err)
	}

	dbConfig := &config.DatabaseConfig{
		Name:        cfg.Name,
		Type:        "postgres",
		Version:     cfg.Version,
		Port:        availablePort,
		User:        cfg.User,
		Password:    cfg.Password,
		Public:      cfg.Public,
		ContainerID: containerID,
		Created:     time.Now(),
	}

	if err := m.store.Save(dbConfig); err != nil {
		return fmt.Errorf("failed to save database config: %w", err)
	}

	fmt.Printf("✅ PostgreSQL database '%s' created successfully!\n", cfg.Name)
	fmt.Printf("   Container ID: %s\n", containerID[:12])
	fmt.Printf("   Port: %d\n", availablePort)
	host := "localhost"
	if cfg.Public {
		host = "<your-server-ip>"
		fmt.Printf("   Public: Yes (accessible externally)\n")
	} else {
		fmt.Printf("   Public: No (localhost only)\n")
	}
	fmt.Printf("   Connection: psql -h %s -p %d -U %s -d %s\n", host, availablePort, cfg.User, cfg.Name)

	return nil
}

func (m *Manager) CreateMySQL(cfg *MySQLConfig) error {
	if m.dockerService == nil {
		return fmt.Errorf("docker service not available")
	}

	if err := m.dockerService.IsDockerRunning(); err != nil {
		return fmt.Errorf("docker is not running: %w", err)
	}

	port := cfg.Port
	if port == 0 {
		port = m.config.Default.MySQL.Port
	}

	availablePort, err := m.dockerService.FindAvailablePort(port)
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}

	containerName := fmt.Sprintf("spindb-mysql-%s", cfg.Name)
	image := fmt.Sprintf("mysql:%s", cfg.Version)

	ctx := context.Background()

	fmt.Printf("Pulling MySQL image %s...\n", image)
	if err := m.dockerService.PullImage(ctx, image); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	dataDir := filepath.Join(m.config.Storage.DataDir, "mysql", cfg.Name)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	containerConfig := &docker.ContainerConfig{
		Name:  containerName,
		Image: image,
		Env: []string{
			fmt.Sprintf("MYSQL_DATABASE=%s", cfg.Name),
			fmt.Sprintf("MYSQL_USER=%s", cfg.User),
			fmt.Sprintf("MYSQL_PASSWORD=%s", cfg.Password),
			fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", cfg.Password),
		},
		Ports: map[string]string{
			"3306": strconv.Itoa(availablePort),
		},
		Volumes: []string{
			docker.CreateVolumeMount(dataDir, "/var/lib/mysql"),
		},
		Public: cfg.Public,
	}

	fmt.Printf("Creating MySQL container %s...\n", containerName)
	containerID, err := m.dockerService.CreateContainer(ctx, containerConfig)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	fmt.Printf("Starting MySQL container...\n")
	if err := m.dockerService.StartContainer(ctx, containerID); err != nil {
		m.dockerService.RemoveContainer(ctx, containerID, true)
		return fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("Waiting for MySQL to be ready...\n")
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:%d)/%s", cfg.User, cfg.Password, availablePort, cfg.Name)
	if err := m.connTester.WaitForDatabase("mysql", dsn, 60*time.Second); err != nil {
		return fmt.Errorf("MySQL failed to start: %w", err)
	}

	dbConfig := &config.DatabaseConfig{
		Name:        cfg.Name,
		Type:        "mysql",
		Version:     cfg.Version,
		Port:        availablePort,
		User:        cfg.User,
		Password:    cfg.Password,
		Public:      cfg.Public,
		ContainerID: containerID,
		Created:     time.Now(),
	}

	if err := m.store.Save(dbConfig); err != nil {
		return fmt.Errorf("failed to save database config: %w", err)
	}

	fmt.Printf("✅ MySQL database '%s' created successfully!\n", cfg.Name)
	fmt.Printf("   Container ID: %s\n", containerID[:12])
	fmt.Printf("   Port: %d\n", availablePort)
	host := "localhost"
	if cfg.Public {
		host = "<your-server-ip>"
		fmt.Printf("   Public: Yes (accessible externally)\n")
	} else {
		fmt.Printf("   Public: No (localhost only)\n")
	}
	fmt.Printf("   Connection: mysql -h %s -P %d -u %s -p%s %s\n", host, availablePort, cfg.User, cfg.Password, cfg.Name)

	return nil
}

func (m *Manager) CreateSQLite(cfg *SQLiteConfig) error {
	fmt.Printf("Creating SQLite database: %s\n", cfg.FilePath)

	dir := filepath.Dir(cfg.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(cfg.FilePath)
	if err != nil {
		return fmt.Errorf("failed to create SQLite file: %w", err)
	}
	file.Close()

	if err := m.connTester.TestSQLite(cfg.FilePath); err != nil {
		return fmt.Errorf("failed to test SQLite connection: %w", err)
	}

	dbConfig := &config.DatabaseConfig{
		Name:     filepath.Base(cfg.FilePath),
		Type:     "sqlite",
		FilePath: cfg.FilePath,
		Created:  time.Now(),
	}

	if err := m.store.Save(dbConfig); err != nil {
		return fmt.Errorf("failed to save database config: %w", err)
	}

	fmt.Printf("✅ SQLite database created successfully!\n")
	fmt.Printf("   File: %s\n", cfg.FilePath)
	fmt.Printf("   Connection: sqlite3 %s\n", cfg.FilePath)

	return nil
}

func (m *Manager) ListDatabases(dbType string) error {
	databases, err := m.store.List(dbType)
	if err != nil {
		return fmt.Errorf("failed to load databases: %w", err)
	}

	if len(databases) == 0 {
		fmt.Println("No databases found.")
		return nil
	}

	fmt.Printf("SpinDB Databases:\n\n")
	for _, db := range databases {
		status := m.getStatus(&db)
		fmt.Printf("📊 %s (%s)\n", db.Name, db.Type)
		fmt.Printf("   Status: %s\n", status)
		if db.Port > 0 {
			fmt.Printf("   Port: %d\n", db.Port)
			if db.Type != "sqlite" {
				if db.Public {
					fmt.Printf("   Access: Public\n")
				} else {
					fmt.Printf("   Access: Private\n")
				}
			}
		}
		if db.FilePath != "" {
			fmt.Printf("   File: %s\n", db.FilePath)
		}
		fmt.Printf("   Created: %s\n", db.Created.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	return nil
}

func (m *Manager) getStatus(db *config.DatabaseConfig) string {
	if db.Type == "sqlite" {
		if _, err := os.Stat(db.FilePath); err != nil {
			return "❌ File not found"
		}
		return "✅ Available"
	}

	if db.ContainerID == "" {
		return "❓ Unknown"
	}

	if m.dockerService == nil {
		return "❓ Docker unavailable"
	}

	ctx := context.Background()
	running, err := m.dockerService.IsContainerRunning(ctx, db.ContainerID)
	if err != nil {
		return "❌ Error checking status"
	}

	if running {
		return "✅ Running"
	}
	return "⏸️ Stopped"
}

func (m *Manager) Connect(name string, testOnly bool) error {
	databases, err := m.store.List("")
	if err != nil {
		return fmt.Errorf("failed to load databases: %w", err)
	}

	var targetDB *config.DatabaseConfig
	for _, db := range databases {
		if db.Name == name {
			targetDB = &db
			break
		}
	}

	if targetDB == nil {
		return fmt.Errorf("database '%s' not found", name)
	}

	if testOnly {
		return m.testConnection(targetDB)
	}

	return m.openConnection(targetDB)
}

func (m *Manager) testConnection(db *config.DatabaseConfig) error {
	fmt.Printf("Testing connection to %s database '%s'...\n", db.Type, db.Name)

	switch db.Type {
	case "postgres":
		err := m.connTester.TestPostgres("localhost", db.Port, db.User, db.Password, db.Name)
		if err != nil {
			fmt.Printf("❌ Connection failed: %v\n", err)
			return err
		}
	case "mysql":
		err := m.connTester.TestMySQL("localhost", db.Port, db.User, db.Password, db.Name)
		if err != nil {
			fmt.Printf("❌ Connection failed: %v\n", err)
			return err
		}
	case "sqlite":
		err := m.connTester.TestSQLite(db.FilePath)
		if err != nil {
			fmt.Printf("❌ Connection failed: %v\n", err)
			return err
		}
	}

	fmt.Printf("✅ Connection successful!\n")
	return nil
}

func (m *Manager) openConnection(db *config.DatabaseConfig) error {
	db.LastUsed = time.Now()
	m.store.Save(db)

	var cmd *exec.Cmd
	var clientCmd string
	var installInstructions string

	switch db.Type {
	case "postgres":
		clientCmd = "psql"
		installInstructions = m.getInstallInstructions("postgresql-client", "psql")
		cmd = exec.Command("psql", "-h", "localhost", "-p", strconv.Itoa(db.Port), "-U", db.User, "-d", db.Name)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", db.Password))
	case "mysql":
		clientCmd = "mysql"
		installInstructions = m.getInstallInstructions("mysql-client", "mysql")
		cmd = exec.Command("mysql", "-h", "localhost", "-P", strconv.Itoa(db.Port), "-u", db.User, fmt.Sprintf("-p%s", db.Password), db.Name)
	case "sqlite":
		clientCmd = "sqlite3"
		installInstructions = m.getInstallInstructions("sqlite3", "sqlite3")
		cmd = exec.Command("sqlite3", db.FilePath)
	default:
		return fmt.Errorf("unsupported database type: %s", db.Type)
	}

	if _, err := exec.LookPath(clientCmd); err != nil {
		fmt.Printf("❌ Database client '%s' not found in PATH.\n\n", clientCmd)
		fmt.Printf("To connect interactively to your database, you need to install the client:\n\n")
		fmt.Printf("%s\n\n", installInstructions)
		fmt.Printf("Alternatively, you can:\n")
		fmt.Printf("• Use '--test-only' flag to test the connection without opening a shell\n")
		fmt.Printf("• Use Docker to run the client:\n")

		switch db.Type {
		case "postgres":
			fmt.Printf("  docker run -it --rm postgres:%s psql -h host.docker.internal -p %d -U %s -d %s\n",
				db.Version, db.Port, db.User, db.Name)
		case "mysql":
			fmt.Printf("  docker run -it --rm mysql:%s mysql -h host.docker.internal -P %d -u %s -p%s %s\n",
				db.Version, db.Port, db.User, db.Password, db.Name)
		case "sqlite":
			fmt.Printf("  docker run -it --rm -v %s:/db alpine sqlite3 /db\n", db.FilePath)
		}

		return fmt.Errorf("client '%s' not available", clientCmd)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Opening %s shell for database '%s'...\n", db.Type, db.Name)
	return cmd.Run()
}

func (m *Manager) getInstallInstructions(packageName, clientName string) string {
	var instructions []string

	instructions = append(instructions, fmt.Sprintf("Ubuntu/Debian: sudo apt update && sudo apt install %s", packageName))

	switch clientName {
	case "psql":
		instructions = append(instructions, "CentOS/RHEL: sudo yum install postgresql")
		instructions = append(instructions, "Fedora: sudo dnf install postgresql")
		instructions = append(instructions, "macOS: brew install postgresql")
		instructions = append(instructions, "Alpine: apk add postgresql-client")
	case "mysql":
		instructions = append(instructions, "CentOS/RHEL: sudo yum install mysql")
		instructions = append(instructions, "Fedora: sudo dnf install mysql")
		instructions = append(instructions, "macOS: brew install mysql-client")
		instructions = append(instructions, "Alpine: apk add mysql-client")
	case "sqlite3":
		instructions = append(instructions, "CentOS/RHEL: sudo yum install sqlite")
		instructions = append(instructions, "Fedora: sudo dnf install sqlite")
		instructions = append(instructions, "macOS: brew install sqlite")
		instructions = append(instructions, "Alpine: apk add sqlite")
	}

	result := ""
	for _, instruction := range instructions {
		result += "• " + instruction + "\n"
	}

	return strings.TrimSuffix(result, "\n")
}

func (m *Manager) ShowInfo(name string, showCredentials bool) error {
	databases, err := m.store.List("")
	if err != nil {
		return fmt.Errorf("failed to load databases: %w", err)
	}

	var targetDB *config.DatabaseConfig
	for _, db := range databases {
		if db.Name == name {
			targetDB = &db
			break
		}
	}

	if targetDB == nil {
		return fmt.Errorf("database '%s' not found", name)
	}

	fmt.Printf("Database Information: %s\n", targetDB.Name)
	fmt.Printf("─────────────────────────────────────\n")
	fmt.Printf("Type:         %s\n", targetDB.Type)
	fmt.Printf("Version:      %s\n", targetDB.Version)
	fmt.Printf("Status:       %s\n", m.getStatus(targetDB))
	fmt.Printf("Created:      %s\n", targetDB.Created.Format("2006-01-02 15:04:05"))

	if !targetDB.LastUsed.IsZero() {
		fmt.Printf("Last Used:    %s\n", targetDB.LastUsed.Format("2006-01-02 15:04:05"))
	}

	if targetDB.Port > 0 {
		fmt.Printf("Port:         %d\n", targetDB.Port)
		if targetDB.Type != "sqlite" {
			if targetDB.Public {
				fmt.Printf("Access:       Public (externally accessible)\n")
			} else {
				fmt.Printf("Access:       Private (localhost only)\n")
			}
		}
	}

	if targetDB.FilePath != "" {
		fmt.Printf("File Path:    %s\n", targetDB.FilePath)
	}

	if showCredentials && targetDB.Type != "sqlite" {
		fmt.Printf("User:         %s\n", targetDB.User)
		fmt.Printf("Password:     %s\n", targetDB.Password)
	}

	if targetDB.ContainerID != "" {
		fmt.Printf("Container ID: %s\n", targetDB.ContainerID)
	}

	return nil
}

func (m *Manager) Delete(name, file string, force bool) error {
	var targetDB *config.DatabaseConfig

	if name != "" {
		databases, err := m.store.List("")
		if err != nil {
			return fmt.Errorf("failed to load databases: %w", err)
		}

		for _, db := range databases {
			if db.Name == name {
				targetDB = &db
				break
			}
		}

		if targetDB == nil {
			return fmt.Errorf("database '%s' not found", name)
		}
	}

	if targetDB != nil {
		if !force {
			fmt.Printf("Are you sure you want to delete database '%s'? This action cannot be undone.\n", targetDB.Name)
			fmt.Print("Type 'yes' to confirm: ")
			var confirmation string
			fmt.Scanln(&confirmation)
			if strings.ToLower(confirmation) != "yes" {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		if targetDB.ContainerID != "" && m.dockerService != nil {
			ctx := context.Background()
			fmt.Printf("Stopping and removing container...\n")
			m.dockerService.StopContainer(ctx, targetDB.ContainerID)
			m.dockerService.RemoveContainer(ctx, targetDB.ContainerID, true)
		}

		if err := m.store.Delete(targetDB.Name, targetDB.Type); err != nil {
			return fmt.Errorf("failed to remove database config: %w", err)
		}

		fmt.Printf("✅ Database '%s' deleted successfully!\n", targetDB.Name)
	}

	return nil
}

func (m *Manager) Start(name string) error {
	databases, err := m.store.List("")
	if err != nil {
		return fmt.Errorf("failed to load databases: %w", err)
	}

	var targetDB *config.DatabaseConfig
	for _, db := range databases {
		if db.Name == name {
			targetDB = &db
			break
		}
	}

	if targetDB == nil {
		return fmt.Errorf("database '%s' not found", name)
	}

	if targetDB.Type == "sqlite" {
		return fmt.Errorf("SQLite databases don't need to be started")
	}

	if targetDB.ContainerID == "" {
		return fmt.Errorf("no container ID found for database '%s'", name)
	}

	if m.dockerService == nil {
		return fmt.Errorf("docker service not available")
	}

	ctx := context.Background()
	fmt.Printf("Starting database '%s'...\n", name)

	if err := m.dockerService.StartContainer(ctx, targetDB.ContainerID); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("✅ Database '%s' started successfully!\n", name)
	return nil
}

func (m *Manager) Stop(name string) error {
	databases, err := m.store.List("")
	if err != nil {
		return fmt.Errorf("failed to load databases: %w", err)
	}

	var targetDB *config.DatabaseConfig
	for _, db := range databases {
		if db.Name == name {
			targetDB = &db
			break
		}
	}

	if targetDB == nil {
		return fmt.Errorf("database '%s' not found", name)
	}

	if targetDB.Type == "sqlite" {
		return fmt.Errorf("SQLite databases don't need to be stopped")
	}

	if targetDB.ContainerID == "" {
		return fmt.Errorf("no container ID found for database '%s'", name)
	}

	if m.dockerService == nil {
		return fmt.Errorf("docker service not available")
	}

	ctx := context.Background()
	fmt.Printf("Stopping database '%s'...\n", name)

	if err := m.dockerService.StopContainer(ctx, targetDB.ContainerID); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	fmt.Printf("✅ Database '%s' stopped successfully!\n", name)
	return nil
}

func (m *Manager) Restart(name string) error {
	if err := m.Stop(name); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	return m.Start(name)
}
