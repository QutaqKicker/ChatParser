package queryFilters

type UserFilter struct {
	commonFilter
	Name string `column:"text" relation:"="`
}
