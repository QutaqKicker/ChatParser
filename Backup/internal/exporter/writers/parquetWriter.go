package writers

import (
	"context"
	"github.com/QutaqKicker/ChatParser/Common/contracts"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	"github.com/parquet-go/parquet-go"
	"os"
)

type ParquetWriter struct{}

func (ParquetWriter) WriteFile(_ context.Context, writeDir string, messages []*chatv1.ChatMessage) error {
	schema := parquet.SchemaOf(new(contracts.ParquetMessageRow))
	fileName := messages[0].Created.String() + ".parquet" //TODO Проверить
	file, err := os.Create(writeDir + "/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := parquet.NewGenericWriter[contracts.ParquetMessageRow](file, schema)
	defer writer.Close()

	rows := mapChatMessagesToParquetRows(messages)

	_, err = writer.Write(rows)
	if err != nil {
		return err
	}

	return nil
}

func mapChatMessagesToParquetRows(messages []*chatv1.ChatMessage) []contracts.ParquetMessageRow {
	result := make([]contracts.ParquetMessageRow, len(messages))
	for i := 0; i < len(messages); i++ {
		result[i] = contracts.ParquetMessageRow{
			Id:               messages[i].Id,
			ChatId:           messages[i].ChatId,
			ChatName:         messages[i].ChatName,
			UserId:           messages[i].UserId,
			UserName:         messages[i].UserName,
			ReplyToMessageId: messages[i].ReplyToMessageId,
			Text:             messages[i].Text,
			Created:          messages[i].Created.AsTime(),
		}
	}

	return result
}
