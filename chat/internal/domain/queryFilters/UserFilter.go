package queryFilters

import "time"

type UserFilter struct {
	Id         string    `column:"id" relation:"="`
	Name       string    `column:"text" relation:"="`
	MinCreated time.Time `column:"created" relation:"<"`
	MaxCreated time.Time `column:"created" relation:">"`
}
