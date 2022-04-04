package mysql

import (
	"fmt"
	"strings"

	"github.com/laijunbin/go-migrate/pkg/interfaces"
)

type Blueprint struct {
	alters      []string
	columns     []string
	dropColumns []string
}

func NewBlueprint() interfaces.Blueprint {
	return &Blueprint{}
}

func (bp *Blueprint) Id(name string, length int) {
	bp.columns = append(bp.columns, fmt.Sprintf("`%s` INT(%d) NOT NULL", name, length))
	bp.alters = append(bp.alters, fmt.Sprintf("ADD PRIMARY KEY (`%s`)", name))
	bp.alters = append(bp.alters, fmt.Sprintf("MODIFY `%s` INT(%d) NOT NULL AUTO_INCREMENT", name, length))
}

func (bp *Blueprint) String(name string, length int) interfaces.Blueprint {
	bp.columns = append(bp.columns, fmt.Sprintf("`%s` VARCHAR(%d) NOT NULL", name, length))
	return bp
}

func (bp *Blueprint) Text(name string) interfaces.Blueprint {
	bp.columns = append(bp.columns, fmt.Sprintf("`%s` TEXT NOT NULL", name))
	return bp
}

func (bp *Blueprint) Integer(name string, length int) interfaces.Blueprint {
	bp.columns = append(bp.columns, fmt.Sprintf("`%s` INT(%d) NOT NULL", name, length))
	return bp
}

func (bp *Blueprint) Date(name string) interfaces.Blueprint {
	bp.columns = append(bp.columns, fmt.Sprintf("`%s` DATE NOT NULL", name))
	return bp
}

func (bp *Blueprint) Boolean(name string) interfaces.Blueprint {
	bp.columns = append(bp.columns, fmt.Sprintf("`%s` TINYINT NOT NULL", name))
	return bp
}

func (bp *Blueprint) DateTime(name string) interfaces.Blueprint {
	bp.columns = append(bp.columns, fmt.Sprintf("`%s` DATETIME NOT NULL", name))
	return bp
}

func (bp *Blueprint) Timestamps() {
	bp.columns = append(bp.columns, "`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` DATETIME NULL")
}

func (bp *Blueprint) Nullable() interfaces.Blueprint {
	bp.columns[len(bp.columns)-1] = strings.Replace(bp.columns[len(bp.columns)-1], "NOT NULL", "", 1)
	return bp
}

func (bp *Blueprint) Default(value interface{}) interfaces.Blueprint {
	bp.columns[len(bp.columns)-1] += fmt.Sprintf(" DEFAULT '%v'", value)
	return bp
}

func (bp *Blueprint) DropColumn(column string) {
	bp.dropColumns = append(bp.dropColumns, fmt.Sprintf("`%s`", column))
}

func (bp *Blueprint) getAlters() []string {
	return bp.alters
}

func (bp *Blueprint) getColumns() []string {
	return bp.columns
}

func (bp *Blueprint) getDropColumns() []string {
	return bp.dropColumns
}
