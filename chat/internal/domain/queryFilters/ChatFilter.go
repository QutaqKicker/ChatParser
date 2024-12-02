package queryFilters

type ChatFilter struct {
	commonFilter
	Name string `column:"text" relation:"="`
}
