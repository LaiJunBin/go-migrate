package mysql

import (
	"fmt"

	"github.com/laijunbin/go-migrate/pkg/interfaces"
	sk "github.com/laijunbin/go-solve-kit"
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

func (bp *Blueprint) Unique(column ...string) interfaces.Blueprint {
	if len(column) == 0 {
		bp.metadata[len(bp.metadata)-1].Unique = true
	} else {
		for _, c := range column {
			bp.metadata = append(bp.metadata, &meta{
				Name:   c,
				Unique: true,
			})
		}
	}
	return bp
}

func (bp *Blueprint) Index(column ...string) interfaces.Blueprint {
	if len(column) == 0 {
		bp.metadata[len(bp.metadata)-1].Index = true
	} else {
		for _, c := range column {
			bp.metadata = append(bp.metadata, &meta{
				Name:  c,
				Index: true,
			})
		}
	}
	return bp
}

func (bp *Blueprint) Default(value interface{}) interfaces.Blueprint {
	bp.metadata[len(bp.metadata)-1].Default = fmt.Sprintf("'%v'", value)
	return bp
}

func (bp *Blueprint) Foreign(name string) interfaces.ForeignBlueprint {
	fb := newForeignBlueprint().(*foreignBlueprint)
	bp.metadata = append(bp.metadata, &meta{
		Name:    name,
		Foreign: fb.meta,
	})
	return fb
}

func (bp *Blueprint) Primary(name ...string) interfaces.Blueprint {
	bp.metadata = append(bp.metadata, &meta{
		Name:    sk.FromStringArray(name).Join("`, `").ValueOf(),
		Primary: true,
	})
	return bp
}

func (bp *Blueprint) DropColumn(column string) {
	bp.metadata = append(bp.metadata, &meta{
		Name: column,
		Type: "DROP",
	})
}

func (bp *Blueprint) DropUnique(name string) {
	bp.metadata = append(bp.metadata, &meta{
		Name:   name,
		Type:   "DROP",
		Unique: true,
	})
}
func (bp *Blueprint) DropIndex(name string) {
	bp.metadata = append(bp.metadata, &meta{
		Name:  name,
		Type:  "DROP",
		Index: true,
	})
}
func (bp *Blueprint) DropForeign(name string) {
	bp.metadata = append(bp.metadata, &meta{
		Name:    name,
		Type:    "DROP",
		Foreign: newForeignBlueprint().(*foreignBlueprint).meta,
	})
}
func (bp *Blueprint) DropPrimary() {
	bp.metadata = append(bp.metadata, &meta{
		Type:    "DROP",
		Primary: true,
	})
}

func (bp *Blueprint) GetSqls(table string, operation operation) []string {
	return operation.generateSql(table, bp.metadata)
}
