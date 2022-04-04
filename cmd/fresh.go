package cmd

import (
	"fmt"

	"github.com/laijunbin/go-migrate/config"
	"github.com/spf13/cobra"
)

var freshCmd = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fresh called")
		migrator := config.Migrator
		migrator.DropAllTable()
		migrateCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(freshCmd)
}
