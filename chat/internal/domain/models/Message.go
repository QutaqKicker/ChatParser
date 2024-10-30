package models

import "time"

type Message struct {
	Id      int32
	ChatId  int32 `column:"chat_id"`
	UserId  int32 `column:"user_id"`
	Text    string
	Created time.Time
}

func (m *Message) TableName() string {
	return "messages"
}
