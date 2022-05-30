package mysql

import "fmt"

type foreignMeta struct {
	Reference string
	Table     string
	OnUpdate  string
	OnDelete  string
}

func (fm *foreignMeta) generateSql(table string, name string) string {
	s := fmt.Sprintf("CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`)",
		fmt.Sprintf("fk_%s_%s", table, name),
		name,
		fm.Table,
		fm.Reference,
	)

	if fm.OnUpdate != "" {
		s += fmt.Sprintf(" ON UPDATE %s", fm.OnUpdate)
	}

	if fm.OnDelete != "" {
		s += fmt.Sprintf(" ON DELETE %s", fm.OnDelete)
	}

	return s
}
