package models

import "time"

type User struct {
	Id      string `auto-generated:"true"`
	Name    string
	Created time.Time
}

func (u User) TableName() string {
	return "users"
}

// FieldValuesAsArray / GetKeyByName fields values as any array. Needs for safe inserting
func (u User) FieldValuesAsArray() []any {
	return []any{u.Id, u.Name, u.Created}
}
