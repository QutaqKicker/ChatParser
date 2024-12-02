package queryFilters

type MessageFilter struct {
	commonFilter
	SubText string   `column:"text" relation:"like"`
	UserIds []string `column:"user_id" relation:"in"`
	ChatIds []int    `column:"chat_id" relation:"in"`
}
