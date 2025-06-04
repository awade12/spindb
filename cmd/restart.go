package cmd

import (
	"fmt"

	"github.com/awade12/spindb/internal/db"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart [database-name]",
	Short: "Restart a database instance",
	Long:  `Restart a database container (PostgreSQL or MySQL)`,
	Args:  cobra.ExactArgs(1),
	RunE:  restartDatabase,
}

func init() {
	rootCmd.AddCommand(restartCmd)
}

func restartDatabase(cmd *cobra.Command, args []string) error {
	manager := db.NewManager()
	name := args[0]

	fmt.Printf("Restarting database '%s'...\n", name)
	return manager.Restart(name)
}
