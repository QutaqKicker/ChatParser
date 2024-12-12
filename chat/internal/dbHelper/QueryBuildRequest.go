package dbHelper

// QueryBuildRequest запрос на создание select sql запроса
type QueryBuildRequest struct {
	// Filter фильтр запроса.
	Filter any
	// SortColumnIndexes индексы колонок, по которым нужно отсортировать запрос. Если число отрицательное, desc, иначе asc
	SortColumnIndexes []int
	// SelectType Тип селекта. Если выбран Special, его нужно задать вручную в SpecialSelect, иначе его можно оставить пустым. По умолчанию выбран All
	SelectType    SelectType
	SpecialSelect string
	// Take Сколько записей из запроса взять
	Take int
	// Skip Сколько записей из запроса пропустить
	Skip int
}

func NewRequest() *QueryBuildRequest {
	return &QueryBuildRequest{}
}

func (r *QueryBuildRequest) WithFilter(filter any) *QueryBuildRequest {
	r.Filter = filter
	return r
}

func (r *QueryBuildRequest) WithSorts(sortColumnIndexes []int) *QueryBuildRequest {
	r.SortColumnIndexes = sortColumnIndexes
	return r
}

func (r *QueryBuildRequest) SetSelectType(selectType SelectType, specialSelect string) *QueryBuildRequest {
	r.SelectType = selectType
	r.SpecialSelect = specialSelect
	return r
}

func (r *QueryBuildRequest) NeedTake(take int) *QueryBuildRequest {
	r.Take = take
	return r
}

func (r *QueryBuildRequest) NeedSkip(skip int) *QueryBuildRequest {
	r.Skip = skip
	return r
}

type SelectType int8

const (
	All SelectType = iota
	Count
	Sum
	Special
)
