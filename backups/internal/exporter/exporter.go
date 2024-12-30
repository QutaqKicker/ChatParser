package exporter

import (
	"backups/internal/exporter/writers"
	"context"
	backupv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/backup"
	chatv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/chat"
	commonv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/common"
	"sync"
)

var messagesBatchSize = int32(100000)

type fileWriter interface {
	WriteFile(ctx context.Context, writeDir string, messages []chatv1.ChatMessage) error
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

func (e *Exporter) ExportToDir(ctx context.Context, exportDir string, messageFilter *commonv1.MessagesFilter) error {
	writersWg := sync.WaitGroup{}

	taken := 0
	for {
		messages, err := e.chatClient.GetMessages(ctx, &chatv1.SearchMessagesRequest{Filter: messageFilter,
			Skip: int32(taken),
			Take: messagesBatchSize})

		if err != nil {
			return err
		}
		if len(messages.Messages) == 0 {
			break
		}

		if


	}

}
