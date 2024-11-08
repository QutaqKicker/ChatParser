package readers

import (
	"chat/internal/domain/models"
	"context"
	"encoding/json"
	"golang.org/x/net/html"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
	wg              *sync.WaitGroup
	outMessagesChan chan<- models.Message
}

func NewJsonReader(wg *sync.WaitGroup, outMessagesChan chan<- models.Message) *HtmlReader {
	return &HtmlReader{wg, outMessagesChan}
}

func (r *JsonReader) ReaderType() models.DumpType {
	return models.Json
}

func (r *JsonReader) ReadMessages(ctx context.Context, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	var content JsonChatContent

	byteArr, _ := io.ReadAll(file)

	err = json.Unmarshal(byteArr, &content)

	doc, err := html.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	bodyNode, err := searchNode(doc, NodeData, "body")
	if err != nil {
		log.Fatal(err)
	}

	chatName, err := getChatNameFromBodyNode(bodyNode)
	if err != nil {
		log.Fatal(err)
	}

	//TODO getChatIdFromName
	var chatId = getChatIdFromName(chatName)

	historyNode, err := searchNode(bodyNode, NodeClass, "history")
	if err != nil {
		log.Fatal(err)
	}

	var lastSenderId string
	numbersRegexp := regexp.MustCompile("[0-9]+")

	for messageNode := historyNode.FirstChild; messageNode != nil; messageNode = messageNode.NextSibling {
		messageId, _ := getAttributeValueByName(messageNode, "id")
		messageBodyNode, err := searchNode(messageNode, NodeClass, "body")
		if err == nil {
			var message models.Message
			message.ChatId = chatId

			Id, _ := strconv.Atoi(strings.Replace(messageId, "message", "", 1))
			message.Id = int32(Id)

			messageBodyChild := getElementNodeChild(messageBodyNode)
			createdDateValue, err := getAttributeValueByName(messageBodyChild, "title")
			if err != nil {
				log.Fatal(err)
			}
			createdDateValue = strings.Replace(createdDateValue, "UTC", "MSK", 1)

			message.Created, _ = time.Parse("02.01.2006 15:04:05 MST-07:00", createdDateValue) //TODO вынести формат в константы
			messageBodyChild = getElementNodeSibling(messageBodyChild)

			nextClassName, err := getAttributeValueByName(messageBodyChild, "class")
			if err != nil {
				log.Fatal(err)
			}

			if nextClassName == "from_name" {
				userId := strings.TrimSpace(messageBodyChild.FirstChild.Data) //getUserIdFromName(messageBodyChild.Data) //TODO Check
				lastSenderId = userId
				messageBodyChild = getElementNodeSibling(messageBodyChild)
				nextClassName, err = getAttributeValueByName(messageBodyChild, "class")
			}
			message.UserId = lastSenderId

			if nextClassName == "media_wrap clearfix" {
				messageBodyChild = getElementNodeSibling(messageBodyChild)
				nextClassName, err = getAttributeValueByName(messageBodyChild, "class")
			}

			if nextClassName == "reply_to details" {
				hrefValue, err := getAttributeValueByName(getElementNodeChild(messageBodyChild), "href")
				if err != nil {
					log.Fatal(err)
				}
				repliedMessageId, err := strconv.Atoi(numbersRegexp.FindString(hrefValue))

				message.ReplyToMessageId = int32(repliedMessageId)
				messageBodyChild = getElementNodeSibling(messageBodyChild)
				nextClassName, err = getAttributeValueByName(messageBodyChild, "class")
			}

			if nextClassName == "text" {
				text := strings.Builder{}
				for textNode := messageBodyChild.FirstChild; textNode != nil; textNode = textNode.NextSibling {
					if textNode.Type == html.TextNode {
						text.WriteString(strings.TrimSpace(textNode.Data) + "\n")
					}
					if textNode.Type == html.ElementNode {
						if textNode.Data == "br" {
							text.WriteString("\n")
						} else if textNode.Data == "a" {
							text.WriteString(textNode.FirstChild.Data + " ")
						}
					}
				}
				message.Text = strings.TrimSpace(text.String())
			}
			r.outMessagesChan <- message
		}
	}
	r.wg.Done()
}
