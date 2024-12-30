package writers

import (
	"context"
	chatv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/chat"
	"github.com/parquet-go/parquet-go"
	"os"
	"time"
)

type ParquetWriter struct{}

type parquetMessageRow struct {
	Id               int32
	ChatId           int32
	ChatName         string
	UserId           string
	UserName         string
	ReplyToMessageId int32
	Text             string
	Created          time.Time
}

func (ParquetWriter) WriteFile(_ context.Context, writeDir string, messages []*chatv1.ChatMessage) error {
	schema := parquet.SchemaOf(new(parquetMessageRow))
	fileName := messages[0].Created.String() + ".parquet" //TODO Проверить
	file, err := os.Create(writeDir + "/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := parquet.NewGenericWriter[parquetMessageRow](file, schema)
	defer writer.Close()

	rows := mapChatMessagesToParquetRows(messages)

	_, err = writer.Write(rows)
	if err != nil {
		return err
	}

	return nil
}

func mapChatMessagesToParquetRows(messages []*chatv1.ChatMessage) []parquetMessageRow {
	result := make([]parquetMessageRow, len(messages))
	for i := 0; i < len(messages); i++ {
		result[i] = parquetMessageRow{
			messages[i].Id,
			messages[i].ChatId,
			messages[i].ChatName,
			messages[i].UserId,
			messages[i].UserName,
			messages[i].ReplyToMessageId,
			messages[i].Text,
			messages[i].Created.AsTime(),
		}
	}

	return result
}
