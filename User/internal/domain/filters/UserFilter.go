package filters

import "time"

type UserFilter struct {
	MinCreatedDate time.Time `column:"created" relation:">"`
	MaxCreatedDate time.Time `column:"created" relation:"<"`
	Name           string    `column:"name" relation:"="`
}

func NewUserFilter() *UserFilter {
	return &UserFilter{}
}

func (f *UserFilter) WhereMinCreatedDate(value time.Time) *UserFilter {
	f.MinCreatedDate = value
	return f
}

func (f *UserFilter) WhereMaxCreatedDate(value time.Time) *UserFilter {
	f.MaxCreatedDate = value
	return f
}

func (f *UserFilter) WhereName(value string) *UserFilter {
	f.Name = value
	return f
}
