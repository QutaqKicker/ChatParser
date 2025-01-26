package readers

import (
	"chat/internal/domain/models"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/contracts"
	"io"
	"os"
	"strconv"
	"time"
)

type CsvReader struct {
	outMessagesChan chan<- models.Message
	errorsChan      chan<- error
}

func NewCsvReader(outMessagesChan chan<- models.Message, errorsChan chan<- error) *CsvReader {
	return &CsvReader{outMessagesChan, errorsChan}
}

func (r *CsvReader) ReaderType() models.DumpType {
	return models.Csv
}

func (r *CsvReader) ReadMessages(ctx context.Context, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		r.errorsChan <- err
		return
	}

	reader := csv.NewReader(file)
	reader.Comma = ';'
	headerColumns, err := reader.Read()

	err = checkHeaderColumns(headerColumns)
	if err != nil {
		r.errorsChan <- err
		return
	}

	for {
		select {
		case <-ctx.Done():
			r.errorsChan <- ctx.Err()
			return

		default:
			row, err := reader.Read()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					r.errorsChan <- err
				}
				return
			}

			message := models.Message{}

			id, err := strconv.Atoi(row[0])
			if err != nil {
				r.errorsChan <- err
				return
			}

			chatId, err := strconv.Atoi(row[1])
			if err != nil {
				r.errorsChan <- err
				return
			}

			replyToMessageId, err := strconv.Atoi(row[5])
			if err != nil {
				r.errorsChan <- err
				return
			}

			created, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", row[7])
			if err != nil {
				r.errorsChan <- err
				return
			}

			message.Id = int32(id)
			message.ChatId = int32(chatId)
			message.ChatName = row[2]
			message.UserId = row[3]
			message.UserName = row[4]
			message.ReplyToMessageId = int32(replyToMessageId)
			message.Text = row[6]
			message.Created = created

			r.outMessagesChan <- message
		}
	}
}

func checkHeaderColumns(columns []string) error {
	if len(columns) != len(contracts.CsvHeaderColumns) {
		return fmt.Errorf("incorrect columns count. expected: %d", len(contracts.CsvHeaderColumns))
	} else {
		for i := 0; i < len(columns); i++ {
			if columns[i] != contracts.CsvHeaderColumns[i] {
				return fmt.Errorf("incorrect csv format. error at column with index %d", i)
			}
		}
	}

	return nil
}
