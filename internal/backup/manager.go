package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/awade12/spindb/internal/config"
)

type BackupManager struct {
	backupDir string
	store     *config.DatabaseStore
}

type BackupInfo struct {
	Name       string    `yaml:"name"`
	Database   string    `yaml:"database"`
	Type       string    `yaml:"type"`
	Size       int64     `yaml:"size"`
	CreatedAt  time.Time `yaml:"created_at"`
	FilePath   string    `yaml:"file_path"`
	Compressed bool      `yaml:"compressed"`
}

func NewBackupManager() *BackupManager {
	cfg := config.Load()
	store := config.NewDatabaseStore()

	os.MkdirAll(cfg.Storage.BackupDir, 0755)

	return &BackupManager{
		backupDir: cfg.Storage.BackupDir,
		store:     store,
	}
}

func (bm *BackupManager) CreateBackup(dbName string, options *BackupOptions) (*BackupInfo, error) {
	registry, err := bm.store.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load database registry: %w", err)
	}

	var db *config.DatabaseConfig
	for _, database := range registry.Databases {
		if database.Name == dbName {
			db = &database
			break
		}
	}

	if db == nil {
		return nil, fmt.Errorf("database '%s' not found", dbName)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s_%s", dbName, timestamp)

	var backupPath string
	var backupErr error

	switch db.Type {
	case "postgres":
		backupPath, backupErr = bm.backupPostgres(db, backupName, options)
	case "mysql":
		backupPath, backupErr = bm.backupMySQL(db, backupName, options)
	case "sqlite":
		backupPath, backupErr = bm.backupSQLite(db, backupName, options)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", db.Type)
	}

	if backupErr != nil {
		return nil, fmt.Errorf("failed to create backup: %w", backupErr)
	}

	fileInfo, err := os.Stat(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup file info: %w", err)
	}

	backupInfo := &BackupInfo{
		Name:       backupName,
		Database:   dbName,
		Type:       db.Type,
		Size:       fileInfo.Size(),
		CreatedAt:  time.Now(),
		FilePath:   backupPath,
		Compressed: options != nil && options.Compress,
	}

	return backupInfo, nil
}

func (bm *BackupManager) RestoreBackup(backupName, targetDbName string) error {
	backupPath := filepath.Join(bm.backupDir, backupName)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file '%s' not found", backupName)
	}

	dbType := bm.detectBackupType(backupName)

	switch dbType {
	case "postgres":
		return bm.restorePostgres(backupPath, targetDbName)
	case "mysql":
		return bm.restoreMySQL(backupPath, targetDbName)
	case "sqlite":
		return bm.restoreSQLite(backupPath, targetDbName)
	default:
		return fmt.Errorf("unable to determine backup type for: %s", backupName)
	}
}

func (bm *BackupManager) ListBackups() ([]*BackupInfo, error) {
	files, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []*BackupInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		dbName, timestamp := bm.parseBackupName(file.Name())
		if dbName == "" {
			continue
		}

		parsedTime, _ := time.Parse("20060102_150405", timestamp)

		backup := &BackupInfo{
			Name:      strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
			Database:  dbName,
			Type:      bm.detectBackupType(file.Name()),
			Size:      info.Size(),
			CreatedAt: parsedTime,
			FilePath:  filepath.Join(bm.backupDir, file.Name()),
		}

		backups = append(backups, backup)
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	return backups, nil
}

func (bm *BackupManager) DeleteBackup(backupName string) error {
	backupPath := filepath.Join(bm.backupDir, backupName)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup '%s' not found", backupName)
	}

	return os.Remove(backupPath)
}

func (bm *BackupManager) backupPostgres(db *config.DatabaseConfig, backupName string, options *BackupOptions) (string, error) {
	backupPath := filepath.Join(bm.backupDir, backupName+".sql")
	if options != nil && options.Compress {
		backupPath += ".gz"
	}

	cmd := []string{
		"pg_dump",
		"-h", "localhost",
		"-p", fmt.Sprintf("%d", db.Port),
		"-U", db.User,
		"-d", db.Name,
		"--no-password",
	}

	if options != nil && options.SchemaOnly {
		cmd = append(cmd, "--schema-only")
	}

	if options != nil && options.DataOnly {
		cmd = append(cmd, "--data-only")
	}

	pgCmd := exec.Command(cmd[0], cmd[1:]...)
	pgCmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", db.Password))

	if options != nil && options.Compress {
		gzipCmd := exec.Command("gzip")

		pgCmd.Stdout, _ = gzipCmd.StdinPipe()
		gzipOutput, err := os.Create(backupPath)
		if err != nil {
			return "", err
		}
		defer gzipOutput.Close()

		gzipCmd.Stdout = gzipOutput

		if err := gzipCmd.Start(); err != nil {
			return "", err
		}

		if err := pgCmd.Run(); err != nil {
			return "", err
		}

		pgCmd.Stdout.(*os.File).Close()

		if err := gzipCmd.Wait(); err != nil {
			return "", err
		}
	} else {
		output, err := os.Create(backupPath)
		if err != nil {
			return "", err
		}
		defer output.Close()

		pgCmd.Stdout = output

		if err := pgCmd.Run(); err != nil {
			return "", err
		}
	}

	return backupPath, nil
}

func (bm *BackupManager) backupMySQL(db *config.DatabaseConfig, backupName string, options *BackupOptions) (string, error) {
	backupPath := filepath.Join(bm.backupDir, backupName+".sql")
	if options != nil && options.Compress {
		backupPath += ".gz"
	}

	cmd := []string{
		"mysqldump",
		"-h", "localhost",
		"-P", fmt.Sprintf("%d", db.Port),
		"-u", db.User,
		fmt.Sprintf("-p%s", db.Password),
		db.Name,
	}

	if options != nil && options.SchemaOnly {
		cmd = append(cmd, "--no-data")
	}

	if options != nil && options.DataOnly {
		cmd = append(cmd, "--no-create-info")
	}

	mysqlCmd := exec.Command(cmd[0], cmd[1:]...)

	if options != nil && options.Compress {
		gzipCmd := exec.Command("gzip")

		mysqlCmd.Stdout, _ = gzipCmd.StdinPipe()
		gzipOutput, err := os.Create(backupPath)
		if err != nil {
			return "", err
		}
		defer gzipOutput.Close()

		gzipCmd.Stdout = gzipOutput

		if err := gzipCmd.Start(); err != nil {
			return "", err
		}

		if err := mysqlCmd.Run(); err != nil {
			return "", err
		}

		mysqlCmd.Stdout.(*os.File).Close()

		if err := gzipCmd.Wait(); err != nil {
			return "", err
		}
	} else {
		output, err := os.Create(backupPath)
		if err != nil {
			return "", err
		}
		defer output.Close()

		mysqlCmd.Stdout = output

		if err := mysqlCmd.Run(); err != nil {
			return "", err
		}
	}

	return backupPath, nil
}

func (bm *BackupManager) backupSQLite(db *config.DatabaseConfig, backupName string, options *BackupOptions) (string, error) {
	backupPath := filepath.Join(bm.backupDir, backupName+".db")
	if options != nil && options.Compress {
		backupPath += ".gz"
	}

	if options != nil && options.Compress {
		cmd := exec.Command("gzip", "-c", db.FilePath)
		output, err := os.Create(backupPath)
		if err != nil {
			return "", err
		}
		defer output.Close()

		cmd.Stdout = output
		if err := cmd.Run(); err != nil {
			return "", err
		}
	} else {
		srcFile, err := os.Open(db.FilePath)
		if err != nil {
			return "", err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(backupPath)
		if err != nil {
			return "", err
		}
		defer dstFile.Close()

		if _, err := srcFile.WriteTo(dstFile); err != nil {
			return "", err
		}
	}

	return backupPath, nil
}

func (bm *BackupManager) restorePostgres(backupPath, targetDbName string) error {
	registry, err := bm.store.Load()
	if err != nil {
		return fmt.Errorf("failed to load database registry: %w", err)
	}

	var db *config.DatabaseConfig
	for _, database := range registry.Databases {
		if database.Name == targetDbName {
			db = &database
			break
		}
	}

	if db == nil {
		return fmt.Errorf("target database '%s' not found", targetDbName)
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(backupPath, ".gz") {
		zcat := exec.Command("zcat", backupPath)
		psql := exec.Command("psql",
			"-h", "localhost",
			"-p", fmt.Sprintf("%d", db.Port),
			"-U", db.User,
			"-d", db.Name,
		)
		psql.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", db.Password))

		psql.Stdin, _ = zcat.StdoutPipe()

		if err := zcat.Start(); err != nil {
			return err
		}

		if err := psql.Run(); err != nil {
			return err
		}

		return zcat.Wait()
	} else {
		cmd = exec.Command("psql",
			"-h", "localhost",
			"-p", fmt.Sprintf("%d", db.Port),
			"-U", db.User,
			"-d", db.Name,
			"-f", backupPath,
		)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", db.Password))
	}

	return cmd.Run()
}

func (bm *BackupManager) restoreMySQL(backupPath, targetDbName string) error {
	registry, err := bm.store.Load()
	if err != nil {
		return fmt.Errorf("failed to load database registry: %w", err)
	}

	var db *config.DatabaseConfig
	for _, database := range registry.Databases {
		if database.Name == targetDbName {
			db = &database
			break
		}
	}

	if db == nil {
		return fmt.Errorf("target database '%s' not found", targetDbName)
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(backupPath, ".gz") {
		zcat := exec.Command("zcat", backupPath)
		mysql := exec.Command("mysql",
			"-h", "localhost",
			"-P", fmt.Sprintf("%d", db.Port),
			"-u", db.User,
			fmt.Sprintf("-p%s", db.Password),
			db.Name,
		)

		mysql.Stdin, _ = zcat.StdoutPipe()

		if err := zcat.Start(); err != nil {
			return err
		}

		if err := mysql.Run(); err != nil {
			return err
		}

		return zcat.Wait()
	} else {
		cmd = exec.Command("mysql",
			"-h", "localhost",
			"-P", fmt.Sprintf("%d", db.Port),
			"-u", db.User,
			fmt.Sprintf("-p%s", db.Password),
			db.Name,
		)

		input, err := os.Open(backupPath)
		if err != nil {
			return err
		}
		defer input.Close()

		cmd.Stdin = input
	}

	return cmd.Run()
}

func (bm *BackupManager) restoreSQLite(backupPath, targetDbName string) error {
	registry, err := bm.store.Load()
	if err != nil {
		return fmt.Errorf("failed to load database registry: %w", err)
	}

	var db *config.DatabaseConfig
	for _, database := range registry.Databases {
		if database.Name == targetDbName {
			db = &database
			break
		}
	}

	if db == nil {
		return fmt.Errorf("target database '%s' not found", targetDbName)
	}

	if strings.HasSuffix(backupPath, ".gz") {
		cmd := exec.Command("gunzip", "-c", backupPath)
		output, err := os.Create(db.FilePath)
		if err != nil {
			return err
		}
		defer output.Close()

		cmd.Stdout = output
		return cmd.Run()
	} else {
		srcFile, err := os.Open(backupPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(db.FilePath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = srcFile.WriteTo(dstFile)
		return err
	}
}

func (bm *BackupManager) detectBackupType(filename string) string {
	if strings.Contains(filename, "_postgres_") || strings.Contains(filename, "_pg_") {
		return "postgres"
	}
	if strings.Contains(filename, "_mysql_") {
		return "mysql"
	}
	if strings.HasSuffix(filename, ".db") || strings.HasSuffix(filename, ".db.gz") {
		return "sqlite"
	}
	return "unknown"
}

func (bm *BackupManager) parseBackupName(filename string) (dbName, timestamp string) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	if strings.HasSuffix(base, ".gz") {
		base = strings.TrimSuffix(base, ".gz")
	}

	parts := strings.Split(base, "_")
	if len(parts) >= 3 {
		timestamp = parts[len(parts)-2] + "_" + parts[len(parts)-1]
		dbName = strings.Join(parts[:len(parts)-2], "_")
	}

	return
}

type BackupOptions struct {
	Compress   bool
	SchemaOnly bool
	DataOnly   bool
}
