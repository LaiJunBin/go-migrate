package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MockDriver struct {
	sqls    []string
	isClose bool
}

func NewMockDriver() (*MockDriver, error) {
	return &MockDriver{isClose: false}, nil
}

func (d *MockDriver) Execute(sql string) (sql.Result, error) {
	d.sqls = append(d.sqls, sql)
	return nil, nil
}

func (d *MockDriver) Query(sql string) (*sql.Rows, error) {
	d.sqls = append(d.sqls, sql)
	return nil, nil
}

func (d *MockDriver) Select(dest interface{}, sql string) error {
	d.sqls = append(d.sqls, sql)
	return nil
}

func (d *MockDriver) Debug() *sqlx.DB {
	return nil
}

func (d *MockDriver) Close() error {
	d.isClose = true
	return nil
}

func (d *MockDriver) GetSqls() []string {
	return d.sqls
}

func (d *MockDriver) IsClose() bool {
	return d.isClose
}
