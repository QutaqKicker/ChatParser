package models

import (
	"github.com/google/uuid"
	"time"
)

type Log struct {
	Id          uuid.UUID `auto-generated:"true"`
	ServiceName string    `column:"service_name"`
	Type        int
	Message     string
	Created     time.Time
}

func (u Log) TableName() string {
	return "logs"
}

// FieldValuesAsArray / GetKeyByName fields values as any array. Needs for safe inserting
func (u Log) FieldValuesAsArray() []any {
	return []any{u.Id, u.ServiceName, u.Type, u.Message, u.Created}
}
