package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Reset and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("refresh called")
		resetCmd.Run(cmd, args)
		migrateCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
