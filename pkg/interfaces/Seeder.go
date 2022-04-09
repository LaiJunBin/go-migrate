package interfaces

type Seeder interface {
	error
	Seed(data ...map[string]interface{}) error
}
