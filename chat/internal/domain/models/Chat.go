package models

import "time"

type Chat struct {
	Id      int
	Name    string
	Created time.Time
}

func (c Chat) TableName() string {
	return "chats"
}

// FieldValuesAsArray / Get fields values as any array. Needs for safe inserting
func (c Chat) FieldValuesAsArray() []any {
	return []any{c.Id, c.Name, c.Created}
}
