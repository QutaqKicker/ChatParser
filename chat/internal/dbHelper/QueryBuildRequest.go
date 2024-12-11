package dbHelper

// QueryBuildRequest запрос на создание select sql запроса
type QueryBuildRequest struct {
	Filter        any
	Sorter        []string
	SelectType    SelectType
	SpecialSelect string
}

type SelectType int8

const (
	All SelectType = iota
	Count
	Sum
	Special
)
