package parser

import (
	"chat/internal/db"
	"chat/internal/domain/models"
	"chat/internal/domain/queryFilters"
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
	rawMessagesChan := make(chan models.Message, 30)
	messagesChan := make(chan models.Message, 30)
	errorsChan := make(chan error, 30)
	readersWg := &sync.WaitGroup{}
	dumpType, err := GetDumpType(dumpDir)
	if err != nil {
		log.Fatal(err)
	}
	switch dumpType {
	case models.Html:
		reader = readers.NewHtmlReader(rawMessagesChan, errorsChan)
	case models.Json:
		reader = readers.NewJsonReader(rawMessagesChan, errorsChan)
	}

	files, err := os.ReadDir(dumpDir)

	if err != nil {
		return err
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
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
		ProcessRawMessages(ctx, tx, rawMessagesChan, messagesChan)
		insertsWg.Done()
	}()

	insertsWg.Add(1)
	go func() {
		InsertMessages(ctx, tx, messagesChan)
		insertsWg.Done()
	}()

	errorsWg := &sync.WaitGroup{}
	errorsWg.Add(1)
	go func() {
		ProcessErrors(p.logger, errorsChan)
		errorsWg.Done()
	}()

	readersWg.Wait()
	close(rawMessagesChan)
	close(errorsChan)
	insertsWg.Wait()
	errorsWg.Wait()

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func InsertMessages(ctx context.Context, tx *sql.Tx, messagesChan <-chan models.Message) {
	insertQuery, err := tx.Prepare(db.BuildInsert[models.Message](false))
	if err != nil {
		log.Fatal(err)
	}
	defer insertQuery.Close()

	for message := range messagesChan {
		_, err := insertQuery.ExecContext(ctx, message.FieldValuesAsArray()...)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type cacheForParse[T comparable] struct {
	ma map[string]T
}

var seenUsers = make(map[string]string)
var seenChats = make(map[string]string)

func ProcessRawMessages(ctx context.Context, tx *sql.Tx, inRawMessagesChan <-chan models.Message, outMessagesChan chan<- models.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case rawMessage := <-inRawMessagesChan:
			if rawMessage.ChatId == 0 {
				chatQuery := db.BuildQuery(db.QueryBuildRequest[queryFilters.ChatFilter]{Filter: queryFilters.ChatFilter{Name: rawMessage.ChatName}})

			} else { //На случай, если пришли четко определенные

			}

			if rawMessage.UserId == "" {

			}

			outMessagesChan <- rawMessage
		}
	}
	close(outMessagesChan)
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
