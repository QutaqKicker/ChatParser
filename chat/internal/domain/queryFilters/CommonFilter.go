package queryFilters

import "time"

type commonFilter struct {
	Id             int       `column:"id" relation:"="`
	MinCreatedDate time.Time `column:"created" relation:"<"`
	MaxCreatedDate time.Time `column:"created" relation:">"`
}
