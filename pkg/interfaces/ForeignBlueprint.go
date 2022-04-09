package interfaces

type ForeignBlueprint interface {
	Reference(name string) ForeignBlueprint
	On(table string) ForeignBlueprint
	OnUpdate(action string) ForeignBlueprint
	OnDelete(action string) ForeignBlueprint
}
