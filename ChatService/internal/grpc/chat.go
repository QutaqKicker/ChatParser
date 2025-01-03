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

type ChatService interface {
	ParseFromDir(ctx context.Context,
		diPath string) error
	SearchMessages(ctx context.Context,
		min time.Time,
		max time.Time,
		userIds []string) ([]*chatv1.ChatMessage, error)
	GetStatistics(ctx context.Context,
		userIds []string) (bool, error)
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

type serverAPI struct {
	chatv1.UnimplementedChatServer
	chat *services.ChatService
}

func Register(gRPC *grpc.Server, log *slog.Logger, db *sql.DB) {
	chatv1.RegisterChatServer(gRPC, &serverAPI{chat: services.NewChatService(log, db)})
}

func (s *serverAPI) ParseFromDir(ctx context.Context, req *chatv1.ParseFromDirRequest) (*chatv1.ParseFromDirResponse, error) {
	if req.DirPath == "" {
		return nil, status.Error(codes.InvalidArgument, "dirPath is empty")
	}

	err := s.chat.Parse(ctx, req.DirPath)
	ok := true
	if err != nil {
		ok = false
	}
	return &chatv1.ParseFromDirResponse{Ok: ok}, err
}

func (s *serverAPI) GetMessages(ctx context.Context, req *chatv1.SearchMessagesRequest) (*chatv1.GetMessagesResponse, error) {
	if req.Filter == nil {
		return nil, status.Error(codes.InvalidArgument, "all filters is empty")
	}

	//messages, err := s.ChatService.SearchMessages(ctx, req.MinDate.AsTime(), req.MaxDate.AsTime(), req.UserIds)
	return &chatv1.GetMessagesResponse{}, nil
}

func (s *serverAPI) DeleteMessages(ctx context.Context, req *chatv1.SearchMessagesRequest) (*chatv1.DeleteMessageResponse, error) {
	if req.Filter == nil {
		return nil, status.Error(codes.InvalidArgument, "all filters is empty")
	}

	//messages, err := s.ChatService.SearchMessages(ctx, req.MinDate.AsTime(), req.MaxDate.AsTime(), req.UserIds)
	return &chatv1.DeleteMessageResponse{}, nil
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
