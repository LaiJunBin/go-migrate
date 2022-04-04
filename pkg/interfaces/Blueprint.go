package interfaces

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
