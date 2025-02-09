package readers

import (
	"chat/internal/domain/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type TgText struct {
	Content string
}

type TgArrayInText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (t *TgText) UnmarshalJSON(data []byte) error {
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		t.Content = stringValue
		return nil
	}

	var arrayValue []interface{}
	if err := json.Unmarshal(data, &arrayValue); err == nil {
		sb := strings.Builder{}
		for _, value := range arrayValue {
			switch tv := value.(type) {
			case string:
				sb.WriteString(tv)
			case map[string]interface{}:
				sb.WriteString(tv["text"].(string))
			}
		}
		t.Content = sb.String()
		return nil
	}

	return fmt.Errorf("invalid value for Text")
}

type JsonChatContent struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	ID       int    `json:"id"`
	Messages []struct {
		ID             int           `json:"id"`
		Type           string        `json:"type"`
		Date           string        `json:"date"`
		DateUnixtime   string        `json:"date_unixtime"`
		Actor          string        `json:"actor,omitempty"`
		ActorID        string        `json:"actor_id,omitempty"`
		Action         string        `json:"action,omitempty"`
		Inviter        string        `json:"inviter,omitempty"`
		Text           TgText        `json:"text"`
		TextEntities   []interface{} `json:"text_entities"`
		From           string        `json:"from,omitempty"`
		FromID         string        `json:"from_id,omitempty"`
		Edited         string        `json:"edited,omitempty"`
		EditedUnixtime string        `json:"edited_unixtime,omitempty"`
		Reactions      []struct {
			Type   string `json:"type"`
			Count  int    `json:"count"`
			Emoji  string `json:"emoji"`
			Recent []struct {
				From   string `json:"from"`
				FromID string `json:"from_id"`
				Date   string `json:"date"`
			} `json:"recent"`
		} `json:"reactions,omitempty"`
		ReplyToMessageID int      `json:"reply_to_message_id,omitempty"`
		ForwardedFrom    string   `json:"forwarded_from,omitempty"`
		Photo            string   `json:"photo,omitempty"`
		Width            int      `json:"width,omitempty"`
		Height           int      `json:"height,omitempty"`
		File             string   `json:"file,omitempty"`
		MediaType        string   `json:"media_type,omitempty"`
		MimeType         string   `json:"mime_type,omitempty"`
		DurationSeconds  int      `json:"duration_seconds,omitempty"`
		FileName         string   `json:"file_name,omitempty"`
		Thumbnail        string   `json:"thumbnail,omitempty"`
		StickerEmoji     string   `json:"sticker_emoji,omitempty"`
		ViaBot           string   `json:"via_bot,omitempty"`
		Members          []string `json:"members,omitempty"`
		MediaSpoiler     bool     `json:"media_spoiler,omitempty"`
		Performer        string   `json:"performer,omitempty"`
		Title            string   `json:"title,omitempty"`
		Poll             struct {
			Question    string `json:"question"`
			Closed      bool   `json:"closed"`
			TotalVoters int    `json:"total_voters"`
			Answers     []struct {
				Text   string `json:"text"`
				Voters int    `json:"voters"`
				Chosen bool   `json:"chosen"`
			} `json:"answers"`
		} `json:"poll,omitempty"`
		LocationInformation struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location_information,omitempty"`
		MessageID        int    `json:"message_id,omitempty"`
		ReplyToPeerID    string `json:"reply_to_peer_id,omitempty"`
		NewTitle         string `json:"new_title,omitempty"`
		NewIconEmojiID   int    `json:"new_icon_emoji_id,omitempty"`
		InlineBotButtons [][]struct {
			Type string `json:"type"`
			Text string `json:"text"`
			Data string `json:"data"`
		} `json:"inline_bot_buttons,omitempty"`
	} `json:"messageActions"`
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
			if rawMessage.FromID == "" {
				continue
			}

			var message models.Message

			message.Id = int32(rawMessage.ID)

			dateUnix, err := strconv.ParseInt(rawMessage.DateUnixtime, 10, 64)
			if err != nil {
				r.errorsChan <- err
				return
			}

			message.Created = time.Unix(dateUnix, 0) //TODO maybe need a timezone

			message.ChatId = int32(content.ID)
			message.ChatName = content.Name

			message.UserId = rawMessage.FromID
			message.UserName = rawMessage.From

			message.ReplyToMessageId = int32(rawMessage.ReplyToMessageID)

			message.Text = rawMessage.Text.Content
			r.outMessagesChan <- message
		}
	}
}
