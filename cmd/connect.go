package cmd

import (
	"github.com/awade12/spindb/internal/db"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to or test a database connection",
	Long:  `Connect to a managed database or test the connection`,
	RunE:  connectToDatabase,
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringP("name", "n", "", "Database name (required)")
	connectCmd.Flags().Bool("test-only", false, "Only test the connection, don't open interactive session")
	connectCmd.MarkFlagRequired("name")
}

func connectToDatabase(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	testOnly, _ := cmd.Flags().GetBool("test-only")

	manager := db.NewManager()
	return manager.Connect(name, testOnly)
}
