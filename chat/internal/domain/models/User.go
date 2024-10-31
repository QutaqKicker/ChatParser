package models

import "time"

type User struct {
	Id      string
	Name    string
	Created time.Time
}
