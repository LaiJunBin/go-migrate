package model

type Seeder struct {
	Err error
}

func (s *Seeder) Error() string {
	if s.Err == nil {
		return ""
	}

	return s.Err.Error()
}

func NewSeeder(err error) *Seeder {
	return &Seeder{
		Err: err,
	}
}
