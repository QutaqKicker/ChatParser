package services

import (
	"chat/internal/caches"
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"chat/internal/parser"
	"context"
	"database/sql"
	"github.com/QutaqKicker/ChatParser/Common/dbHelper"
	"github.com/QutaqKicker/ChatParser/Common/myKafka"
	"log/slog"
)

type ChatService struct {
	log *slog.Logger
	db  *sql.DB
	MessageSaver
	MessageProvider
	userMessageCounterProducer *myKafka.UserMessageCounterProducer
}

type HtmlParser interface {
	ParseFromDir(ctx context.Context, dumpDir string) (<-chan models.Message, error)
}

type MessageSaver interface {
	SaveMessages(
		ctx context.Context,
		message []models.Message) error
}

type MessageProvider interface {
	GetMessages(
		filter filters.MessageFilter) ([]models.Message, error)
}

func NewChatService(
	log *slog.Logger,
	db *sql.DB,
	userMessageCounterProducer *myKafka.UserMessageCounterProducer) *ChatService {
	return &ChatService{
		log:                        log,
		db:                         db,
		userMessageCounterProducer: userMessageCounterProducer,
	}
}

func (s *ChatService) Parse(ctx context.Context,
	dirPath string) error {
	const op = "ChatService.Parse"
	log := s.log.With(
		slog.String("op", op),
		slog.String("dirPath", dirPath))

	log.Info("parsing " + dirPath)

	parser1 := parser.New(s.log, s.db, s.userMessageCounterProducer)
	err := parser1.ParseFromDir(ctx, dirPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *ChatService) SearchMessages(ctx context.Context, request *dbHelper.SelectBuildRequest) ([]models.Message, error) {
	selectQuery, selectParams := dbHelper.BuildQuery[models.Message](*request)
	rows, err := s.db.QueryContext(ctx, selectQuery, selectParams...)
	if err != nil {
		return nil, err
	}

	messages, err := dbHelper.RowsToEntities[models.Message](rows)
	if err != nil {
		return messages, err
	}

	for i := 0; i < len(messages); i++ {
		messages[i].ChatName, _ = caches.ChatsCache.GetByKey(s.db, messages[i].ChatId)
		messages[i].UserName, _ = caches.UsersCache.GetByKey(s.db, messages[i].UserId)
	}

	return messages, nil
}

func (s *ChatService) DeleteMessages(ctx context.Context, request *dbHelper.SelectBuildRequest) error {
	deleteQuery, deleteParams := dbHelper.BuildDelete[models.Message](*request)
	_, err := s.db.ExecContext(ctx, deleteQuery, deleteParams...)
	return err
}
