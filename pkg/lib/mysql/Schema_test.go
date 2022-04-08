package mysql_test

import (
	"testing"

	"github.com/laijunbin/go-migrate/pkg/interfaces"

	"github.com/laijunbin/go-migrate/pkg/lib/mysql"
	sk "github.com/laijunbin/go-solve-kit"
)

func checkDriverClosed(t *testing.T, driver *mysql.MockDriver) {
	if !driver.IsClose() {
		t.Fatal("driver not closed!")
	}
}

func checkSqlsMatch(t *testing.T, sqls []string, expectedSqls []string) {
	if len(sqls) != len(expectedSqls) {
		t.Fatal("sqls length and expected sqls length not match!")
	}

	sk.FromStringArray(sqls).ForEach(func(s sk.String, i int) {
		if expectedSqls[i] != s.ValueOf() {
			t.Fatalf("unknow sql:\nresult: %v\nexpect: %v", s, expectedSqls[i])
		}
	})
}

func TestCreate(t *testing.T) {
	driver, _ := mysql.NewMockDriver()
	defer checkDriverClosed(t, driver)

	schema := &mysql.Schema_test{}
	schema.Create(driver, "users", func(table interfaces.Blueprint) {
		table.Id("id", 10)
		table.Text("description").Nullable()
		table.Integer("amount", 10).Default(0)
		table.String("name", 100)
		table.Boolean("enable").Default(0)
		table.Date("birthday")
		table.DateTime("last_login_at").Nullable()
		table.Timestamps()
	})

	expectedSqls := []string{
		"CREATE TABLE `users` (`id` INT(10) NOT NULL,`description` TEXT ,`amount` INT(10) NOT NULL DEFAULT '0',`name` VARCHAR(100) NOT NULL,`enable` TINYINT NOT NULL DEFAULT '0',`birthday` DATE NOT NULL,`last_login_at` DATETIME ,`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` DATETIME NULL);",
		"ALTER TABLE `users` ADD PRIMARY KEY (`id`);",
		"ALTER TABLE `users` MODIFY `id` INT(10) NOT NULL AUTO_INCREMENT;",
	}

	sqls := driver.GetSqls()

	t.Logf("sqls: %v", sqls)

	checkSqlsMatch(t, sqls, expectedSqls)
}

func TestTable(t *testing.T) {
	driver, _ := mysql.NewMockDriver()
	defer checkDriverClosed(t, driver)

	schema := &mysql.Schema_test{}
	schema.Table(driver, "users", func(table interfaces.Blueprint) {
		table.Integer("name", 10)
		table.String("price", 100)
		table.DropColumn("description")
		table.DropColumn("enable")
	})

	expectedSqls := []string{
		"ALTER TABLE `users` ADD `name` INT(10) NOT NULL, ADD `price` VARCHAR(100) NOT NULL;",
		"ALTER TABLE `users` DROP `description`, DROP `enable`;",
	}

	sqls := driver.GetSqls()

	t.Logf("sqls: %v", sqls)

	checkSqlsMatch(t, sqls, expectedSqls)
}

func TestDropIfExists(t *testing.T) {
	driver, _ := mysql.NewMockDriver()
	defer checkDriverClosed(t, driver)

	schema := &mysql.Schema_test{}
	schema.DropIfExists(driver, "users")

	expectedSqls := []string{
		"DROP TABLE IF EXISTS users;",
	}

	sqls := driver.GetSqls()

	t.Logf("sqls: %v", sqls)

	checkSqlsMatch(t, sqls, expectedSqls)
}
