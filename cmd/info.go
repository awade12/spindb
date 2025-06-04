package cmd

import (
	"github.com/awade12/spindb/internal/db"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show database information",
	Long:  `Display detailed information about a managed database`,
	RunE:  showDatabaseInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringP("name", "n", "", "Database name (required)")
	infoCmd.Flags().Bool("show-credentials", false, "Show database credentials")
	infoCmd.MarkFlagRequired("name")
}

func showDatabaseInfo(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	showCreds, _ := cmd.Flags().GetBool("show-credentials")

	manager := db.NewManager()
	return manager.ShowInfo(name, showCreds)
}
