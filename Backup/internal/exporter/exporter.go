package exporter

import (
	"backups/internal/exporter/writers"
	"context"
	"errors"
	backupv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/backup"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	commonv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/common"
	"sync"
)

var messagesBatchSize = int32(100000)

type fileWriter interface {
	WriteFile(ctx context.Context, writeDir string, messages []*chatv1.ChatMessage) error
}

type Exporter struct {
	chatClient chatv1.ChatClient
	exportDir  string
}

func NewExporter(chatClient chatv1.ChatClient, exportDir string) Exporter {
	return Exporter{chatClient: chatClient, exportDir: exportDir}
}

func (e Exporter) ExportToDir(ctx context.Context, exportType backupv1.ExportType, messageFilter *commonv1.MessagesFilter) error {
	var writer fileWriter
	switch exportType {
	case backupv1.ExportType_CSV:
		writer = writers.CsvWriter{}
	case backupv1.ExportType_PARQUET:
		writer = writers.ParquetWriter{}
	}

	var err error

	writersWg := sync.WaitGroup{}
	taken := int32(0)
	for {
		messagesResponse, rpcErr := e.chatClient.GetMessages(ctx, &chatv1.SearchMessagesRequest{
			Filter: messageFilter,
			Skip:   taken,
			Take:   messagesBatchSize,
		})

		if rpcErr != nil {
			return rpcErr
		}
		if len(messagesResponse.Messages) == 0 {
			break
		}

		writersWg.Add(1)
		errorMutex := sync.Mutex{}
		go func() {
			writeErr := writer.WriteFile(ctx, e.exportDir, messagesResponse.Messages)
			if writeErr != nil {

				errorMutex.Lock()
				if err != nil {
					err = errors.Join(err, writeErr)
				} else {
					err = writeErr
				}
				errorMutex.Unlock()

			}
			writersWg.Done()
		}()
		taken += messagesBatchSize
	}

	writersWg.Wait()
	return err
}
