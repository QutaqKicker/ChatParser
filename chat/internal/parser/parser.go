package parser

import (
	"chat/internal/domain/models"
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type DumpReader interface {
	ReadMessages(ctx context.Context, fileDir string, outChan <-chan models.Message, wg *sync.WaitGroup)
	ReaderType() models.DumpType
}

type Parser struct {
	logger *slog.Logger
	reader *DumpReader
}

func New(logger *slog.Logger, reader *DumpReader) *Parser {
	return &Parser{
		logger,
		reader,
	}
}

func (p *Parser) ParseDir(ctx context.Context, dumpDir string) (<-chan models.Message, error) {
	files, err := os.ReadDir(dumpDir)

	if err != nil {
		return nil, err
	}

	outChan := make(chan models.Message, 30)

	wg := &sync.WaitGroup{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), string((*p.reader).ReaderType())) {
			wg.Add(1)
			go (*p.reader).ReadMessages(ctx, dumpDir+"/"+file.Name(), outChan, wg)
		}
	}

	return outChan, nil
}

func GetDumpType(dumpDir string) (models.DumpType, error) {
	files, err := os.ReadDir(dumpDir)
	if err != nil {
		return models.Html, err
	}

	for _, value := range files {
		if value.IsDir() {
			continue
		}

		if strings.HasSuffix(value.Name(), models.Html) {
			return models.Html, nil
		}
		if strings.HasSuffix(value.Name(), models.Json) {
			return models.Json, nil
		}
		if strings.HasSuffix(value.Name(), models.Csv) {
			return models.Csv, nil
		}
		if strings.HasSuffix(value.Name(), models.Parquet) {
			return models.Parquet, nil
		}
	}

	return models.Html, errors.New("dumps on selected dir does not exists")
}
