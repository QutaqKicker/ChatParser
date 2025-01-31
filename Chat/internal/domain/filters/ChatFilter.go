package filters

import "time"

type ChatFilter struct {
	Id             int32     `column:"id" relation:"="`
	MinCreatedDate time.Time `column:"created" relation:">"`
	MaxCreatedDate time.Time `column:"created" relation:"<"`
	Name           string    `column:"text" relation:"="`
}

func NewChatFilter() *ChatFilter {
	return &ChatFilter{}
}

func (f *ChatFilter) WhereId(value int32) *ChatFilter {
	f.Id = value
	return f
}

func (f *ChatFilter) WhereMinCreatedDate(value time.Time) *ChatFilter {
	f.MinCreatedDate = value
	return f
}

func (f *ChatFilter) WhereMaxCreatedDate(value time.Time) *ChatFilter {
	f.MaxCreatedDate = value
	return f
}

func (f *ChatFilter) WhereName(value string) *ChatFilter {
	f.Name = value
	return f
}
