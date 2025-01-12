package contracts

import "time"

type ParquetMessageRow struct {
	Id               int32
	ChatId           int32
	ChatName         string
	UserId           string
	UserName         string
	ReplyToMessageId int32
	Text             string
	Created          time.Time
}
