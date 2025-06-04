package config

import "time"

type DatabaseConfig struct {
	Name        string    `yaml:"name"`
	Type        string    `yaml:"type"`
	Version     string    `yaml:"version,omitempty"`
	Port        int       `yaml:"port,omitempty"`
	User        string    `yaml:"user,omitempty"`
	Password    string    `yaml:"password,omitempty"`
	FilePath    string    `yaml:"file_path,omitempty"`
	Public      bool      `yaml:"public,omitempty"`
	ContainerID string    `yaml:"container_id,omitempty"`
	Created     time.Time `yaml:"created"`
	LastUsed    time.Time `yaml:"last_used,omitempty"`
}
