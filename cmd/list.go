package cmd

import (
	"github.com/awade12/spindb/internal/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed databases",
	Long:  `Display all databases currently managed by SpinDB`,
	RunE:  listDatabases,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("type", "t", "", "Filter by database type (postgres, mysql, sqlite)")
}

func listDatabases(cmd *cobra.Command, args []string) error {
	dbType, _ := cmd.Flags().GetString("type")

	manager := db.NewManager()
	return manager.ListDatabases(dbType)
}
