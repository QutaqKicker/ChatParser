package parser

import (
	"chat/internal/domain/models"
	"chat/internal/domain/queryBuilders"
	"chat/internal/parser/readers"
	"context"
	"database/sql"
	"errors"
	"log"
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
	db     *sql.DB
}

func New(logger *slog.Logger, db *sql.DB) *Parser {
	return &Parser{
		logger,
		db,
	}
}

type Reader interface {
	ReaderType() models.DumpType
	ReadMessages(ctx context.Context, fileName string)
}

func (p *Parser) ParseFromDir(ctx context.Context, dumpDir string) error {
	var reader Reader
	messagesChan := make(chan models.Message, 30)
	errorsChan := make(chan error, 30)
	readersWg := &sync.WaitGroup{}
	dumpType, err := GetDumpType(dumpDir)
	if err != nil {
		log.Fatal(err)
	}
	switch dumpType {
	case models.Html:
		reader = readers.NewHtmlReader(messagesChan, errorsChan)
	case models.Json:
		reader = readers.NewJsonReader(messagesChan, errorsChan)
	}

	files, err := os.ReadDir(dumpDir)

	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), string(reader.ReaderType())) {
			readersWg.Add(1)
			go func() {
				reader.ReadMessages(ctx, dumpDir+"/"+file.Name())
				readersWg.Done()
			}()
		}
	}

	insertsWg := &sync.WaitGroup{}
	insertsWg.Add(1)
	go func() {
		InsertMessages(ctx, p.db, messagesChan)
		insertsWg.Done()
	}()

	errorsWg := &sync.WaitGroup{}
	errorsWg.Add(1)
	go func() {
		ProcessErrors(p.logger, errorsChan)
		errorsWg.Done()
	}()

	readersWg.Wait()
	close(messagesChan)
	close(errorsChan)
	insertsWg.Wait()
	errorsWg.Wait()

	return nil
}

func InsertMessages(ctx context.Context, db *sql.DB, messagesChan <-chan models.Message) {
	insertQuery := queryBuilders.BuildInsert[models.Message](false)
	for message := range messagesChan {
		_, err := db.ExecContext(ctx, insertQuery, message.FieldValuesAsArray()...)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ProcessErrors(log *slog.Logger, errorsChan <-chan error) {
	for err := range errorsChan {
		log.Error(err.Error())
	}
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
