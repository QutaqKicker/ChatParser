package models

import (
	"github.com/QutaqKicker/ChatParser/common/dbHelper"
	"time"
)

type Message struct {
	Id               int32
	ChatId           int32    `column:"chat_id"`
	ChatName         string   `not-mapped:"true"`
	UserId           string   `column:"user_id"`
	UserName         string   `not-mapped:"true"`
	ReplyToMessageId int32    `column:"reply_to_message_id"`
	RepliedMessage   *Message `not-mapped:"true"`
	Text             string
	Created          time.Time
}

func (m Message) TableName() string {
	return "messages"
}

// FieldValuesAsArray / GetKeyByName fields values as any array. Needs for safe inserting
func (m Message) FieldValuesAsArray() []any {
	return []any{m.Id, m.ChatId, m.UserId, m.ReplyToMessageId, m.Text, m.Created}
}

type MessageSorter []dbHelper.SortField

func (m Message) Sorter() *MessageSorter {
	return &MessageSorter{}
}

func (s *MessageSorter) ByCreated(direction dbHelper.SortDirection) *MessageSorter {
	*s = append(*s, dbHelper.SortField{FieldName: "created", Direction: direction})
	return s
}

func (s *MessageSorter) ByUser(direction dbHelper.SortDirection) *MessageSorter {
	*s = append(*s, dbHelper.SortField{FieldName: "user_id", Direction: direction})
	return s
}
