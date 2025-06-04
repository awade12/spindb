package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/awade12/spindb/internal/environment"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage database environments",
	Long:  `Create, switch, and manage database environments for isolation and organization`,
}

var envCreateCmd = &cobra.Command{
	Use:   "create [environment-name]",
	Short: "Create a new environment",
	Long:  `Create a new environment for organizing databases`,
	Args:  cobra.ExactArgs(1),
	RunE:  createEnvironment,
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environments",
	Long:  `List all available environments with their status`,
	RunE:  listEnvironments,
}

var envSwitchCmd = &cobra.Command{
	Use:   "switch [environment-name]",
	Short: "Switch to an environment",
	Long:  `Switch to the specified environment`,
	Args:  cobra.ExactArgs(1),
	RunE:  switchEnvironment,
}

var envDeleteCmd = &cobra.Command{
	Use:   "delete [environment-name]",
	Short: "Delete an environment",
	Long:  `Delete the specified environment`,
	Args:  cobra.ExactArgs(1),
	RunE:  deleteEnvironment,
}

var envShowCmd = &cobra.Command{
	Use:   "show [environment-name]",
	Short: "Show environment details",
	Long:  `Display detailed information about an environment`,
	Args:  cobra.ExactArgs(1),
	RunE:  showEnvironment,
}

var envAddCmd = &cobra.Command{
	Use:   "add [environment-name] [database-name]",
	Short: "Add database to environment",
	Long:  `Add a database to the specified environment`,
	Args:  cobra.ExactArgs(2),
	RunE:  addDatabaseToEnvironment,
}

var envRemoveCmd = &cobra.Command{
	Use:   "remove [environment-name] [database-name]",
	Short: "Remove database from environment",
	Long:  `Remove a database from the specified environment`,
	Args:  cobra.ExactArgs(2),
	RunE:  removeDatabaseFromEnvironment,
}

var envBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Bulk operations on environment databases",
	Long:  `Perform bulk operations on all databases in an environment`,
}

var envBulkStartCmd = &cobra.Command{
	Use:   "start [environment-name]",
	Short: "Start all databases in environment",
	Long:  `Start all databases in the specified environment`,
	Args:  cobra.ExactArgs(1),
	RunE:  bulkStartEnvironment,
}

var envBulkStopCmd = &cobra.Command{
	Use:   "stop [environment-name]",
	Short: "Stop all databases in environment",
	Long:  `Stop all databases in the specified environment`,
	Args:  cobra.ExactArgs(1),
	RunE:  bulkStopEnvironment,
}

var envBulkRestartCmd = &cobra.Command{
	Use:   "restart [environment-name]",
	Short: "Restart all databases in environment",
	Long:  `Restart all databases in the specified environment`,
	Args:  cobra.ExactArgs(1),
	RunE:  bulkRestartEnvironment,
}

var envIsolateCmd = &cobra.Command{
	Use:   "isolate [environment-name]",
	Short: "Isolate environment",
	Long:  `Stop all databases in the environment to isolate it`,
	Args:  cobra.ExactArgs(1),
	RunE:  isolateEnvironment,
}

var envActivateCmd = &cobra.Command{
	Use:   "activate [environment-name]",
	Short: "Activate environment",
	Long:  `Start all databases in the environment to activate it`,
	Args:  cobra.ExactArgs(1),
	RunE:  activateEnvironment,
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envCreateCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envSwitchCmd)
	envCmd.AddCommand(envDeleteCmd)
	envCmd.AddCommand(envShowCmd)
	envCmd.AddCommand(envAddCmd)
	envCmd.AddCommand(envRemoveCmd)
	envCmd.AddCommand(envBulkCmd)
	envCmd.AddCommand(envIsolateCmd)
	envCmd.AddCommand(envActivateCmd)

	envBulkCmd.AddCommand(envBulkStartCmd)
	envBulkCmd.AddCommand(envBulkStopCmd)
	envBulkCmd.AddCommand(envBulkRestartCmd)

	envCreateCmd.Flags().StringP("description", "d", "", "Environment description")
	envDeleteCmd.Flags().Bool("force", false, "Force delete even if environment contains databases")
}

func createEnvironment(cmd *cobra.Command, args []string) error {
	name := args[0]
	description, _ := cmd.Flags().GetString("description")

	manager := environment.NewEnvironmentManager()

	if err := manager.CreateEnvironment(name, description); err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}

	fmt.Printf("✅ Environment '%s' created successfully!\n", name)
	return nil
}

func listEnvironments(cmd *cobra.Command, args []string) error {
	manager := environment.NewEnvironmentManager()

	environments, err := manager.ListEnvironments()
	if err != nil {
		return fmt.Errorf("failed to list environments: %w", err)
	}

	if len(environments) == 0 {
		fmt.Println("No environments found. Create one with 'spindb env create <name>'")
		return nil
	}

	current := manager.GetCurrentEnvironment()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tDESCRIPTION\tDATABASES\tACTIVE\tCREATED")
	fmt.Fprintln(w, "----\t-----------\t---------\t------\t-------")

	for _, env := range environments {
		active := ""
		if env.Name == current {
			active = "* CURRENT"
		} else if env.Active {
			active = "ACTIVE"
		}

		created := env.CreatedAt.Format("2006-01-02")
		dbCount := fmt.Sprintf("%d", len(env.Databases))

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			env.Name,
			env.Description,
			dbCount,
			active,
			created,
		)
	}

	return w.Flush()
}

func switchEnvironment(cmd *cobra.Command, args []string) error {
	name := args[0]
	manager := environment.NewEnvironmentManager()

	if err := manager.SwitchEnvironment(name); err != nil {
		return fmt.Errorf("failed to switch environment: %w", err)
	}

	fmt.Printf("✅ Switched to environment '%s'\n", name)
	return nil
}

func deleteEnvironment(cmd *cobra.Command, args []string) error {
	name := args[0]
	force, _ := cmd.Flags().GetBool("force")

	manager := environment.NewEnvironmentManager()

	if err := manager.DeleteEnvironment(name, force); err != nil {
		return fmt.Errorf("failed to delete environment: %w", err)
	}

	fmt.Printf("✅ Environment '%s' deleted successfully!\n", name)
	return nil
}

func showEnvironment(cmd *cobra.Command, args []string) error {
	name := args[0]
	manager := environment.NewEnvironmentManager()

	env, err := manager.LoadEnvironment(name)
	if err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}

	fmt.Printf("Environment: %s\n", env.Name)
	fmt.Printf("Description: %s\n", env.Description)
	fmt.Printf("Active: %t\n", env.Active)
	fmt.Printf("Created: %s\n", env.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", env.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Databases: %d\n", len(env.Databases))

	if len(env.Databases) > 0 {
		fmt.Printf("\nDatabases in this environment:\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tVERSION\tPORT")
		fmt.Fprintln(w, "----\t----\t-------\t----")

		for _, db := range env.Databases {
			port := fmt.Sprintf("%d", db.Port)
			if db.Port == 0 {
				port = "-"
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				db.Name,
				db.Type,
				db.Version,
				port,
			)
		}
		w.Flush()
	}

	return nil
}

func addDatabaseToEnvironment(cmd *cobra.Command, args []string) error {
	envName := args[0]
	dbName := args[1]

	manager := environment.NewEnvironmentManager()

	if err := manager.AddDatabaseToEnvironment(envName, dbName); err != nil {
		return fmt.Errorf("failed to add database to environment: %w", err)
	}

	fmt.Printf("✅ Database '%s' added to environment '%s'\n", dbName, envName)
	return nil
}

func removeDatabaseFromEnvironment(cmd *cobra.Command, args []string) error {
	envName := args[0]
	dbName := args[1]

	manager := environment.NewEnvironmentManager()

	if err := manager.RemoveDatabaseFromEnvironment(envName, dbName); err != nil {
		return fmt.Errorf("failed to remove database from environment: %w", err)
	}

	fmt.Printf("✅ Database '%s' removed from environment '%s'\n", dbName, envName)
	return nil
}

func bulkStartEnvironment(cmd *cobra.Command, args []string) error {
	envName := args[0]
	manager := environment.NewEnvironmentManager()

	env, err := manager.LoadEnvironment(envName)
	if err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}

	var dbNames []string
	for dbName := range env.Databases {
		dbNames = append(dbNames, dbName)
	}

	fmt.Printf("Starting %d databases in environment '%s'...\n", len(dbNames), envName)

	result := manager.BulkStart(envName, dbNames)

	if len(result.Success) > 0 {
		fmt.Printf("✅ Successfully started: %s\n", strings.Join(result.Success, ", "))
	}

	if len(result.Failed) > 0 {
		fmt.Printf("❌ Failed to start:\n")
		for db, err := range result.Failed {
			fmt.Printf("   %s: %s\n", db, err)
		}
	}

	return nil
}

func bulkStopEnvironment(cmd *cobra.Command, args []string) error {
	envName := args[0]
	manager := environment.NewEnvironmentManager()

	env, err := manager.LoadEnvironment(envName)
	if err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}

	var dbNames []string
	for dbName := range env.Databases {
		dbNames = append(dbNames, dbName)
	}

	fmt.Printf("Stopping %d databases in environment '%s'...\n", len(dbNames), envName)

	result := manager.BulkStop(envName, dbNames)

	if len(result.Success) > 0 {
		fmt.Printf("✅ Successfully stopped: %s\n", strings.Join(result.Success, ", "))
	}

	if len(result.Failed) > 0 {
		fmt.Printf("❌ Failed to stop:\n")
		for db, err := range result.Failed {
			fmt.Printf("   %s: %s\n", db, err)
		}
	}

	return nil
}

func bulkRestartEnvironment(cmd *cobra.Command, args []string) error {
	envName := args[0]
	manager := environment.NewEnvironmentManager()

	env, err := manager.LoadEnvironment(envName)
	if err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}

	var dbNames []string
	for dbName := range env.Databases {
		dbNames = append(dbNames, dbName)
	}

	fmt.Printf("Restarting %d databases in environment '%s'...\n", len(dbNames), envName)

	result := manager.BulkRestart(envName, dbNames)

	if len(result.Success) > 0 {
		fmt.Printf("✅ Successfully restarted: %s\n", strings.Join(result.Success, ", "))
	}

	if len(result.Failed) > 0 {
		fmt.Printf("❌ Failed to restart:\n")
		for db, err := range result.Failed {
			fmt.Printf("   %s: %s\n", db, err)
		}
	}

	return nil
}

func isolateEnvironment(cmd *cobra.Command, args []string) error {
	envName := args[0]
	manager := environment.NewEnvironmentManager()

	if err := manager.IsolateEnvironment(envName); err != nil {
		return fmt.Errorf("failed to isolate environment: %w", err)
	}

	return nil
}

func activateEnvironment(cmd *cobra.Command, args []string) error {
	envName := args[0]
	manager := environment.NewEnvironmentManager()

	if err := manager.ActivateEnvironment(envName); err != nil {
		return fmt.Errorf("failed to activate environment: %w", err)
	}

	return nil
}
