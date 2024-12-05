package services

import (
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"chat/internal/parser"
	"context"
	"database/sql"
	"log/slog"
	"time"
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

// New returns new instance of chat service
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
	const op = "chatService.Parse"
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

func (s *ChatService) SearchMessages(ctx context.Context,
	min time.Time,
	max time.Time,
	userIds []string) ([]*models.Message, error) {
	panic("not implemented")
}

func (s *ChatService) GetStatistics(ctx context.Context,
	userIds []string) (bool, error) {
	panic("not implemented")
}
