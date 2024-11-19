package db

type QueryBuildRequest[TFilter any] struct {
	Filter        *TFilter
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
