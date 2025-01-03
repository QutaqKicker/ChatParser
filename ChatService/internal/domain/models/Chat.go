package models

import "time"

type Chat struct {
	Id      int32 `auto-generated:"true"`
	Name    string
	Created time.Time
}

func (c Chat) TableName() string {
	return "ChatService"
}

// FieldValuesAsArray / GetKeyByName fields values as any array. Needs for safe inserting
func (c Chat) FieldValuesAsArray() []any {
	return []any{c.Id, c.Name, c.Created}
}
