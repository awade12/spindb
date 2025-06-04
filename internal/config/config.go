package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Default DefaultConfig `yaml:"default"`
	Docker  DockerConfig  `yaml:"docker"`
	Storage StorageConfig `yaml:"storage"`
}

type DefaultConfig struct {
	Postgres PostgresDefaults `yaml:"postgres"`
	MySQL    MySQLDefaults    `yaml:"mysql"`
}

type PostgresDefaults struct {
	Version string `yaml:"version"`
	Port    int    `yaml:"port"`
	User    string `yaml:"user"`
}

type MySQLDefaults struct {
	Version string `yaml:"version"`
	Port    int    `yaml:"port"`
	User    string `yaml:"user"`
}

type DockerConfig struct {
	Host           string `yaml:"host"`
	CleanupTimeout string `yaml:"cleanup_timeout"`
}

type StorageConfig struct {
	DataDir   string `yaml:"data_dir"`
	BackupDir string `yaml:"backup_dir"`
}

func Load() *Config {
	home, _ := os.UserHomeDir()

	return &Config{
		Default: DefaultConfig{
			Postgres: PostgresDefaults{
				Version: "15",
				Port:    5432,
				User:    "postgres",
			},
			MySQL: MySQLDefaults{
				Version: "8.0",
				Port:    3306,
				User:    "root",
			},
		},
		Docker: DockerConfig{
			Host:           "unix:///var/run/docker.sock",
			CleanupTimeout: "30s",
		},
		Storage: StorageConfig{
			DataDir:   filepath.Join(home, ".spindb", "data"),
			BackupDir: filepath.Join(home, ".spindb", "backups"),
		},
	}
}
