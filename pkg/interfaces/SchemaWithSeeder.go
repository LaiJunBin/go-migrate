package interfaces

type SchemaWithSeeder interface {
	Create(table string, schemaFunc func(Blueprint)) Seeder
	Table(table string, schemaFunc func(Blueprint)) error
	DropIfExists(table string) error
}
