# Go-Migrate

> Make your database migrations as easy as Laravel.

# Contents
* [Install](#install)
* [Quick start](#quick-start)
* [Commands](#commands)
* [API](#api)
* [These APIs how interact?](#these-apis-how-interact)

# Install
```
$ go install github.com/laijunbin/go-migrate@latest
```

# Quick start

Let us quickly start trying to use go-migrate to migrate your database, which uses mysql as an example.

> Assume you already installed the package.

1. Create an empty directory for the go project and enter it.
```
$ mkdir go-app && cd go-app
```

2. Init go module.
```
$ go mod init <module-path>
```

3. Init go-migrate
```
$ go-migrate init mysql
```

4. Open `cmd/migrate.go` to set your database config or immediately execute migrate operation.
```
$ go-migrate migrate
```

5. Finalize, you can see the users table in your mysql database.

# Commands

The available commands are as follows:
* [init](#init)
* [new](#new)
* [migrate](#migrate)
* [rollback](#rollback)
* [reset](#reset)
* [fresh](#fresh)
* [refresh](#refresh)


## init
```
$ go-migrate init <db>
```

`<db>` currently support the following:
* mysql

The `init` command will init go-migrate context in your project and create a sample migration file like below.
```go
// omit...
type UsersTable struct{}

func CreateUsersTable() interfaces.Migration {
	return &UsersTable{}
}

func (t *UsersTable) Up() error {
	return mysql.Schema.Create("users", func(table interfaces.Blueprint) {
		table.Id("id", 10)
		table.String("username", 100)
		table.String("password", 100)
		table.Timestamps()
	})
}

func (t *UsersTable) Down() error {
	return mysql.Schema.DropIfExists("users")
}
```

Also, you can modify DatabaseConfig in `cmd/migrate.go`.

The default like follows:
```go
config.DatabaseConfig{
    Host:     "127.0.0.1",
    Port:     3306,
    Username: "root",
    Password: "",
    Dbname:   "test",
}
```


## new
```
$ go-migrate new <filename>
```

The `new` command will create a migration file, filename has the following rules:

* create_`<table>`_table: Generate a migration file, Up() includes `Schema.Create` and Down() includes `Schema.DropIfExists`.
* xxx_to_`<table>`_table: Generate a migration file, Up() and Down() includes `Schema.Table` both.
* otherwise: Generate a bare migration file which only Up() and Down().

## migrate
```
$ go-migrate migrate
```

The `migrate` command will execute all migrate operations, that except already executed migrate and record the state to the database.

## rollback
```
$ go-migrate rollback
```

The `rollback` command will rollback migrate of the one batch.

## reset
```
$ go-migrate reset
```

The `reset` command will rollback all migrate operations.

## fresh
```
$ go-migrate fresh
```

The `fresh` command will drop all tables and re-run all migrations.

## refresh
```
$ go-migrate refresh
```

The `refresh` command will rollback all migrate operations and re-run all migrations.

# API
* [Model](#model)
	* [Migration](#migration-model)
* [Interface](#interface)
	* [Migration](#migration-interface)
	* [Migrator](#migrator)
	* [Schema](#schema)
	* [Blueprint](#blueprint)

# Model

## Migration Model
```go
type Migration struct {
	Id        int
	Migration string
	Batch     int
}
```

The model represents a row of migration records.

# Interface

## Migration Interface
```go
type Migration interface {
	Up() error
	Down() error
}
```

The interface defines that every migration file should implement the function.
* Up: When executing the `migrate` command.
* Down: When executing the `rollback` command.

## Migrator
```go
type Migrator interface {
	CheckTable() (bool, error)
	CreateTable() error
	DropTableIfExists() error
	DropAllTable() error
	GetMigrations() ([]model.Migration, error)
	WriteRecord(migration string, batch int) error
	DeleteRecord(id int) error
}
```

The interface defines follows: 
> Assume use `migrations` as the table for storing migration records.
* CheckTable: Check if the `migrations` table exists.
* CreateTable: Create the `migrations` table.
* DropTableIfExists: Drop the `migrations` table
* DropAllTable: Drop all table.
* GetMigrations: Get all migrations.
* WriteRecord: Write a record to `migrations` table.
* DeleteRecord: Delete a record of `migration` table.

## Schema
```go
type Schema interface {
	Create(table string, schemaFunc func(Blueprint)) error
	Table(table string, schemaFunc func(Blueprint)) error
	DropIfExists(table string) error
}
```

The interface defines follows: 
* Create: Define how to create the table and execute it.
* Table: Define how to alter the table and execute it.
* DropIfExists: Drop table if exists.

## Blueprint
```go
type Blueprint interface {
	Id(name string, length int)
	String(name string, length int) Blueprint
	Text(name string) Blueprint
	Integer(name string, length int) Blueprint
	Date(name string) Blueprint
	Boolean(name string) Blueprint
	DateTime(name string) Blueprint
	Nullable() Blueprint
	Default(value interface{}) Blueprint
	DropColumn(column string)
	Timestamps()
}
```

The interface defines follows: 
* Id: Create an `int` equivalent column and auto increment.
* String: Create a `string` equivalent column.
* Text: Create a `text` equivalent column.
* Integer: Create an `int` equivalent column.
* Date: Create a `date` equivalent column.
* Boolean: Create a `boolean` equivalent column.
* DateTime: Create a `datetime` equivalent column.
* Nullable: The columns that are created will be `nullable`
* Default: Set a default value to the column.
* DropColumn: Drop a column.
* Timestamps: Create created_at field with the current time as default and updated_at TIMESTAMP equivalent columns.  

# These APIs how interact?
* `Command` uses `Migrator` to check migrations status and calls `Migration`.
* `Migration` use `Schema` in Up() and Down().
* `Schema` called `Blueprint` to get the blueprint and execute it to alter the database.
* `Blueprint` defines how to generate columns.

## For user
You just need to focus on the Up() and Down() function, so you just need to understand that `Blueprint` and `Schema`.

## For developer
You need to create a directory in the `pkg/lib/` and implement the above function.