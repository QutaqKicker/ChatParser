package grpc

import (
	"context"
	"fmt"
	routerv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/Router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"time"
)

type Chat interface {
	ParseHtml(ctx context.Context,
		diPath string)
	SearchByDate(ctx context.Context,
		min time.Time,
		max time.Time)
	SearchByUser(ctx context.Context,
		userId string)
	GetStatistics(ctx context.Context,
		userId string)
}

type Router interface {
	ParseHtml(ctx context.Context,
		diPath string) (bool, error)
}

type serverAPI struct {
	routerv1.UnimplementedRouterServer
	router Router
}

func Register(gRPC *grpc.Server) {
	routerv1.RegisterRouterServer(gRPC, &serverAPI{})
}

func (s *serverAPI) ParseHtml(ctx context.Context, req *routerv1.ParseHtmlRequest) (*routerv1.ParseHtmlResponse, error) {
	if req.DirPath == "" {
		return nil, status.Error(codes.InvalidArgument, "dirPath is empty")
	}
	isSuccess, err := s.router.ParseHtml(ctx, req.DirPath)
	return &routerv1.ParseHtmlResponse{IsSuccess: isSuccess}, err
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
	const op = "RouterApp.Run"

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
	const op = "RouterApp.Stop"
	a.gRPCServer.GracefulStop()
}
