package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Template struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Type        string            `yaml:"type"`
	Version     string            `yaml:"version"`
	Config      map[string]string `yaml:"config"`
	CreatedAt   time.Time         `yaml:"created_at"`
	Tags        []string          `yaml:"tags,omitempty"`
}

type TemplateStore struct {
	templatesDir string
}

func NewTemplateStore() *TemplateStore {
	home, _ := os.UserHomeDir()
	templatesDir := filepath.Join(home, ".spindb", "templates")
	os.MkdirAll(templatesDir, 0755)

	return &TemplateStore{
		templatesDir: templatesDir,
	}
}

func (ts *TemplateStore) Save(template *Template) error {
	template.CreatedAt = time.Now()

	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	filename := fmt.Sprintf("%s.yaml", template.Name)
	filepath := filepath.Join(ts.templatesDir, filename)

	return ioutil.WriteFile(filepath, data, 0644)
}

func (ts *TemplateStore) Load(name string) (*Template, error) {
	filename := fmt.Sprintf("%s.yaml", name)
	filepath := filepath.Join(ts.templatesDir, filename)

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("template '%s' not found: %w", name, err)
	}

	var template Template
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &template, nil
}

func (ts *TemplateStore) List() ([]*Template, error) {
	files, err := ioutil.ReadDir(ts.templatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	var templates []*Template
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			name := file.Name()[:len(file.Name())-5]
			template, err := ts.Load(name)
			if err != nil {
				continue
			}
			templates = append(templates, template)
		}
	}

	return templates, nil
}

func (ts *TemplateStore) Delete(name string) error {
	filename := fmt.Sprintf("%s.yaml", name)
	filepath := filepath.Join(ts.templatesDir, filename)

	return os.Remove(filepath)
}

func (ts *TemplateStore) Exists(name string) bool {
	filename := fmt.Sprintf("%s.yaml", name)
	filepath := filepath.Join(ts.templatesDir, filename)

	_, err := os.Stat(filepath)
	return err == nil
}

func (ts *TemplateStore) Export(name, exportPath string) error {
	template, err := ts.Load(name)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	return ioutil.WriteFile(exportPath, data, 0644)
}

func (ts *TemplateStore) Import(importPath string) error {
	data, err := ioutil.ReadFile(importPath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	var template Template
	if err := yaml.Unmarshal(data, &template); err != nil {
		return fmt.Errorf("failed to parse template file: %w", err)
	}

	return ts.Save(&template)
}

func GetPredefinedTemplates() []*Template {
	return []*Template{
		{
			Name:        "postgres-dev",
			Description: "PostgreSQL development environment",
			Type:        "postgres",
			Version:     "15",
			Config: map[string]string{
				"user":     "dev_user",
				"password": "dev_password",
				"port":     "5432",
			},
			Tags: []string{"dev", "postgres"},
		},
		{
			Name:        "postgres-test",
			Description: "PostgreSQL testing environment",
			Type:        "postgres",
			Version:     "15",
			Config: map[string]string{
				"user":     "test_user",
				"password": "test_password",
				"port":     "5433",
			},
			Tags: []string{"test", "postgres"},
		},
		{
			Name:        "mysql-dev",
			Description: "MySQL development environment",
			Type:        "mysql",
			Version:     "8.0",
			Config: map[string]string{
				"user":     "dev_user",
				"password": "dev_password",
				"port":     "3306",
			},
			Tags: []string{"dev", "mysql"},
		},
		{
			Name:        "mysql-test",
			Description: "MySQL testing environment",
			Type:        "mysql",
			Version:     "8.0",
			Config: map[string]string{
				"user":     "test_user",
				"password": "test_password",
				"port":     "3307",
			},
			Tags: []string{"test", "mysql"},
		},
	}
}
