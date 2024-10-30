package filters

import "time"

type QueryBuildRequest struct {
	Filter        *MessageFilter
	Sorter        []string
	SelectType    SelectType
	SpecialSelect string
}

type MessageFilter struct {
	Id             int
	MinCreatedDate time.Time `column:"created" relation:"<"`
	MaxCreatedDate time.Time `column:"created" relation:">"`
	SubText        string    `column:"text" relation:"like"`
	UserIds        []int     `column:"user_id" relation:"in"`
	ChatIds        []int     `column:"user_id" relation:"in"`
}

type SelectType int8

const (
	All SelectType = iota
	Count
	Sum
	Special
)
