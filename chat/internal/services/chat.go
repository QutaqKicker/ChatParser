package services

import (
	"chat/internal/domain/models"
	"context"
	"log/slog"
	"time"
)

type ChatService struct {
	log *slog.Logger
	HtmlParser
	MessageSaver
	MessageProvider
}

type HtmlParser interface {
	ParseFromDir(
		ctx context.Context,
		directory string) error
}

type MessageSaver interface {
	SaveMessages(
		ctx context.Context,
		message []models.Message) error
}

type MessageProvider interface {
	GetMessages(
		filter models.MessageFilter) ([]models.Message, error)
}

// New returns new instance of chat service
func New(
	log *slog.Logger,
	htmlParser HtmlParser,
	messageSaver MessageSaver,
	messageProvider MessageProvider) *ChatService {
	return &ChatService{
		log:             log,
		HtmlParser:      htmlParser,
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

	panic("not implemented")
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
