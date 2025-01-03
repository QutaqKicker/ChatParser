package exporter

import (
	"backups/internal/exporter/writers"
	"context"
	"errors"
	backupv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/backup"
	chatv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/chat"
	commonv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/common"
	"sync"
)

var messagesBatchSize = int32(100000)

type fileWriter interface {
	WriteFile(ctx context.Context, writeDir string, messages []*chatv1.ChatMessage) error
}

type Exporter struct {
	writer     fileWriter
	chatClient chatv1.ChatClient
}

func NewExporter(exportType backupv1.ExportType) Exporter {
	var writer fileWriter
	switch exportType {
	case backupv1.ExportType_CSV:
		writer = writers.CsvWriter{}
	case backupv1.ExportType_PARQUET:
		writer = writers.ParquetWriter{}
	}

	return Exporter{writer: writer}
}

func (e *Exporter) ExportToDir(ctx context.Context, messageFilter *commonv1.MessagesFilter) error {
	writersWg := sync.WaitGroup{}

	//TODO Вынести адрес экспорта в конфиг
	exportDir := "C:\\Projects\\TestingData\\Export"

	taken := int32(0)
	for {
		messagesResponse, err := e.chatClient.GetMessages(ctx, &chatv1.SearchMessagesRequest{
			Filter: messageFilter,
			Skip:   taken,
			Take:   messagesBatchSize,
		})

		if err != nil {
			return err
		}
		if len(messagesResponse.Messages) == 0 {
			break
		}

		writersWg.Add(1)
		errorMutex := sync.Mutex{}
		go func() {
			writeErr := e.writer.WriteFile(ctx, exportDir, messagesResponse.Messages)
			if writeErr != nil {

				errorMutex.Lock()
				if err != nil {
					err = writeErr
				} else {
					err = errors.Join(err, writeErr)
				}
				errorMutex.Unlock()

			}
			writersWg.Done()
		}()
		taken += messagesBatchSize
	}

	writersWg.Wait()
	return nil
}
