package models

import "time"

type ChatMessage struct {
	userId     string
	text       string
	createDate time.Time
}
