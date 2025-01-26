package writers

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/contracts"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	"os"
	"strconv"
	"time"
)

type CsvWriter struct{}

func (CsvWriter) WriteFile(ctx context.Context, writeDir string, messages []*chatv1.ChatMessage) (err error) {
	if len(messages) == 0 {
		return nil
	}

	fileName := messages[len(messages)-1].Created.AsTime().Format(time.DateOnly) + ".csv" //TODO Проверить
	file, err := os.Create(fmt.Sprintf("%s\\%s", writeDir, fileName))
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()

	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = ';'
	csvWriter.UseCRLF = true
	defer csvWriter.Flush()

	err = csvWriter.Write(contracts.CsvHeaderColumns)
	if err != nil {
		return err
	}

	for i := 0; i < len(messages); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err = csvWriter.Write(FieldValuesAsArray(messages[i]))
			if err != nil {
				return fmt.Errorf("error on writing message with id %d. error: %w", messages[i].Id, err)
			}
		}
	}

	return nil
}

func FieldValuesAsArray(m *chatv1.ChatMessage) []string {
	return []string{strconv.Itoa(int(m.Id)),
		strconv.Itoa(int(m.ChatId)),
		m.ChatName,
		m.UserId,
		m.UserName,
		strconv.Itoa(int(m.ReplyToMessageId)),
		m.Text,
		m.Created.AsTime().String()}
}
