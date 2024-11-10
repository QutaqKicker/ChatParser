package readers

import (
	"chat/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var numbersRegexp = regexp.MustCompile("[0-9]+")

type HtmlReader struct {
	outMessagesChan chan<- models.Message
	errorsChan      chan<- error
}

func NewHtmlReader(outMessagesChan chan<- models.Message, errorsChan chan<- error) *HtmlReader {
	return &HtmlReader{outMessagesChan, errorsChan}
}

func (r *HtmlReader) ReaderType() models.DumpType {
	return models.Html
}

func (r *HtmlReader) ReadMessages(ctx context.Context, fileName string) {
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

	processMessageNode := func(node *html.Node) func() {
		var lastSenderId string

		return func() {
			message, err := parseMessageNode(node)
			if err != nil {
				if message != nil && message.Id != 0 {
					r.errorsChan <- fmt.Errorf("could not parse message with id %d. error: %w", message.Id, err)
				} else {
					r.errorsChan <- fmt.Errorf("error on parse message of chat %s. error: %w", chatName, err)
				}
			} else {
				message.ChatId = chatId
				if message.UserId != "" {
					lastSenderId = message.UserId
				} else {
					message.UserId = lastSenderId
				}

				r.outMessagesChan <- *message
			}
		}
	}

	for messageNode := historyNode.FirstChild; messageNode != nil; messageNode = messageNode.NextSibling {
		select {
		case <-ctx.Done():
			r.errorsChan <- ctx.Err()

		default:
			processMessageNode(messageNode)
		}
	}
}

func parseMessageNode(node *html.Node) (*models.Message, error) {
	messageBodyNode, err := searchNode(node, NodeClass, "body")
	if err == nil {
		var message models.Message

		message.Id, err = getMessageId(node) //messageNode, not a messageBodyNode
		if err != nil {
			return &message, err
		}

		messageBodyChild := getElementNodeChild(messageBodyNode)

		message.Created, err = getMessageCreated(&messageBodyNode)
		if err != nil {
			return &message, err
		}

		nextClassName, err := getAttributeValueByName(messageBodyChild, "class")
		if err != nil {
			return &message, err
		}

		message.UserId, err = getMessageUserId(&messageBodyChild, &nextClassName)
		if err != nil {
			//Ничего не делаем, идем дальше. Id автора пуст если это не первое сообщение в серии сообщений от одного пользователя.
			//Возьмем автора с первого сообщения серии
		}

		err = SkipUnusedMessageBodyTags(&messageBodyChild, &nextClassName)
		if err != nil {
			return &message, err
		}

		message.ReplyToMessageId, err = getReplyToMessageId(&messageBodyChild, &nextClassName)
		if err != nil {
			return &message, err
		}

		message.Text, err = getMessageText(messageBodyChild, nextClassName)
		return &message, err
	}
	return nil, nil
}

func getMessageId(messageNode *html.Node) (int32, error) {
	messageIdAttrValue, _ := getAttributeValueByName(messageNode, "id")
	messageId, err := strconv.Atoi(strings.Replace(messageIdAttrValue, "message", "", 1))
	if err != nil {
		return 0, err
	}
	return int32(messageId), nil
}

func getMessageCreated(node **html.Node) (time.Time, error) {
	createdDateValue, err := getAttributeValueByName(*node, "title")
	if err != nil {
		return time.Time{}, err
	}
	createdDateValue = strings.Replace(createdDateValue, "UTC", "MSK", 1)

	*node = getElementNodeSibling(*node)
	return time.Parse("02.01.2006 15:04:05 MST-07:00", createdDateValue) //TODO вынести формат в константы
}

func getMessageUserId(node **html.Node, className *string) (string, error) {
	if *className == "from_name" {
		userId := strings.TrimSpace((*node).FirstChild.Data) //getUserIdFromName(messageBodyChild.Data) //TODO Check
		*node = getElementNodeSibling(*node)
		var err error
		*className, err = getAttributeValueByName(*node, "class")
		return userId, err
	}
	return "", errors.New("user id does not exists on this message, maybe that user id exists on previous message")
}

func SkipUnusedMessageBodyTags(node **html.Node, className *string) error {
	if *className == "media_wrap clearfix" {
		*node = getElementNodeSibling(*node)
		var err error
		*className, err = getAttributeValueByName(*node, "class")
		return err
	}
	return nil
}

func getReplyToMessageId(node **html.Node, className *string) (int32, error) {
	if *className == "reply_to details" {
		hrefValue, err := getAttributeValueByName(getElementNodeChild(*node), "href")
		if err != nil {
			return 0, err
		}

		repliedMessageId, err := strconv.Atoi(numbersRegexp.FindString(hrefValue))
		*node = getElementNodeSibling(*node)
		*className, err = getAttributeValueByName(*node, "class")

		return int32(repliedMessageId), err
	}
	return 0, nil
}

func getMessageText(node *html.Node, className string) (string, error) {
	if className == "text" {
		text := strings.Builder{}

		for textNode := node.FirstChild; textNode != nil; textNode = textNode.NextSibling {
			switch textNode.Type {
			case html.TextNode:
				text.WriteString(strings.TrimSpace(textNode.Data) + "\n")
			case html.ElementNode:
				if textNode.Data == "br" {
					text.WriteString("\n")
				} else if textNode.Data == "a" {
					text.WriteString(textNode.FirstChild.Data + " ")
				}
			}
		}

		return strings.TrimSpace(text.String()), nil
	}

	return "", errors.New("message text does not exists")
}

func getAttributeValueByName(node *html.Node, attrName string) (string, error) {
	if node != nil {
		for _, attr := range node.Attr {
			if attr.Key == attrName {
				return attr.Val, nil
			}
		}
	}
	return "", fmt.Errorf("attribute with name %s does not exists", attrName)
}

func getElementNodeSibling(node *html.Node) *html.Node {
	for currentNode := node.NextSibling; currentNode != nil; currentNode = currentNode.NextSibling {
		if currentNode.Type == html.ElementNode {
			return currentNode
		}
	}
	return nil
}

func getElementNodeChild(node *html.Node) *html.Node {
	return getElementNodeSibling(node.FirstChild)
}

func getChatNameFromBodyNode(bodyNode *html.Node) (string, error) {
	pageHeaderNode, err := searchNode(bodyNode, NodeClass, "page_header")
	if err != nil {
		return "", err
	}
	for pageHeaderNodeChild := pageHeaderNode.FirstChild; pageHeaderNodeChild != nil; pageHeaderNodeChild = pageHeaderNodeChild.NextSibling {
		if pageHeaderNodeChild.Type == html.ElementNode {
			if class, err := getAttributeValueByName(pageHeaderNodeChild, NodeClass); err == nil && class == "text bold" {
				return pageHeaderNodeChild.FirstChild.Data, nil
			}
			pageHeaderNodeChild = pageHeaderNodeChild.FirstChild
		}
	}
	return "", err
}

func searchNode(node *html.Node, searchType searchType, value string) (*html.Node, error) {
	for currNode := node.FirstChild; currNode != nil; currNode = currNode.NextSibling {
		switch searchType {
		case NodeData:
			if currNode.Data == value {
				return currNode, nil
			}
		case NodeClass:
			for _, attr := range currNode.Attr {
				if attr.Key == "class" && attr.Val == value {
					return currNode, nil
				}
			}
		}

		if currNode.Type == html.ElementNode {
			searchResult, _ := searchNode(currNode, searchType, value)
			if searchResult != nil {
				return searchResult, nil
			}
		}
	}
	return nil, fmt.Errorf("node with %s = %s does not exists", searchType, value)
}

type searchType string

const (
	NodeData  searchType = "data"
	NodeClass            = "class"
)
