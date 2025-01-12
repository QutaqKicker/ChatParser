package models

import "time"

type User struct {
	Name          string
	MessagesCount int `column:"messages_count"`
	Created       time.Time
}

func (u User) TableName() string {
	return "User"
}

// FieldValuesAsArray / GetKeyByName fields values as any array. Needs for safe inserting
func (u User) FieldValuesAsArray() []any {
	return []any{u.Name, u.MessagesCount, u.Created}
}
