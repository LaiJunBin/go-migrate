package mysql

import (
	"fmt"

	"github.com/laijunbin/go-migrate/pkg/interfaces"
)

type Blueprint struct {
	metadata []*meta
}

func NewBlueprint() interfaces.Blueprint {
	return &Blueprint{}
}

func (bp *Blueprint) Id(name string, length int) {
	bp.metadata = append(bp.metadata, &meta{
		Name:          name,
		Type:          "INT",
		Length:        length,
		AutoIncrement: true,
		Primary:       true,
	})
}

func (bp *Blueprint) String(name string, length int) interfaces.Blueprint {
	bp.metadata = append(bp.metadata, &meta{
		Name:   name,
		Type:   "VARCHAR",
		Length: length,
	})
	return bp
}

func (bp *Blueprint) Text(name string) interfaces.Blueprint {
	bp.metadata = append(bp.metadata, &meta{
		Name: name,
		Type: "TEXT",
	})
	return bp
}

func (bp *Blueprint) Integer(name string, length int) interfaces.Blueprint {
	bp.metadata = append(bp.metadata, &meta{
		Name:   name,
		Type:   "INT",
		Length: length,
	})
	return bp
}

func (bp *Blueprint) Date(name string) interfaces.Blueprint {
	bp.metadata = append(bp.metadata, &meta{
		Name: name,
		Type: "DATE",
	})
	return bp
}

func (bp *Blueprint) Boolean(name string) interfaces.Blueprint {
	bp.metadata = append(bp.metadata, &meta{
		Name: name,
		Type: "TINYINT",
	})
	return bp
}

func (bp *Blueprint) DateTime(name string) interfaces.Blueprint {
	bp.metadata = append(bp.metadata, &meta{
		Name: name,
		Type: "DATETIME",
	})
	return bp
}

func (bp *Blueprint) Timestamps() {
	bp.metadata = append(bp.metadata, &meta{
		Name:    "created_at",
		Type:    "DATETIME",
		Default: "CURRENT_TIMESTAMP",
	})

	bp.metadata = append(bp.metadata, &meta{
		Name:     "updated_at",
		Type:     "DATETIME",
		Nullable: true,
		Default:  "NULL",
	})
}

func (bp *Blueprint) Nullable() interfaces.Blueprint {
	bp.metadata[len(bp.metadata)-1].Nullable = true
	return bp
}

func (bp *Blueprint) Default(value interface{}) interfaces.Blueprint {
	bp.metadata[len(bp.metadata)-1].Default = fmt.Sprintf("'%v'", value)
	return bp
}

func (bp *Blueprint) DropColumn(column string) {
	bp.metadata = append(bp.metadata, &meta{
		Name: column,
		Type: "DROP",
	})
}

func (bp *Blueprint) GetSqls(table string, operation operation) []string {
	return operation.generateSql(table, bp.metadata)
}
