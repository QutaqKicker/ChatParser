package parser

import (
	"chat/internal/caches"
	"chat/internal/domain/models"
	"chat/internal/parser/readers"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/dbHelper"
	"github.com/QutaqKicker/ChatParser/Common/myKafka"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Parser осуществляет чтение дампа сообщений, их дозаполнение, подсчет и отправку в микросервис Users и сохранение данных в БД
type Parser struct {
	logger                     *slog.Logger
	db                         *sql.DB
	userMessageCounterProducer *myKafka.UserMessageCounterProducer
}

func New(logger *slog.Logger, db *sql.DB, userMessageCounterProducer *myKafka.UserMessageCounterProducer) *Parser {
	return &Parser{
		logger,
		db,
		userMessageCounterProducer,
	}
}

// Reader интерфейс для читателей файлов дампа
type Reader interface {
	ReaderType() models.DumpType
	ReadMessages(ctx context.Context, fileName string)
}

// ParseFromDir парсит все файлы в директории. Выбирает формат дампа на основании формата первого подходящего для парса файлов.
// Файлы, соответствующие выбранному формату игнорируются в данном парсе. Наверное, стоит это поправить в дальнейшем
func (p *Parser) ParseFromDir(ctx context.Context, dumpDir string) error {
	rawMessagesChan := make(chan models.Message, 30)
	messagesChan := make(chan models.Message, 30)
	errorsChan := make(chan error, 30)

	dumpType, err := getDumpType(dumpDir)
	if err != nil {
		return err
	}

	var reader Reader
	switch dumpType {
	case models.Html:
		reader = readers.NewHtmlReader(rawMessagesChan, errorsChan)
	case models.Json:
		reader = readers.NewJsonReader(rawMessagesChan, errorsChan)
	case models.Parquet:
		reader = readers.NewParquetReader(rawMessagesChan, errorsChan)
	case models.Csv:
		reader = readers.NewCsvReader(rawMessagesChan, errorsChan)
	}

	files, err := os.ReadDir(dumpDir)
	if err != nil {
		return err
	}

	readersWg := &sync.WaitGroup{}
	readersSemaphore := make(chan struct{}, 5)
	//Запускаем чтение сообщений из файлов в директориях. Больше пяти горутин может помешать производительности
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), string(reader.ReaderType())) {
			readersWg.Add(1)
			go func() {
				readersSemaphore <- struct{}{}
				reader.ReadMessages(ctx, dumpDir+"/"+file.Name())
				<-readersSemaphore
				readersWg.Done()
			}()
		}
	}

	tx, err := p.db.BeginTx(ctx, nil)

	writersWg := &sync.WaitGroup{}
	//Заполняем сырые сообщения. У них могут быть не заполнены айдишники чата и юзера
	writersWg.Add(1)
	go func() {
		processRawMessages(ctx, tx, rawMessagesChan, messagesChan, p.userMessageCounterProducer)
		writersWg.Done()
	}()

	writersWg.Add(1)
	//Заполненные сообщения инсертим в БД
	go func() {
		insertMessages(ctx, tx, messagesChan)
		writersWg.Done()
	}()

	writersWg.Add(1)
	//Некритичные ошибки в ходе парса отдаем в логгер
	go func() {
		processErrors(p.logger, errorsChan)
		writersWg.Done()
	}()

	readersWg.Wait()
	close(rawMessagesChan)
	close(errorsChan)
	writersWg.Wait()

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// processRawMessages заполняет айди чата и юзера в незаполненных сообщениях и ведет подсчет для отправки в микросервис Users
func processRawMessages(ctx context.Context, tx *sql.Tx, inRawMessagesChan <-chan models.Message, outMessagesChan chan<- models.Message, userMessageCounterProducer *myKafka.UserMessageCounterProducer) {
	messagesCountPerUser := make(map[string]*atomic.Uint64)
chanLoop:
	for {
		select {
		case <-ctx.Done():
			return
		case rawMessage, ok := <-inRawMessagesChan:
			if !ok {
				break chanLoop
			}

			if rawMessage.ChatId == 0 {
				if chatId, ok := caches.ChatsCache.GetByName(tx, rawMessage.ChatName); ok {
					rawMessage.ChatId = chatId
				} else {
					rawMessage.ChatId = caches.ChatsCache.Set(tx, rawMessage.ChatName, rawMessage.ChatId)
				}
			} else {
				//Если такого чата в кэше нет или айди этого чата не совпадает с тем, что был в сообщении, апсертим этот чат
				if chatId, ok := caches.ChatsCache.GetByName(tx, rawMessage.ChatName); !ok || rawMessage.ChatId != chatId {
					caches.ChatsCache.Set(tx, rawMessage.ChatName, rawMessage.ChatId)
				}
			}

			if rawMessage.UserId == "" {
				if userId, ok := caches.UsersCache.GetByName(tx, rawMessage.UserName); ok {
					rawMessage.UserId = userId
				} else {
					rawMessage.UserId = caches.UsersCache.Set(tx, rawMessage.UserName, rawMessage.UserId)
				}
			} else {
				//Если такого юзера в кэше нет или айди этого юзера не совпадает с тем, что был в сообщении, апсертим этого юзера
				if userId, ok := caches.UsersCache.GetByName(tx, rawMessage.UserName); !ok || rawMessage.UserId != userId {
					caches.UsersCache.Set(tx, rawMessage.UserName, rawMessage.UserId)
				}
			}

			outMessagesChan <- rawMessage

			if counter, ok := messagesCountPerUser[rawMessage.UserName]; !ok {
				counter = &atomic.Uint64{}
				counter.Add(1)
				messagesCountPerUser[rawMessage.UserName] = counter
			} else {
				counter.Add(1)
			}
		}
	}
	close(outMessagesChan)

	for key, value := range messagesCountPerUser {
		err := userMessageCounterProducer.Send(myKafka.UserMessageCountRequest{
			UserName:     key,
			MessageCount: int32(value.Load()),
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

// insertMessages Инсертит в БД прочитанные заполненные сообщения
func insertMessages(ctx context.Context, tx *sql.Tx, messagesChan <-chan models.Message) {
	insertQuery, err := tx.Prepare(dbHelper.BuildInsert[models.Message](false))
	if err != nil {
		log.Fatal(err)
	}
	defer insertQuery.Close()

	messageCount := atomic.Uint64{}
	var startTime time.Time
	once := sync.Once{}
	for message := range messagesChan {
		once.Do(func() { startTime = time.Now() })

		_, err := insertQuery.ExecContext(ctx, message.FieldValuesAsArray()...)
		if err != nil {
			log.Fatal(err)
		}

		messageCount.Add(1)
		fmt.Printf("\rinserted %d messages. speed: %f messageActions per second", messageCount.Load(), float64(messageCount.Load())/time.Since(startTime).Seconds())
	}
}

func processErrors(log *slog.Logger, errorsChan <-chan error) {
	for err := range errorsChan {
		log.Error(err.Error())
	}
}

func getDumpType(dumpDir string) (models.DumpType, error) {
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
