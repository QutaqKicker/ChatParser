package grpc

import (
	"chat/internal/services"
	"context"
	"database/sql"
	"fmt"
	chatv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"time"
)

type Chat interface {
	ParseHtml(ctx context.Context,
		diPath string) (bool, error)
	SearchMessages(ctx context.Context,
		min time.Time,
		max time.Time,
		userIds []string) ([]*chatv1.ChatMessage, error)
	GetStatistics(ctx context.Context,
		userIds []string) (bool, error)
}

type serverAPI struct {
	chatv1.UnimplementedChatServer
	chat *services.ChatService
}

func Register(gRPC *grpc.Server, log *slog.Logger, db *sql.DB) {
	chatv1.RegisterChatServer(gRPC, &serverAPI{chat: services.NewChatService(log, db)})
}

func (s *serverAPI) ParseHtml(ctx context.Context, req *chatv1.ParseHtmlRequest) (*chatv1.ParseHtmlResponse, error) {
	if req.DirPath == "" {
		return nil, status.Error(codes.InvalidArgument, "dirPath is empty")
	}

	err := s.chat.Parse(ctx, req.DirPath)
	isSuccess := true
	if err != nil {
		isSuccess = false
	}
	return &chatv1.ParseHtmlResponse{IsSuccess: isSuccess}, err
}

func (s *serverAPI) SearchMessages(ctx context.Context, req *chatv1.SearchMessagesRequest) (*chatv1.SearchMessagesResponse, error) {
	if req.MinDate == nil && req.MaxDate == nil && req.UserIds == nil {
		return nil, status.Error(codes.InvalidArgument, "all queryFilters is empty")
	}

	//messages, err := s.chat.SearchMessages(ctx, req.MinDate.AsTime(), req.MaxDate.AsTime(), req.UserIds)
	return &chatv1.SearchMessagesResponse{Messages: nil}, nil
}

func (s *serverAPI) GetStatistics(ctx context.Context, req *chatv1.GetStatisticsRequest) (*chatv1.GetStatisticsResponse, error) {
	isSuccess, err := s.chat.GetStatistics(ctx, req.UserIds)
	return &chatv1.GetStatisticsResponse{IsSuccess: isSuccess}, err
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	db         *sql.DB
	port       int
}

func New(log *slog.Logger, db *sql.DB, port int) *App {
	grpcServer := grpc.NewServer()
	Register(grpcServer, log, db)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		db:         db,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "ChatApp.Run"

	log := a.log.With(slog.String("op", op))

	log.Info("starting Grpc server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "ChatApp.Stop"
	log := a.log.With(slog.String("op", op))
	log.Info("stopping Grpc server")
	a.gRPCServer.GracefulStop()
}
