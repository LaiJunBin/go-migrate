package mysql

import (
	"fmt"
	"sort"

	"github.com/laijunbin/go-migrate/pkg/interfaces"
	mysql_interfaces "github.com/laijunbin/go-migrate/pkg/lib/mysql/interfaces"
	"github.com/laijunbin/go-migrate/pkg/model"
	sk "github.com/laijunbin/go-solve-kit"
)

type seeder struct {
	*model.Seeder
	table string
}

type Seeder_test struct {
	*model.Seeder
	table  string
	driver mysql_interfaces.Driver
}

func NewSeeder(table string, err error) interfaces.Seeder {
	return &seeder{
		Seeder: model.NewSeeder(err),
		table:  table,
	}
}

func NewTestSeeder(driver mysql_interfaces.Driver, table string, err error) interfaces.Seeder {
	return &Seeder_test{
		Seeder: model.NewSeeder(err),
		table:  table,
		driver: driver,
	}
}

func (s *Seeder_test) Seed(data ...map[string]interface{}) error {
	return runSeed(s.driver, s.table, data...)
}

func (s *seeder) Seed(data ...map[string]interface{}) error {
	if s.Error() != "" {
		return s.Err
	}

	driver, err := NewDriver()
	if err != nil {
		return err
	}

	return runSeed(driver, s.table, data...)
}

func runSeed(driver mysql_interfaces.Driver, table string, data ...map[string]interface{}) error {
	for i := 0; i < len(data); i++ {
		keys := []string{}
		values := []string{}
		for k, _ := range data[i] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for j, k := range keys {
			values = append(values, fmt.Sprintf("'%s'", data[i][k]))
			keys[j] = fmt.Sprintf("`%s`", keys[j])
		}

		sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s);",
			table,
			sk.FromStringArray(keys).Join(", "),
			sk.FromStringArray(values).Join(", "),
		)

		if _, err := driver.Execute(sql); err != nil {
			return err
		}
	}

	return nil
}
