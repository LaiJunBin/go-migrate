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
	Unique        bool
	Index         bool
	Primary       bool
	AutoIncrement bool
	Default       interface{}
	Foreign       *foreignMeta
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

		if m.Foreign != nil {
			return m.Foreign.generateSql(m.Name)
		}

		s := ""
		if m.Type != "" {
			s += fmt.Sprintf("`%s` %s", m.Name, m.Type)
		}

		if m.Length != 0 {
			s += fmt.Sprintf("(%d)", m.Length)
		}

		if s != "" && !m.Nullable {
			s += " NOT NULL"
		}

		if m.AutoIncrement {
			s += " AUTO_INCREMENT"
		}

		if m.Primary {
			if s != "" {
				s += ", "
			}
			s += fmt.Sprintf("PRIMARY KEY (`%s`)", m.Name)
		}

		if m.Unique {
			if s != "" {
				s += ", "
			}
			s += fmt.Sprintf("UNIQUE (`%s`)", m.Name)
		}

		if m.Index {
			if s != "" {
				s += ", "
			}
			s += fmt.Sprintf("INDEX (`%s`)", m.Name)
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
			if m.Primary {
				return "DROP PRIMARY KEY"
			}

			if m.Index || m.Unique {
				return fmt.Sprintf("DROP INDEX `%s`", m.Name)
			}

			if m.Foreign != nil {
				return fmt.Sprintf("DROP FOREIGN KEY `%[1]s`, DROP INDEX `%[1]s`", fmt.Sprintf("fk_%s", m.Name))
			}

			return fmt.Sprintf("DROP `%s`", m.Name)
		}

		if m.Foreign != nil {
			return fmt.Sprintf("ADD %s", m.Foreign.generateSql(m.Name))
		}

		s := ""
		if m.Type != "" {
			s += fmt.Sprintf("ADD `%s` %s", m.Name, m.Type)
		}

		if m.Length != 0 {
			s += fmt.Sprintf("(%d)", m.Length)
		}

		if s != "" && !m.Nullable {
			s += " NOT NULL"
		}

		if m.AutoIncrement {
			s += " AUTO_INCREMENT"
		}

		if m.Primary {
			if s != "" {
				s += ", "
			}
			s += fmt.Sprintf("ADD PRIMARY KEY (`%s`)", m.Name)
		}

		if m.Default != nil {
			s += fmt.Sprintf(" DEFAULT %v", m.Default)
		}

		if m.Unique {
			if s != "" {
				s += ", "
			}
			s += fmt.Sprintf("ADD UNIQUE (`%s`)", m.Name)
		}

		if m.Index {
			if s != "" {
				s += ", "
			}
			s += fmt.Sprintf("ADD INDEX (`%s`)", m.Name)
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
