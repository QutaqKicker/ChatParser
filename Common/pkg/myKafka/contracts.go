package myKafka

import "time"

type CreateLogRequest struct {
	ServiceName string
	Type        int
	Message     string
	Created     time.Time
}

type UserMessageCountRequest struct {
	UserName     string
	MessageCount int
}
