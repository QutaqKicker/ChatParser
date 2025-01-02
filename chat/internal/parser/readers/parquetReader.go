package readers

import (
	"chat/internal/domain/models"
	"context"
	"github.com/QutaqKicker/ChatParser/common/contracts"
	"github.com/parquet-go/parquet-go"
)

type ParquetReader struct {
	outMessagesChan chan<- models.Message
	errorsChan      chan<- error
}

func NewParquetReader(outMessagesChan chan<- models.Message, errorsChan chan<- error) *ParquetReader {
	return &ParquetReader{outMessagesChan, errorsChan}
}

func (r *ParquetReader) ReaderType() models.DumpType {
	return models.Parquet
}

func (r *ParquetReader) ReadMessages(ctx context.Context, fileName string) {
	rows, err := parquet.ReadFile[contracts.ParquetMessageRow](fileName)
	if err != nil {
		r.errorsChan <- err
	}

	for _, row := range rows {
		select {
		case <-ctx.Done():
			r.errorsChan <- ctx.Err()
			return

		default:
			message := models.Message{
				Id:               row.Id,
				ChatId:           row.ChatId,
				ChatName:         row.ChatName,
				UserId:           row.UserId,
				UserName:         row.UserName,
				ReplyToMessageId: row.ReplyToMessageId,
				Text:             row.Text,
				Created:          row.Created,
			}
		}
	}
}
