package models

import "time"

type Message struct {
	Id               int32
	ChatId           int32    `column:"chat_id"`
	UserId           string   `column:"user_id"`
	ReplyToMessageId int32    `column:"reply_to_message_id"`
	RepliedMessage   *Message `not-mapped:"true"`
	Text             string
	Created          time.Time
}

func (m *Message) TableName() string {
	return "messages"
}
