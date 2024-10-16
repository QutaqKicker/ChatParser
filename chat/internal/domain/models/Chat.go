package models

import "time"

type Message struct {
	Id      int
	Name    string
	Created time.Time
}
