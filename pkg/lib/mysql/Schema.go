package mysql

import (
	"fmt"
	"strings"

	"github.com/laijunbin/go-migrate/pkg/interfaces"
)

type schema struct {
}

func newSchema() interfaces.Schema {
	return &schema{}
}

func (s *schema) Create(table string, schemaFunc func(interfaces.Blueprint)) error {
	driver, err := NewDriver()
	if err != nil {
		return err
	}
	defer driver.Close()

	blueprint := NewBlueprint().(*Blueprint)
	schemaFunc(blueprint)
	sql := fmt.Sprintf("CREATE TABLE `%s` (%s);", table, strings.Join(blueprint.getColumns(), ","))
	if _, err := driver.Execute(sql); err != nil {
		return err
	}

	for _, alter := range blueprint.getAlters() {
		sql := fmt.Sprintf("ALTER TABLE `%s` %s;", table, alter)
		if _, err := driver.Execute(sql); err != nil {
			return err
		}
	}

	return nil
}

func (s *schema) Table(table string, schemaFunc func(interfaces.Blueprint)) error {
	driver, err := NewDriver()
	if err != nil {
		return err
	}
	defer driver.Close()

	blueprint := NewBlueprint().(*Blueprint)
	schemaFunc(blueprint)

	columns := blueprint.getColumns()
	if len(columns) > 0 {
		sql := fmt.Sprintf("ALTER TABLE `%s` ADD %s;", table, strings.Join(columns, ", ADD "))
		if _, err := driver.Execute(sql); err != nil {
			return err
		}
	}

	dropColumns := blueprint.getDropColumns()
	if len(dropColumns) > 0 {
		sql := fmt.Sprintf("ALTER TABLE `%s` DROP %s;", table, strings.Join(dropColumns, ", DROP "))
		if _, err := driver.Execute(sql); err != nil {
			return err
		}
	}

	for _, alter := range blueprint.getAlters() {
		sql := fmt.Sprintf("ALTER TABLE `%s` %s;", table, alter)
		if _, err := driver.Execute(sql); err != nil {
			return err
		}
	}

	return nil
}

func (s *schema) DropIfExists(table string) error {
	driver, err := NewDriver()
	if err != nil {
		return err
	}
	defer driver.Close()

	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", table)
	_, err = driver.Execute(sql)
	return err
}

var Schema = newSchema()
