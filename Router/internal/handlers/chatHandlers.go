package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	commonv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/common"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetMessagesCountHandler(logger *slog.Logger, chatClient *chatv1.ChatClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var chatFilter *commonv1.MessagesFilter

		err := json.NewDecoder(r.Body).Decode(chatFilter)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		response, err := (*chatClient).GetMessagesCount(ctx, &chatv1.GetMessagesCountRequest{Filter: chatFilter})

		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprintf(w, strconv.Itoa(int(response.Count)))
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}

func SearchMessagesHandler(logger *slog.Logger, chatClient *chatv1.ChatClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var messagesRequest *chatv1.SearchMessagesRequest

		err := json.NewDecoder(r.Body).Decode(messagesRequest)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		messages, err := (*chatClient).SearchMessages(ctx, messagesRequest)

		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func ParseFromDirHandler(logger *slog.Logger, chatClient *chatv1.ChatClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dirPath := r.URL.Query().Get("dir-path")

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		response, err := (*chatClient).ParseFromDir(ctx, &chatv1.ParseFromDirRequest{DirPath: dirPath})

		if err != nil || !response.Ok {
			if err != nil {
				logger.Error(err.Error())
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
}
