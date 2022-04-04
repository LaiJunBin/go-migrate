package cmd

import (
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/laijunbin/go-migrate/config"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Rollback all database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reset called")
		migrator := config.Migrator
		migrations, err := migrator.GetMigrations()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		sort.SliceStable(migrations, func(i, j int) bool {
			return migrations[j].Id < migrations[i].Id
		})

		for _, m := range migrations {
			for _, v := range config.Migrations {
				migration := reflect.TypeOf(v).String()
				if migration == m.Migration {
					if err := v.Down(); err != nil {
						fmt.Println("Error: ", err)
						os.Exit(1)
					}
					if err := migrator.DeleteRecord(m.Id); err != nil {
						fmt.Println("Error: ", err)
						os.Exit(1)
					}
					fmt.Printf("rollback %s success.\n", migration)
				}
			}
		}

		migrator.DropTableIfExists()
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
