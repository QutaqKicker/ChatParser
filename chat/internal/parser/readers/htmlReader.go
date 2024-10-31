package readers

import (
	"chat/internal/domain/models"
	"context"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type HtmlReader struct {
	wg              *sync.WaitGroup
	outMessagesChan chan<- models.Message
}

func NewHtmlReader(wg *sync.WaitGroup, outMessagesChan chan<- models.Message) *HtmlReader {
	return &HtmlReader{wg, outMessagesChan}
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

	historyNode, err := searchNode(bodyNode, NodeData, "history")
	if err != nil {
		log.Fatal(err)
	}

	var lastSenderId string
	numbersRegexp := regexp.MustCompile("[0-9]+")

	for messageNode := historyNode.FirstChild; messageNode != nil; messageNode = messageNode.NextSibling {
		//TODO Get message id
		messageBodyNode, err := searchNode(messageNode, NodeClass, "body")
		if err == nil {
			var message models.Message
			message.ChatId = chatId

			messageBodyChild := messageBodyNode.FirstChild
			createdDateValue, err := getAttributeValueByName(messageBodyChild, "title")
			if err != nil {
				log.Fatal(err)
			}
			message.Created, _ = time.Parse(createdDateValue, "DD.MM.YYYY HH:MM:SS")
			messageBodyChild = messageBodyChild.NextSibling

			nextClassName, err := getAttributeValueByName(messageBodyChild, "class")
			if err != nil {
				log.Fatal(err)
			}

			if nextClassName == "from_name" {
				userId := getUserIdFromName(messageBodyChild.Data) //TODO Check
				lastSenderId = userId
				messageBodyChild = messageBodyChild.NextSibling
				nextClassName, err = getAttributeValueByName(messageBodyChild, "class")
			}

			if nextClassName == "reply_to details" {
				hrefValue, err := getAttributeValueByName(messageBodyChild.LastChild, "href")
				if err != nil {
					log.Fatal(err)
				}
				repliedMessageId, err := strconv.Atoi(numbersRegexp.FindString(hrefValue))

				message.ReplyToMessageId = int32(repliedMessageId)
				messageBodyChild = messageBodyChild.NextSibling
			}

			if nextClassName == "text" {
				message.UserId = lastSenderId
				message.Text = messageBodyChild.Data
			}
			r.outMessagesChan <- message
		}
	}
	r.wg.Done()
}

func getAttributeValueByName(node *html.Node, attrName string) (string, error) {
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			return attr.Val, nil
		}
	}
	return "", fmt.Errorf("attribute with name %s does not exists", attrName)
}

func getChatNameFromBodyNode(bodyNode *html.Node) (string, error) {
	pageHeaderNode, err := searchNode(bodyNode, NodeClass, "page_header")
	if err != nil {
		return "", err
	}

	return pageHeaderNode.FirstChild.FirstChild.FirstChild.Data, nil
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
