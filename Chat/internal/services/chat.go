package services

import (
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"chat/internal/parser"
	"context"
	"database/sql"
	"github.com/QutaqKicker/ChatParser/common/dbHelper"
	"log/slog"
)

type ChatService struct {
	log *slog.Logger
	db  *sql.DB
	MessageSaver
	MessageProvider
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
	db *sql.DB) *ChatService {
	return &ChatService{
		log: log,
		db:  db,
	}
}

func (s *ChatService) Parse(ctx context.Context,
	dirPath string) error {
	const op = "ChatService.Parse"
	log := s.log.With(
		slog.String("op", op),
		slog.String("dirPath", dirPath))

	log.Info("parsing " + dirPath)

	parser1 := parser.New(s.log, s.db)
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

	return dbHelper.RowsToEntities[models.Message](rows)
}

func (s *ChatService) DeleteMessages(ctx context.Context, request *dbHelper.SelectBuildRequest) error {
	deleteQuery, deleteParams := dbHelper.BuildDelete[models.Message](*request)
	_, err := s.db.ExecContext(ctx, deleteQuery, deleteParams...)
	return err
}
