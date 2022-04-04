package interfaces

import "github.com/laijunbin/go-migrate/pkg/model"

type Migrator interface {
	CheckTable() (bool, error)
	CreateTable() error
	DropTableIfExists() error
	DropAllTable() error
	GetMigrations() ([]model.Migration, error)
	WriteRecord(migration string, batch int) error
	DeleteRecord(id int) error
}
