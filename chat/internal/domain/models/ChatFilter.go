package models

import "time"

type MessageFilter struct {
	MinDate       time.Time
	MaxDate       time.Time
	SubText       string
	UserIds       []int
	ChatIds       []int
	Sorts         []string
	SpecifySelect string
}
