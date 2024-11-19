package queryFilters

import "time"

type MessageFilter struct {
	Id             int       `column:"id" relation:"="`
	MinCreatedDate time.Time `column:"created" relation:"<"`
	MaxCreatedDate time.Time `column:"created" relation:">"`
	SubText        string    `column:"text" relation:"like"`
	UserIds        []int     `column:"user_id" relation:"in"`
	ChatIds        []int     `column:"user_id" relation:"in"`
}
