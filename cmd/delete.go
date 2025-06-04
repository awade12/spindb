package cmd

import (
	"github.com/awade12/spindb/internal/db"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a database instance",
	Long:  `Delete a managed database instance and clean up resources`,
	RunE:  deleteDatabase,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("name", "n", "", "Database name")
	deleteCmd.Flags().StringP("file", "f", "", "SQLite database file path")
	deleteCmd.Flags().Bool("force", false, "Force deletion without confirmation")
}

func deleteDatabase(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	file, _ := cmd.Flags().GetString("file")
	force, _ := cmd.Flags().GetBool("force")

	manager := db.NewManager()
	return manager.Delete(name, file, force)
}
