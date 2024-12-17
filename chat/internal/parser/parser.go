package parser

import (
	"chat/internal/caches"
	"chat/internal/dbHelper"
	"chat/internal/domain/models"
	"chat/internal/parser/readers"
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

// Parser осуществляет чтение дампа сообщений, их дозаполнение, подсчет и отправку в микросервис Users и сохранение данных в БД
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
	} //TODO остальные форматы сюда бахнуть

	files, err := os.ReadDir(dumpDir)
	if err != nil {
		return err
	}

	readersWg := &sync.WaitGroup{}
	//Запускаем чтение сообщений из файлов в директориях
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

	tx, err := p.db.BeginTx(ctx, nil)

	writersWg := &sync.WaitGroup{}
	writersWg.Add(1)
	//Заполняем сырые сообщения. У них могут быть не заполнены айдишники чата и юзера
	go func() {
		processRawMessages(ctx, tx, rawMessagesChan, messagesChan)
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

// UsersWithMessagesCount Сущность для отправки информации по количеству сообщений в микросервис Users
type UsersWithMessagesCount struct {
	UserName      string
	MessagesCount uint64
}

// processRawMessages заполняет айди чата и юзера в незаполненных сообщениях и ведет подсчет для отправки в микросервис Users
func processRawMessages(ctx context.Context, tx *sql.Tx, inRawMessagesChan <-chan models.Message, outMessagesChan chan<- models.Message) {
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
				if chatId, ok := caches.ChatsCache.Get(tx, rawMessage.ChatName); ok {
					rawMessage.ChatId = chatId
				} else {
					rawMessage.ChatId = caches.ChatsCache.Set(tx, rawMessage.ChatName, rawMessage.ChatId)
				}
			} else {
				//Если такого чата в кэше нет или айди этого чата не совпадает с тем, что был в сообщении, апсертим этот чат
				if chatId, ok := caches.ChatsCache.Get(tx, rawMessage.ChatName); !ok || rawMessage.ChatId != chatId {
					caches.ChatsCache.Set(tx, rawMessage.ChatName, rawMessage.ChatId)
				}
			}

			if rawMessage.UserId == "" {
				if userId, ok := caches.UsersCache.Get(tx, rawMessage.UserName); ok {
					rawMessage.UserId = userId
				} else {
					rawMessage.UserId = caches.UsersCache.Set(tx, rawMessage.UserName, rawMessage.UserId)
				}
			} else {
				//Если такого юзера в кэше нет или айди этого юзера не совпадает с тем, что был в сообщении, апсертим этого юзера
				if userId, ok := caches.UsersCache.Get(tx, rawMessage.UserName); !ok || rawMessage.UserId != userId {
					caches.UsersCache.Set(tx, rawMessage.UserName, rawMessage.UserId)
				}
			}

			outMessagesChan <- rawMessage

			if counter, ok := messagesCountPerUser[rawMessage.UserName]; !ok {
				messagesCountPerUser[rawMessage.UserName] = &atomic.Uint64{}
			} else {
				counter.Add(1)
			}
		}
	}
	close(outMessagesChan)

	for key, value := range messagesCountPerUser {
		_ = UsersWithMessagesCount{key, value.Load()}
		//TODO Организовать отправку этого добра в сервис Users и вынести в основной метод. Если этот метод будет у нескольких горутин, данные отправим несколько раз вместо одного
	}
}

// insertMessages Инсертит в БД прочитанные заполненные сообщения
func insertMessages(ctx context.Context, tx *sql.Tx, messagesChan <-chan models.Message) {
	insertQuery, err := tx.Prepare(dbHelper.BuildInsert[models.Message](false))
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
