package templates

var UserMigrationTemplate = `package migrations

import (
	"github.com/laijunbin/go-migrate/config"
	"github.com/laijunbin/go-migrate/pkg/interfaces"
	"github.com/laijunbin/go-migrate/pkg/lib/%[1]s"
)

func init() {
	config.Migrations = append(config.Migrations, CreateUsersTable())
}

type UsersTable struct{}

func CreateUsersTable() interfaces.Migration {
	return &UsersTable{}
}

func (t *UsersTable) Up() error {
	return %[1]s.Schema.Create("users", func(table interfaces.Blueprint) {
		table.Id("id", 10)
		table.String("username", 100)
		table.String("password", 100)
		table.Timestamps()
	})
}

func (t *UsersTable) Down() error {
	return %[1]s.Schema.DropIfExists("users")
}
`
