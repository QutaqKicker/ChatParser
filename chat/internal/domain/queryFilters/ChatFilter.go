package queryFilters

import "time"

type ChatFilter struct {
	Id             int       `column:"id" relation:"="`
	MinCreatedDate time.Time `column:"created" relation:"<"`
	MaxCreatedDate time.Time `column:"created" relation:">"`
	Name           string    `column:"text" relation:"="`
}
