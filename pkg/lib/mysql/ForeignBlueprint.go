package mysql

import (
	"strings"

	"github.com/laijunbin/go-migrate/pkg/interfaces"
)

type foreignBlueprint struct {
	meta *foreignMeta
}

func newForeignBlueprint() interfaces.ForeignBlueprint {
	return &foreignBlueprint{
		meta: &foreignMeta{},
	}
}

func (fb *foreignBlueprint) Reference(name string) interfaces.ForeignBlueprint {
	fb.meta.Reference = name
	return fb
}

func (fb *foreignBlueprint) On(table string) interfaces.ForeignBlueprint {
	fb.meta.Table = table
	return fb
}

func (fb *foreignBlueprint) OnUpdate(action string) interfaces.ForeignBlueprint {
	fb.meta.OnUpdate = strings.ToUpper(action)
	return fb
}

func (fb *foreignBlueprint) OnDelete(action string) interfaces.ForeignBlueprint {
	fb.meta.OnDelete = strings.ToUpper(action)
	return fb
}
