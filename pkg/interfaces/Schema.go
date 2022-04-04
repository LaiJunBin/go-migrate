package interfaces

type Schema interface {
	Create(table string, schemaFunc func(Blueprint)) error
	Table(table string, schemaFunc func(Blueprint)) error
	DropIfExists(table string) error
}
