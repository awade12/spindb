package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseStore struct {
	configDir string
}

type DatabaseRegistry struct {
	Databases []DatabaseConfig `yaml:"databases"`
}

func NewDatabaseStore() *DatabaseStore {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".spindb")

	return &DatabaseStore{
		configDir: configDir,
	}
}

func (ds *DatabaseStore) ensureConfigDir() error {
	return os.MkdirAll(ds.configDir, 0755)
}

func (ds *DatabaseStore) Save(dbConfig *DatabaseConfig) error {
	if err := ds.ensureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	registry, err := ds.Load()
	if err != nil {
		registry = &DatabaseRegistry{Databases: []DatabaseConfig{}}
	}

	found := false
	for i, existing := range registry.Databases {
		if existing.Name == dbConfig.Name && existing.Type == dbConfig.Type {
			registry.Databases[i] = *dbConfig
			found = true
			break
		}
	}

	if !found {
		registry.Databases = append(registry.Databases, *dbConfig)
	}

	return ds.saveRegistry(registry)
}

func (ds *DatabaseStore) Load() (*DatabaseRegistry, error) {
	registryPath := filepath.Join(ds.configDir, "databases.yaml")

	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		return &DatabaseRegistry{Databases: []DatabaseConfig{}}, nil
	}

	data, err := os.ReadFile(registryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry file: %w", err)
	}

	var registry DatabaseRegistry
	if err := yaml.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry file: %w", err)
	}

	return &registry, nil
}

func (ds *DatabaseStore) Get(name, dbType string) (*DatabaseConfig, error) {
	registry, err := ds.Load()
	if err != nil {
		return nil, err
	}

	for _, db := range registry.Databases {
		if db.Name == name && db.Type == dbType {
			return &db, nil
		}
	}

	return nil, fmt.Errorf("database %s of type %s not found", name, dbType)
}

func (ds *DatabaseStore) Delete(name, dbType string) error {
	registry, err := ds.Load()
	if err != nil {
		return err
	}

	var newDatabases []DatabaseConfig
	found := false

	for _, db := range registry.Databases {
		if db.Name == name && db.Type == dbType {
			found = true
			continue
		}
		newDatabases = append(newDatabases, db)
	}

	if !found {
		return fmt.Errorf("database %s of type %s not found", name, dbType)
	}

	registry.Databases = newDatabases
	return ds.saveRegistry(registry)
}

func (ds *DatabaseStore) List(dbType string) ([]DatabaseConfig, error) {
	registry, err := ds.Load()
	if err != nil {
		return nil, err
	}

	if dbType == "" {
		return registry.Databases, nil
	}

	var filtered []DatabaseConfig
	for _, db := range registry.Databases {
		if db.Type == dbType {
			filtered = append(filtered, db)
		}
	}

	return filtered, nil
}

func (ds *DatabaseStore) saveRegistry(registry *DatabaseRegistry) error {
	registryPath := filepath.Join(ds.configDir, "databases.yaml")

	data, err := yaml.Marshal(registry)
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(registryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry file: %w", err)
	}

	return nil
}
