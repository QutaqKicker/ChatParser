package dbHelper

// SelectBuildRequest запрос на создание select sql запроса
type SelectBuildRequest struct {
	// Filter фильтр запроса.
	Filter any
	// Sorts Названия колонок, которые нужно отсортировать и порядок сортировки
	Sorts []SortField
	// SelectType Тип селекта. Если выбран Special, его нужно задать вручную в SpecialSelect, иначе его можно оставить пустым. По умолчанию выбран All
	SelectType    SelectType
	SpecialSelect string
	// Take Сколько записей из запроса взять
	Take int
	// Skip Сколько записей из запроса пропустить
	Skip int
}

func NewRequest() *SelectBuildRequest {
	return &SelectBuildRequest{}
}

func (r *SelectBuildRequest) WithFilter(filter any) *SelectBuildRequest {
	r.Filter = filter
	return r
}

func (r *SelectBuildRequest) WithSorts(sorts []SortField) *SelectBuildRequest {
	r.Sorts = sorts
	return r
}

func (r *SelectBuildRequest) SetSelectType(selectType SelectType, specialSelect string) *SelectBuildRequest {
	r.SelectType = selectType
	r.SpecialSelect = specialSelect
	return r
}

func (r *SelectBuildRequest) NeedTake(take int) *SelectBuildRequest {
	r.Take = take
	return r
}

func (r *SelectBuildRequest) NeedSkip(skip int) *SelectBuildRequest {
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

type SortDirection string

const (
	Asc  SortDirection = "asc"
	Desc               = "desc"
)

type SortField struct {
	FieldName string
	Direction SortDirection
}
