package mysql

import (
	"fmt"
	"strings"

	"github.com/laijunbin/go-migrate/pkg/interfaces"
	"github.com/laijunbin/go-migrate/pkg/model"
)

type migrator struct {
}

func InitMigrator() interfaces.Migrator {
	return &migrator{}
}

func (m *migrator) CheckTable() (bool, error) {
	driver, err := NewDriver()
	if err != nil {
		return false, err
	}
	defer driver.Close()

	sql := "SHOW TABLES LIKE 'migrations'"
	rows, err := driver.Query(sql)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (m *migrator) CreateTable() error {
	driver, err := NewDriver()
	if err != nil {
		return err
	}
	defer driver.Close()

	sqls := []string{
		"CREATE TABLE `migrations` (`id` int(10) UNSIGNED NOT NULL, `migration` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL, `batch` int(11) NOT NULL) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;",
		"ALTER TABLE `migrations` ADD PRIMARY KEY (`id`);",
		"ALTER TABLE `migrations` MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;",
	}

	for _, sql := range sqls {
		if _, err := driver.Execute(sql); err != nil {
			return err
		}
	}

	return nil
}

func (m *migrator) DropTableIfExists() error {
	driver, err := NewDriver()
	if err != nil {
		return err
	}
	defer driver.Close()

	sql := "DROP TABLE IF EXISTS migrations;"
	_, err = driver.Execute(sql)
	return err
}

func (m *migrator) DropAllTable() error {
	driver, err := NewDriver()
	if err != nil {
		return err
	}
	defer driver.Close()

	tables := []string{}
	sql := "SHOW TABLES"
	driver.Select(&tables, sql)

	sql = fmt.Sprintf("DROP TABLE IF EXISTS %s;", strings.Join(tables, ","))
	_, err = driver.Execute(sql)
	return err
}

func (m *migrator) GetMigrations() ([]model.Migration, error) {
	driver, err := NewDriver()
	if err != nil {
		return nil, err
	}
	defer driver.Close()

	migrations := []model.Migration{}
	sql := "SELECT id, migration, batch FROM `migrations`"
	err = driver.Select(&migrations, sql)
	return migrations, err
}

func (m *migrator) WriteRecord(migration string, batch int) error {
	driver, err := NewDriver()
	if err != nil {
		return err
	}
	defer driver.Close()

	sql := fmt.Sprintf("INSERT INTO `migrations`(`migration`, `batch`) VALUES ('%s','%d')", migration, batch)
	_, err = driver.Execute(sql)
	return err
}

func (m *migrator) DeleteRecord(id int) error {
	driver, err := NewDriver()
	if err != nil {
		return err
	}
	defer driver.Close()

	sql := fmt.Sprintf("DELETE FROM `migrations` WHERE id = %d", id)
	_, err = driver.Execute(sql)
	return err
}
