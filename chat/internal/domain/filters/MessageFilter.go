package filters

import "time"

type MessageFilter struct {
	Id             int       `column:"id" relation:"="`
	MinCreatedDate time.Time `column:"created" relation:"<"`
	MaxCreatedDate time.Time `column:"created" relation:">"`
	SubText        string    `column:"text" relation:"like"`
	UserId         string    `column:"user_id" relation:"="`
	UserIds        []string  `column:"user_id" relation:"in"`
	ChatIds        []int     `column:"chat_id" relation:"in"`
}

func NewMessageFilter() *MessageFilter {
	return &MessageFilter{}
}

func (f *MessageFilter) WhereId(value int) *MessageFilter {
	f.Id = value
	return f
}

func (f *MessageFilter) WhereMinCreatedDate(value time.Time) *MessageFilter {
	f.MinCreatedDate = value
	return f
}

func (f *MessageFilter) WhereMaxCreatedDate(value time.Time) *MessageFilter {
	f.MaxCreatedDate = value
	return f
}

func (f *MessageFilter) WhereSubText(value string) *MessageFilter {
	f.SubText = value
	return f
}

func (f *MessageFilter) WhereUserId(value string) *MessageFilter {
	f.UserId = value
	return f
}

func (f *MessageFilter) WhereUserIds(value []string) *MessageFilter {
	f.UserIds = value
	return f
}

func (f *MessageFilter) WhereChatIds(value []int) *MessageFilter {
	f.ChatIds = value
	return f
}
