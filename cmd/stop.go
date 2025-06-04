package cmd

import (
	"fmt"

	"github.com/awade12/spindb/internal/db"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [database-name]",
	Short: "Stop a database instance",
	Long:  `Stop a running database container (PostgreSQL or MySQL)`,
	Args:  cobra.ExactArgs(1),
	RunE:  stopDatabase,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stopDatabase(cmd *cobra.Command, args []string) error {
	manager := db.NewManager()
	name := args[0]

	fmt.Printf("Stopping database '%s'...\n", name)
	return manager.Stop(name)
}
