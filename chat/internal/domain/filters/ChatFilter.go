package filters

import "time"

type MessageFilter struct {
	Id            int
	MinDate       time.Time
	MaxDate       time.Time
	SubText       string
	UserIds       []int
	ChatIds       []int
	Sorts         []string
	SpecifySelect string
}
