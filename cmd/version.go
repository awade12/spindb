package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the current version of SpinDB`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SpinDB v0.1.0")
		fmt.Println("A powerful CLI tool to spin up and manage databases")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
