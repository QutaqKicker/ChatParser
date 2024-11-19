package grpc

import (
	"context"
	"fmt"
	auditv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/audit"
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
		userIds []string) ([]*auditv1.ChatMessage, error)
	GetStatistics(ctx context.Context,
		userIds []string) (bool, error)
}

type serverAPI struct {
	auditv1.UnimplementedAuditServer
	chat Chat
}

func Register(gRPC *grpc.Server) {
	auditv1.RegisterChatServer(gRPC, &serverAPI{})
}

func (s *serverAPI) ParseHtml(ctx context.Context, req *auditv1.ParseHtmlRequest) (*auditv1.ParseHtmlResponse, error) {
	if req.DirPath == "" {
		return nil, status.Error(codes.InvalidArgument, "dirPath is empty")
	}

	isSuccess, err := s.chat.ParseHtml(ctx, req.DirPath)
	return &auditv1.ParseHtmlResponse{IsSuccess: isSuccess}, err
}

func (s *serverAPI) SearchMessages(ctx context.Context, req *auditv1.SearchMessagesRequest) (*auditv1.SearchMessagesResponse, error) {
	if req.MinDate == nil && req.MaxDate == nil && req.UserIds == nil {
		return nil, status.Error(codes.InvalidArgument, "all queryFilters is empty")
	}

	messages, err := s.chat.SearchMessages(ctx, req.MinDate.AsTime(), req.MaxDate.AsTime(), req.UserIds)
	return &auditv1.SearchMessagesResponse{Messages: messages}, err
}

func (s *serverAPI) GetStatistics(ctx context.Context, req *auditv1.GetStatisticsRequest) (*auditv1.GetStatisticsResponse, error) {
	isSuccess, err := s.chat.GetStatistics(ctx, req.UserIds)
	return &auditv1.GetStatisticsResponse{IsSuccess: isSuccess}, err
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *App {
	grpcServer := grpc.NewServer()
	Register(grpcServer)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
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
