package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/laijunbin/go-migrate/config"
)

type driver struct {
	db *sqlx.DB
}

func NewDriver() (*driver, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		config.Config.Username,
		config.Config.Password,
		config.Config.Host,
		config.Config.Port,
		config.Config.Dbname,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect mysql server failed, err:%v", err)
	}

	return &driver{db: db}, nil
}

func (d *driver) Execute(sql string) (sql.Result, error) {
	result, err := d.db.Exec(sql)
	return result, err
}

func (d *driver) Query(sql string) (*sql.Rows, error) {
	return d.db.Query(sql)
}

func (d *driver) Select(dest interface{}, sql string) error {
	return d.db.Select(dest, sql)
}

func (d *driver) Debug() *sqlx.DB {
	return d.db
}

func (d *driver) Close() error {
	return d.db.Close()
}
