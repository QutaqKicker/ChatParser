package filters

import "time"

type UserFilter struct {
	Id             string    `column:"id" relation:"="`
	MinCreatedDate time.Time `column:"created" relation:"<"`
	MaxCreatedDate time.Time `column:"created" relation:">"`
	Name           string    `column:"text" relation:"="`
}

func NewUserFilter() *UserFilter {
	return &UserFilter{}
}

func (f *UserFilter) WhereId(value string) *UserFilter {
	f.Id = value
	return f
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
