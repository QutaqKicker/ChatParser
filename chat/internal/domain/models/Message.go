package models

import "time"

type Message struct {
	Id      int
	ChatId  int `column:"chat_id"`
	UserId  int `column:"user_id"`
	Text    string
	Created time.Time
}
