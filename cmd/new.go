package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/laijunbin/go-migrate/config"
	"github.com/laijunbin/go-migrate/templates"
	"github.com/spf13/cobra"
)

func toCamelCase(s string) string {
	ss := strings.Split(s, "_")
	sss := ""
	for _, v := range ss {
		sss += strings.Title(v)
	}
	return sss
}

func newMigration(filename string) error {
	var data []byte
	if strings.HasSuffix(filename, "_table") {
		if strings.HasPrefix(filename, "create_") {
			table := filename[7 : len(filename)-6]
			data = []byte(fmt.Sprintf(templates.MigrationCreateTemplate, config.Driver, toCamelCase(table), table))
		} else if index := strings.Index(filename, "to_"); index != -1 {
			table := filename[index+3 : len(filename)-6]
			data = []byte(fmt.Sprintf(templates.MigrationAlterTemplate, config.Driver, toCamelCase(filename), table))
		}
	}

	if len(data) == 0 {
		data = []byte(fmt.Sprintf(templates.MigrationTemplate, config.Driver, toCamelCase(filename)))
	}

	now := time.Now()
	filename = fmt.Sprintf("%d_%02d_%02d_%02d%02d%02d_%s.go",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), filename)

	if err := ioutil.WriteFile(
		fmt.Sprintf("migrations/%s", filename), data, 0666); err != nil {
		return fmt.Errorf("can't create migration file. \n%v", err)
	}

	return nil
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "New a migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called")
		filename := args[0]
		if err := newMigration(filename); err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Printf("new %s success.\n", filename)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
