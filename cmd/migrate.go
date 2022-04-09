package cmd

import (
	"fmt"
	"os"
	"reflect"

	"github.com/laijunbin/go-migrate/config"
	"github.com/laijunbin/go-migrate/pkg/interfaces"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Execute database migration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrate called")
		migrator := config.Migrator

		exists, err := migrator.CheckTable()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		if !exists {
			if err := migrator.CreateTable(); err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
		}

		migrations, err := migrator.GetMigrations()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		batch := 0

	nextMigrate:
		for _, v := range config.Migrations {
			migration := reflect.TypeOf(v).String()
			for _, m := range migrations {
				batch = m.Batch
				if migration == m.Migration {
					continue nextMigrate
				}
			}
			if err := v.Up(); err != nil {
				e, ok := err.(interfaces.Seeder)
				if ok {
					if e.Error() != "" {
						fmt.Println("Error: ", err)
						os.Exit(1)
					}
				} else {
					fmt.Println("Error: ", err)
					os.Exit(1)
				}
			}
			if err := migrator.WriteRecord(migration, batch+1); err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			fmt.Printf("migrate %s success.\n", migration)
		}

	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
