package readers

import (
	"chat/internal/domain/models"
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"time"
)

type JsonChatContent struct {
	Name     string `json:"name"`
	ID       int    `json:"id"`
	Messages []struct {
		ID               int    `json:"id"`
		Type             string `json:"type"`
		DateUnixTime     string `json:"date_unixtime"`
		Text             string `json:"text"`
		From             string `json:"from,omitempty"`
		FromID           string `json:"from_id,omitempty"`
		ReplyToMessageID int    `json:"reply_to_message_id,omitempty"`
	} `json:"messages"`
}

type JsonReader struct {
	outMessagesChan chan<- models.Message
	errorsChan      chan<- error
}

func NewJsonReader(outMessagesChan chan<- models.Message, errorsChan chan<- error) *JsonReader {
	return &JsonReader{outMessagesChan, errorsChan}
}

func (r *JsonReader) ReaderType() models.DumpType {
	return models.Json
}

func (r *JsonReader) ReadMessages(ctx context.Context, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		r.errorsChan <- err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			r.errorsChan <- err
		}
	}(file)

	var content JsonChatContent

	byteArr, _ := io.ReadAll(file)

	err = json.Unmarshal(byteArr, &content)
	if err != nil {
		r.errorsChan <- err
	}

	for _, rawMessage := range content.Messages {
		select {
		case <-ctx.Done():
			r.errorsChan <- err
			return
		default:
			var message models.Message

			message.Id = int32(rawMessage.ID)

			dateUnix, err := strconv.ParseInt(rawMessage.DateUnixTime, 10, 64)
			if err != nil {
				r.errorsChan <- err
				return
			}

			message.Created = time.Unix(dateUnix, 0) //TODO maybe need a timezone

			message.UserId = rawMessage.FromID //TODO get user id from context with rawMessage.fromName

			message.ReplyToMessageId = int32(rawMessage.ReplyToMessageID)

			message.Text = rawMessage.Text
			r.outMessagesChan <- message
		}
	}
}
