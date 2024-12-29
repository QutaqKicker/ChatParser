package exporter

import (
	"backups/internal/exporter/writers"
	"context"
	backupv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/backup"
	chatv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/chat"
)

type fileWriter interface {
	WriteFile(ctx context.Context, writeDir string, messages []chatv1.ChatMessage) error
}

type Exporter struct {
	writer fileWriter
}

func NewExporter(exportType backupv1.ExportType) Exporter {
	var writer fileWriter
	switch exportType {
	case backupv1.ExportType_CSV:
		writer = writers.CsvWriter{}
		//TODO parquet
	}

	return Exporter{writer: writer}
}

func (e *Exporter) ExportToDir(ctx context.Context, exportDir string) error {

}
