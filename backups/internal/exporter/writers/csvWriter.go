package writers

import (
	"context"
	"encoding/csv"
	"fmt"
	chatv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/chat"
	"os"
	"strconv"
)

var csvHeaderColumns = []string{"Id",
	"ChatId",
	"ChatName",
	"UserId",
	"UserName",
	"ReplyToMessageId",
	"Text",
	"Created"}

type CsvWriter struct{}

func (CsvWriter) WriteFile(ctx context.Context, writeDir string, messages []chatv1.ChatMessage) error {
	if len(messages) == 0 {
		return nil
	}

	fileName := messages[len(messages)-1].Created.String() + ".csv" //TODO Проверить
	file, err := os.Create(fmt.Sprintf("%s/%s", writeDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = ';'
	csvWriter.UseCRLF = true
	defer csvWriter.Flush()

	err = csvWriter.Write(csvHeaderColumns)
	if err != nil {
		return err
	}

	for i := 0; i < len(messages); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err = csvWriter.Write(FieldValuesAsArray(&messages[i]))
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
		m.Created.String()}
}
