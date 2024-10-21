package models

import "time"

type User struct {
	Id      int
	Name    string
	Created time.Time
}
