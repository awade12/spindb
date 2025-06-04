package cmd

import (
	"fmt"

	"github.com/awade12/spindb/internal/db"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [database-name]",
	Short: "Start a database instance",
	Long:  `Start a stopped database container (PostgreSQL or MySQL)`,
	Args:  cobra.ExactArgs(1),
	RunE:  startDatabase,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startDatabase(cmd *cobra.Command, args []string) error {
	manager := db.NewManager()
	name := args[0]

	fmt.Printf("Starting database '%s'...\n", name)
	return manager.Start(name)
}
