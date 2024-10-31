package services

import (
	"chat/internal/domain/filters"
	"chat/internal/domain/models"
	"chat/internal/parser"
	"context"
	"log/slog"
	"time"
)

type ChatService struct {
	log *slog.Logger
	parser.Parser
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
		filter filters.QueryBuildRequest) ([]models.Message, error)
}

// New returns new instance of chat service
func New(
	log *slog.Logger,
	parser parser.Parser,
	messageSaver MessageSaver,
	messageProvider MessageProvider) *ChatService {
	return &ChatService{
		log:             log,
		Parser:          parser,
		MessageSaver:    messageSaver,
		MessageProvider: messageProvider,
	}
}

func (s *ChatService) Parse(ctx context.Context,
	dirPath string) error {
	const op = "chatService.Parse"
	log := s.log.With(
		slog.String("op", op),
		slog.String("dirPath", dirPath))

	log.Info("parsing " + dirPath)

	err := s.ParseFromDir(ctx, dirPath)
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
