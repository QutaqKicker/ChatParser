package models

import "time"

type ChatFilter struct {
	minDate time.Time
	maxDate time.Time
	userIds []string
}
