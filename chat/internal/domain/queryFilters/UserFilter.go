package queryFilters

import "time"

type UserFilter struct {
	Id             string    `column:"id" relation:"="`
	MinCreatedDate time.Time `column:"created" relation:"<"`
	MaxCreatedDate time.Time `column:"created" relation:">"`
	Name           string    `column:"text" relation:"="`
}
