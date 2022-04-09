package mysql

import (
	"fmt"

	sk "github.com/laijunbin/go-solve-kit"
)

type meta struct {
	Name          string
	Type          string
	Length        int
	Nullable      bool
	Primary       bool
	AutoIncrement bool
	Default       interface{}
}

type operation interface {
	generateSql(table string, metadata []*meta) []string
}

type createOperation struct{}
type alterOperation struct{}

func (o *createOperation) generateSql(table string, metadata []*meta) []string {
	columns := sk.FromInterfaceArray(metadata)
	proceedColumns := columns.Map(func(v sk.Type, i int) interface{} {
		m := v.ValueOf().(*meta)

		if m.Type == "DROP" {
			return nil
		}

		s := fmt.Sprintf("`%s` %s", m.Name, m.Type)

		if m.Length != 0 {
			s += fmt.Sprintf("(%d)", m.Length)
		}

		if !m.Nullable {
			s += " NOT NULL"
		}

		if m.AutoIncrement {
			s += " AUTO_INCREMENT"
		}

		if m.Primary {
			s += fmt.Sprintf(", PRIMARY KEY (`%s`)", m.Name)
		}

		if m.Default != nil {
			s += fmt.Sprintf(" DEFAULT %v", m.Default)
		}

		return s
	})

	columnsStr := proceedColumns.Filter(func(s sk.Type, i int) bool {
		return s.ValueOf() != nil
	}).ToStringArray().Join(", ").ValueOf()
	sql := fmt.Sprintf("CREATE TABLE `%s` (%s);", table, columnsStr)
	return []string{sql}
}

func (o *alterOperation) generateSql(table string, metadata []*meta) []string {
	columns := sk.FromInterfaceArray(metadata)
	sql := fmt.Sprintf("ALTER TABLE `%s` %s;", table, columns.Map(func(v sk.Type, i int) interface{} {
		m := v.ValueOf().(*meta)

		if m.Type == "DROP" {
			return fmt.Sprintf("DROP `%s`", m.Name)
		}

		s := fmt.Sprintf("ADD `%s` %s", m.Name, m.Type)

		if m.Length != 0 {
			s += fmt.Sprintf("(%d)", m.Length)
		}

		if !m.Nullable {
			s += " NOT NULL"
		}

		if m.AutoIncrement {
			s += " AUTO_INCREMENT"
		}

		if m.Primary {
			s += " PRIMARY KEY"
		}

		if m.Default != nil {
			s += fmt.Sprintf(" DEFAULT %v", m.Default)
		}

		return s
	}).ToStringArray().Join(", ").ValueOf())

	return []string{
		sql,
	}
}

var metaOperations = struct {
	CREATE operation
	ALTER  operation
}{
	CREATE: &createOperation{},
	ALTER:  &alterOperation{},
}
