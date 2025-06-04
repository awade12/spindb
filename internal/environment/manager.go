package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/awade12/spindb/internal/config"
	"github.com/awade12/spindb/internal/db"
	"gopkg.in/yaml.v3"
)

type Environment struct {
	Name        string                            `yaml:"name"`
	Description string                            `yaml:"description"`
	Active      bool                              `yaml:"active"`
	Databases   map[string]*config.DatabaseConfig `yaml:"databases"`
	CreatedAt   time.Time                         `yaml:"created_at"`
	UpdatedAt   time.Time                         `yaml:"updated_at"`
}

type EnvironmentManager struct {
	envDir     string
	currentEnv string
	store      *config.DatabaseStore
	dbManager  *db.Manager
}

type BulkOperation struct {
	Environment string
	Databases   []string
	Operation   string
	Success     []string
	Failed      map[string]string
}

func NewEnvironmentManager() *EnvironmentManager {
	home, _ := os.UserHomeDir()
	envDir := filepath.Join(home, ".spindb", "environments")
	os.MkdirAll(envDir, 0755)

	store := config.NewDatabaseStore()
	dbManager := db.NewManager()

	em := &EnvironmentManager{
		envDir:    envDir,
		store:     store,
		dbManager: dbManager,
	}

	em.currentEnv = em.getCurrentEnvironment()
	return em
}

func (em *EnvironmentManager) CreateEnvironment(name, description string) error {
	if em.EnvironmentExists(name) {
		return fmt.Errorf("environment '%s' already exists", name)
	}

	env := &Environment{
		Name:        name,
		Description: description,
		Active:      false,
		Databases:   make(map[string]*config.DatabaseConfig),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return em.saveEnvironment(env)
}

func (em *EnvironmentManager) DeleteEnvironment(name string, force bool) error {
	if name == "default" {
		return fmt.Errorf("cannot delete the default environment")
	}

	if !em.EnvironmentExists(name) {
		return fmt.Errorf("environment '%s' not found", name)
	}

	env, err := em.LoadEnvironment(name)
	if err != nil {
		return err
	}

	if len(env.Databases) > 0 && !force {
		return fmt.Errorf("environment '%s' contains databases, use --force to delete", name)
	}

	if env.Active {
		if err := em.SwitchEnvironment("default"); err != nil {
			return fmt.Errorf("failed to switch to default environment: %w", err)
		}
	}

	envPath := filepath.Join(em.envDir, name+".yaml")
	return os.Remove(envPath)
}

func (em *EnvironmentManager) SwitchEnvironment(name string) error {
	if !em.EnvironmentExists(name) {
		if name != "default" {
			return fmt.Errorf("environment '%s' not found", name)
		}
		if err := em.CreateEnvironment("default", "Default environment"); err != nil {
			return err
		}
	}

	if em.currentEnv != "" {
		currentEnv, err := em.LoadEnvironment(em.currentEnv)
		if err == nil {
			currentEnv.Active = false
			em.saveEnvironment(currentEnv)
		}
	}

	env, err := em.LoadEnvironment(name)
	if err != nil {
		return err
	}

	env.Active = true
	env.UpdatedAt = time.Now()

	if err := em.saveEnvironment(env); err != nil {
		return err
	}

	em.currentEnv = name
	return em.saveCurrentEnvironment(name)
}

func (em *EnvironmentManager) ListEnvironments() ([]*Environment, error) {
	files, err := os.ReadDir(em.envDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read environments directory: %w", err)
	}

	var environments []*Environment
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			name := strings.TrimSuffix(file.Name(), ".yaml")
			env, err := em.LoadEnvironment(name)
			if err != nil {
				continue
			}
			environments = append(environments, env)
		}
	}

	return environments, nil
}

func (em *EnvironmentManager) LoadEnvironment(name string) (*Environment, error) {
	envPath := filepath.Join(em.envDir, name+".yaml")

	data, err := os.ReadFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("environment '%s' not found: %w", name, err)
	}

	var env Environment
	if err := yaml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("failed to parse environment file: %w", err)
	}

	return &env, nil
}

func (em *EnvironmentManager) AddDatabaseToEnvironment(envName, dbName string) error {
	env, err := em.LoadEnvironment(envName)
	if err != nil {
		return err
	}

	registry, err := em.store.Load()
	if err != nil {
		return fmt.Errorf("failed to load database registry: %w", err)
	}

	var dbConfig *config.DatabaseConfig
	for _, db := range registry.Databases {
		if db.Name == dbName {
			dbConfig = &db
			break
		}
	}

	if dbConfig == nil {
		return fmt.Errorf("database '%s' not found", dbName)
	}

	env.Databases[dbName] = dbConfig
	env.UpdatedAt = time.Now()

	return em.saveEnvironment(env)
}

func (em *EnvironmentManager) RemoveDatabaseFromEnvironment(envName, dbName string) error {
	env, err := em.LoadEnvironment(envName)
	if err != nil {
		return err
	}

	if _, exists := env.Databases[dbName]; !exists {
		return fmt.Errorf("database '%s' not found in environment '%s'", dbName, envName)
	}

	delete(env.Databases, dbName)
	env.UpdatedAt = time.Now()

	return em.saveEnvironment(env)
}

func (em *EnvironmentManager) BulkStart(envName string, dbNames []string) *BulkOperation {
	op := &BulkOperation{
		Environment: envName,
		Databases:   dbNames,
		Operation:   "start",
		Success:     []string{},
		Failed:      make(map[string]string),
	}

	env, err := em.LoadEnvironment(envName)
	if err != nil {
		for _, dbName := range dbNames {
			op.Failed[dbName] = err.Error()
		}
		return op
	}

	for _, dbName := range dbNames {
		if _, exists := env.Databases[dbName]; !exists {
			op.Failed[dbName] = "database not found in environment"
			continue
		}

		if err := em.dbManager.Start(dbName); err != nil {
			op.Failed[dbName] = err.Error()
		} else {
			op.Success = append(op.Success, dbName)
		}
	}

	return op
}

func (em *EnvironmentManager) BulkStop(envName string, dbNames []string) *BulkOperation {
	op := &BulkOperation{
		Environment: envName,
		Databases:   dbNames,
		Operation:   "stop",
		Success:     []string{},
		Failed:      make(map[string]string),
	}

	env, err := em.LoadEnvironment(envName)
	if err != nil {
		for _, dbName := range dbNames {
			op.Failed[dbName] = err.Error()
		}
		return op
	}

	for _, dbName := range dbNames {
		if _, exists := env.Databases[dbName]; !exists {
			op.Failed[dbName] = "database not found in environment"
			continue
		}

		if err := em.dbManager.Stop(dbName); err != nil {
			op.Failed[dbName] = err.Error()
		} else {
			op.Success = append(op.Success, dbName)
		}
	}

	return op
}

func (em *EnvironmentManager) BulkRestart(envName string, dbNames []string) *BulkOperation {
	op := &BulkOperation{
		Environment: envName,
		Databases:   dbNames,
		Operation:   "restart",
		Success:     []string{},
		Failed:      make(map[string]string),
	}

	env, err := em.LoadEnvironment(envName)
	if err != nil {
		for _, dbName := range dbNames {
			op.Failed[dbName] = err.Error()
		}
		return op
	}

	for _, dbName := range dbNames {
		if _, exists := env.Databases[dbName]; !exists {
			op.Failed[dbName] = "database not found in environment"
			continue
		}

		if err := em.dbManager.Restart(dbName); err != nil {
			op.Failed[dbName] = err.Error()
		} else {
			op.Success = append(op.Success, dbName)
		}
	}

	return op
}

func (em *EnvironmentManager) GetCurrentEnvironment() string {
	return em.currentEnv
}

func (em *EnvironmentManager) EnvironmentExists(name string) bool {
	envPath := filepath.Join(em.envDir, name+".yaml")
	_, err := os.Stat(envPath)
	return err == nil
}

func (em *EnvironmentManager) saveEnvironment(env *Environment) error {
	envPath := filepath.Join(em.envDir, env.Name+".yaml")

	data, err := yaml.Marshal(env)
	if err != nil {
		return fmt.Errorf("failed to marshal environment: %w", err)
	}

	return os.WriteFile(envPath, data, 0644)
}

func (em *EnvironmentManager) getCurrentEnvironment() string {
	currentPath := filepath.Join(em.envDir, ".current")

	data, err := os.ReadFile(currentPath)
	if err != nil {
		return "default"
	}

	current := strings.TrimSpace(string(data))
	if current == "" {
		return "default"
	}

	return current
}

func (em *EnvironmentManager) saveCurrentEnvironment(name string) error {
	currentPath := filepath.Join(em.envDir, ".current")
	return os.WriteFile(currentPath, []byte(name), 0644)
}

func (em *EnvironmentManager) IsolateEnvironment(envName string) error {
	env, err := em.LoadEnvironment(envName)
	if err != nil {
		return err
	}

	for dbName := range env.Databases {
		fmt.Printf("Stopping database '%s' in environment '%s'\n", dbName, envName)
		if err := em.dbManager.Stop(dbName); err != nil {
			fmt.Printf("Warning: failed to stop database '%s': %v\n", dbName, err)
		}
	}

	fmt.Printf("Environment '%s' isolated successfully\n", envName)
	return nil
}

func (em *EnvironmentManager) ActivateEnvironment(envName string) error {
	env, err := em.LoadEnvironment(envName)
	if err != nil {
		return err
	}

	for dbName := range env.Databases {
		fmt.Printf("Starting database '%s' in environment '%s'\n", dbName, envName)
		if err := em.dbManager.Start(dbName); err != nil {
			fmt.Printf("Warning: failed to start database '%s': %v\n", dbName, err)
		}
	}

	fmt.Printf("Environment '%s' activated successfully\n", envName)
	return nil
}
