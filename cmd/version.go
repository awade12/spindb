package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the current version of SpinDB`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SpinDB %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Built: %s\n", BuildDate)
		fmt.Println("A powerful CLI tool to spin up and manage databases")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
