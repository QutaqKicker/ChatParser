package queryFilters

import "time"

type ChatFilter struct {
	Id         int       `column:"id" relation:"="`
	Name       string    `column:"text" relation:"="`
	MinCreated time.Time `column:"created" relation:"<"`
	MaxCreated time.Time `column:"created" relation:">"`
}
