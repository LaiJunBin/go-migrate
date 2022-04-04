package cmd

import (
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/laijunbin/go-migrate/config"
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback the last database migration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rollback called")

		migrator := config.Migrator

		migrations, err := migrator.GetMigrations()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		sort.SliceStable(migrations, func(i, j int) bool {
			return migrations[j].Id < migrations[i].Id
		})

		if len(migrations) > 0 {
			batch := migrations[0].Batch
			for _, m := range migrations {
				if batch != m.Batch {
					break
				}
				for _, v := range config.Migrations {
					migration := reflect.TypeOf(v).String()
					if migration == m.Migration {
						batch = m.Batch
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

			if batch == 1 {
				migrator.DropTableIfExists()
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}
