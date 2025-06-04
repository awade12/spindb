package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/awade12/spindb/internal/config"
	"github.com/awade12/spindb/internal/db"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage database templates",
	Long:  `Create, list, import, export and manage database templates for quick database setup`,
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates",
	Long:  `List all available database templates (predefined and custom)`,
	RunE:  listTemplates,
}

var templateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new template",
	Long:  `Create a new database template from existing configuration`,
	RunE:  createTemplate,
}

var templateDeleteCmd = &cobra.Command{
	Use:   "delete [template-name]",
	Short: "Delete a template",
	Long:  `Delete a custom database template`,
	Args:  cobra.ExactArgs(1),
	RunE:  deleteTemplate,
}

var templateShowCmd = &cobra.Command{
	Use:   "show [template-name]",
	Short: "Show template details",
	Long:  `Display detailed information about a template`,
	Args:  cobra.ExactArgs(1),
	RunE:  showTemplate,
}

var templateInstallCmd = &cobra.Command{
	Use:   "install [template-name] [database-name]",
	Short: "Create database from template",
	Long:  `Create a new database instance using a template`,
	Args:  cobra.ExactArgs(2),
	RunE:  installTemplate,
}

var templateImportCmd = &cobra.Command{
	Use:   "import [file-path]",
	Short: "Import template from file",
	Long:  `Import a database template from a YAML file`,
	Args:  cobra.ExactArgs(1),
	RunE:  importTemplate,
}

var templateExportCmd = &cobra.Command{
	Use:   "export [template-name] [file-path]",
	Short: "Export template to file",
	Long:  `Export a database template to a YAML file`,
	Args:  cobra.ExactArgs(2),
	RunE:  exportTemplate,
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateCreateCmd)
	templateCmd.AddCommand(templateDeleteCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateInstallCmd)
	templateCmd.AddCommand(templateImportCmd)
	templateCmd.AddCommand(templateExportCmd)

	templateCreateCmd.Flags().StringP("name", "n", "", "Template name (required)")
	templateCreateCmd.Flags().StringP("description", "d", "", "Template description")
	templateCreateCmd.Flags().StringP("type", "t", "", "Database type (postgres/mysql/sqlite) (required)")
	templateCreateCmd.Flags().StringP("version", "v", "", "Database version")
	templateCreateCmd.Flags().StringP("user", "u", "", "Database user")
	templateCreateCmd.Flags().StringP("password", "p", "", "Database password")
	templateCreateCmd.Flags().StringP("port", "", "", "Database port")
	templateCreateCmd.Flags().StringSliceP("tags", "", []string{}, "Template tags")
	templateCreateCmd.MarkFlagRequired("name")
	templateCreateCmd.MarkFlagRequired("type")

	templateInstallCmd.Flags().StringP("password", "p", "", "Override template password")
	templateInstallCmd.Flags().StringP("port", "", "", "Override template port")
	templateInstallCmd.Flags().Bool("public", false, "Make database publicly accessible")
}

func listTemplates(cmd *cobra.Command, args []string) error {
	store := config.NewTemplateStore()

	predefined := config.GetPredefinedTemplates()
	custom, err := store.List()
	if err != nil {
		return fmt.Errorf("failed to load custom templates: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tTYPE\tVERSION\tDESCRIPTION\tSOURCE\tTAGS")
	fmt.Fprintln(w, "----\t----\t-------\t-----------\t------\t----")

	allTemplates := append(predefined, custom...)
	sort.Slice(allTemplates, func(i, j int) bool {
		return allTemplates[i].Name < allTemplates[j].Name
	})

	for _, template := range allTemplates {
		source := "predefined"
		if template.CreatedAt.After(template.CreatedAt.Truncate(24 * 365 * 10)) {
			source = "custom"
		}

		tags := strings.Join(template.Tags, ", ")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			template.Name,
			template.Type,
			template.Version,
			template.Description,
			source,
			tags,
		)
	}

	return w.Flush()
}

func createTemplate(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	dbType, _ := cmd.Flags().GetString("type")
	version, _ := cmd.Flags().GetString("version")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	port, _ := cmd.Flags().GetString("port")
	tags, _ := cmd.Flags().GetStringSlice("tags")

	if dbType != "postgres" && dbType != "mysql" && dbType != "sqlite" {
		return fmt.Errorf("invalid database type: %s (must be postgres, mysql, or sqlite)", dbType)
	}

	store := config.NewTemplateStore()

	if store.Exists(name) {
		return fmt.Errorf("template '%s' already exists", name)
	}

	templateConfig := make(map[string]string)
	if user != "" {
		templateConfig["user"] = user
	}
	if password != "" {
		templateConfig["password"] = password
	}
	if port != "" {
		templateConfig["port"] = port
	}

	template := &config.Template{
		Name:        name,
		Description: description,
		Type:        dbType,
		Version:     version,
		Config:      templateConfig,
		Tags:        tags,
	}

	if err := store.Save(template); err != nil {
		return fmt.Errorf("failed to save template: %w", err)
	}

	fmt.Printf("✅ Template '%s' created successfully!\n", name)
	return nil
}

func deleteTemplate(cmd *cobra.Command, args []string) error {
	name := args[0]
	store := config.NewTemplateStore()

	if !store.Exists(name) {
		return fmt.Errorf("template '%s' not found", name)
	}

	if err := store.Delete(name); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	fmt.Printf("✅ Template '%s' deleted successfully!\n", name)
	return nil
}

func showTemplate(cmd *cobra.Command, args []string) error {
	name := args[0]
	store := config.NewTemplateStore()

	template, err := store.Load(name)
	if err != nil {
		predefined := config.GetPredefinedTemplates()
		for _, t := range predefined {
			if t.Name == name {
				template = t
				break
			}
		}
		if template == nil {
			return fmt.Errorf("template '%s' not found", name)
		}
	}

	fmt.Printf("Template: %s\n", template.Name)
	fmt.Printf("Type: %s\n", template.Type)
	fmt.Printf("Version: %s\n", template.Version)
	fmt.Printf("Description: %s\n", template.Description)

	if len(template.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(template.Tags, ", "))
	}

	if !template.CreatedAt.IsZero() {
		fmt.Printf("Created: %s\n", template.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	if len(template.Config) > 0 {
		fmt.Printf("\nConfiguration:\n")
		for key, value := range template.Config {
			if key == "password" {
				fmt.Printf("  %s: ***\n", key)
			} else {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
	}

	return nil
}

func installTemplate(cmd *cobra.Command, args []string) error {
	templateName := args[0]
	databaseName := args[1]

	passwordOverride, _ := cmd.Flags().GetString("password")
	portOverride, _ := cmd.Flags().GetString("port")
	public, _ := cmd.Flags().GetBool("public")

	store := config.NewTemplateStore()

	var template *config.Template
	var err error

	template, err = store.Load(templateName)
	if err != nil {
		predefined := config.GetPredefinedTemplates()
		for _, t := range predefined {
			if t.Name == templateName {
				template = t
				break
			}
		}
		if template == nil {
			return fmt.Errorf("template '%s' not found", templateName)
		}
	}

	manager := db.NewManager()

	switch template.Type {
	case "postgres":
		port := 5432
		if portOverride != "" {
			if p, err := strconv.Atoi(portOverride); err == nil {
				port = p
			}
		} else if template.Config["port"] != "" {
			if p, err := strconv.Atoi(template.Config["port"]); err == nil {
				port = p
			}
		}

		password := template.Config["password"]
		if passwordOverride != "" {
			password = passwordOverride
		}

		config := &db.PostgresConfig{
			Name:     databaseName,
			User:     template.Config["user"],
			Password: password,
			Port:     port,
			Version:  template.Version,
			Public:   public,
		}

		fmt.Printf("Creating PostgreSQL database '%s' from template '%s'...\n", databaseName, templateName)
		return manager.CreatePostgres(config)

	case "mysql":
		port := 3306
		if portOverride != "" {
			if p, err := strconv.Atoi(portOverride); err == nil {
				port = p
			}
		} else if template.Config["port"] != "" {
			if p, err := strconv.Atoi(template.Config["port"]); err == nil {
				port = p
			}
		}

		password := template.Config["password"]
		if passwordOverride != "" {
			password = passwordOverride
		}

		config := &db.MySQLConfig{
			Name:     databaseName,
			User:     template.Config["user"],
			Password: password,
			Port:     port,
			Version:  template.Version,
		}

		fmt.Printf("Creating MySQL database '%s' from template '%s'...\n", databaseName, templateName)
		return manager.CreateMySQL(config)

	case "sqlite":
		config := &db.SQLiteConfig{
			FilePath: databaseName + ".db",
		}

		fmt.Printf("Creating SQLite database '%s' from template '%s'...\n", databaseName, templateName)
		return manager.CreateSQLite(config)

	default:
		return fmt.Errorf("unsupported database type: %s", template.Type)
	}
}

func importTemplate(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	store := config.NewTemplateStore()

	if err := store.Import(filePath); err != nil {
		return fmt.Errorf("failed to import template: %w", err)
	}

	fmt.Printf("✅ Template imported successfully from %s!\n", filePath)
	return nil
}

func exportTemplate(cmd *cobra.Command, args []string) error {
	templateName := args[0]
	filePath := args[1]

	store := config.NewTemplateStore()

	if err := store.Export(templateName, filePath); err != nil {
		return fmt.Errorf("failed to export template: %w", err)
	}

	fmt.Printf("✅ Template '%s' exported to %s!\n", templateName, filePath)
	return nil
}
