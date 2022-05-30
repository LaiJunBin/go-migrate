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

func TestCreateUsersTable(t *testing.T) {
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
		"CREATE TABLE `users` (`id` INT(10) NOT NULL AUTO_INCREMENT, PRIMARY KEY (`id`), `description` TEXT, `amount` INT(10) NOT NULL DEFAULT '0', `name` VARCHAR(100) NOT NULL, `enable` TINYINT NOT NULL DEFAULT '0', `birthday` DATE NOT NULL, `last_login_at` DATETIME, `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` DATETIME DEFAULT NULL);",
	}

	sqls := driver.GetSqls()

	t.Logf("sqls: %v", sqls)

	checkSqlsMatch(t, sqls, expectedSqls)
}

func TestCreateProductsTable(t *testing.T) {
	driver, _ := mysql.NewMockDriver()
	defer checkDriverClosed(t, driver)

	schema := &mysql.Schema_test{}
	schema.Create(driver, "products", func(table interfaces.Blueprint) {
		table.String("id", 20)
		table.Primary("id")
		table.Integer("user_id", 10)
		table.Foreign("user_id").Reference("id").On("users").OnUpdate("cascade").OnDelete("cascade")
		table.Integer("category_id", 10).Index()
		table.Boolean("enable").Default(1)
		table.Timestamps()
	})

	expectedSqls := []string{
		"CREATE TABLE `products` (`id` VARCHAR(20) NOT NULL, PRIMARY KEY (`id`), `user_id` INT(10) NOT NULL, CONSTRAINT `fk_products_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE CASCADE ON DELETE CASCADE, `category_id` INT(10) NOT NULL, INDEX (`category_id`), `enable` TINYINT NOT NULL DEFAULT '1', `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` DATETIME DEFAULT NULL);",
	}

	sqls := driver.GetSqls()

	t.Logf("sqls: %v", sqls)

	checkSqlsMatch(t, sqls, expectedSqls)
}

func TestAlterUsersTable(t *testing.T) {
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
		"ALTER TABLE `users` ADD `name` INT(10) NOT NULL, ADD `price` VARCHAR(100) NOT NULL, DROP `description`, DROP `enable`;",
	}

	sqls := driver.GetSqls()

	t.Logf("sqls: %v", sqls)

	checkSqlsMatch(t, sqls, expectedSqls)
}

func TestAlterProductsTable(t *testing.T) {
	driver, _ := mysql.NewMockDriver()
	defer checkDriverClosed(t, driver)

	schema := &mysql.Schema_test{}
	schema.Table(driver, "products", func(table interfaces.Blueprint) {
		table.Integer("price", 10)
		table.DropPrimary()
		table.DropForeign("user_id")
		table.DropIndex("category_id")

		table.Index("user_id")
	})

	expectedSqls := []string{
		"ALTER TABLE `products` ADD `price` INT(10) NOT NULL, DROP PRIMARY KEY, DROP FOREIGN KEY `fk_products_user_id`, DROP INDEX `fk_products_user_id`, DROP INDEX `category_id`, ADD INDEX (`user_id`);",
	}

	sqls := driver.GetSqls()

	t.Logf("sqls: %v", sqls)

	checkSqlsMatch(t, sqls, expectedSqls)
}

func TestAlterTablePrimary(t *testing.T) {
	driver, _ := mysql.NewMockDriver()
	defer checkDriverClosed(t, driver)

	schema := &mysql.Schema_test{}
	schema.Table(driver, "products", func(table interfaces.Blueprint) {
		table.Id("auto_id", 10)
		table.Primary("id")
		table.Index("user_id")
		table.Unique("user_id")
	})

	expectedSqls := []string{
		"ALTER TABLE `products` ADD `auto_id` INT(10) NOT NULL AUTO_INCREMENT, ADD PRIMARY KEY (`auto_id`), ADD PRIMARY KEY (`id`), ADD INDEX (`user_id`), ADD UNIQUE (`user_id`);",
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

func TestCreateUsersTableAndSeed(t *testing.T) {
	driver, _ := mysql.NewMockDriver()
	defer checkDriverClosed(t, driver)

	schema := &mysql.Schema_test{}
	schema.Create(driver, "users", func(table interfaces.Blueprint) {
		table.Id("id", 10)
		table.String("username", 100)
		table.String("password", 100)
		table.Timestamps()
	}).Seed([]map[string]interface{}{
		{
			"username": "admin",
			"password": "1234",
		},
		{
			"username": "user01",
			"password": "1234",
		},
	}...)

	expectedSqls := []string{
		"CREATE TABLE `users` (`id` INT(10) NOT NULL AUTO_INCREMENT, PRIMARY KEY (`id`), `username` VARCHAR(100) NOT NULL, `password` VARCHAR(100) NOT NULL, `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` DATETIME DEFAULT NULL);",
		"INSERT INTO `users` (`password`, `username`) VALUES ('1234', 'admin');",
		"INSERT INTO `users` (`password`, `username`) VALUES ('1234', 'user01');",
	}

	sqls := driver.GetSqls()

	t.Logf("sqls: %v", sqls)

	checkSqlsMatch(t, sqls, expectedSqls)
}
