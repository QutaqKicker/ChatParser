package grpc

import (
	"context"
	"fmt"
	userv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/user"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type User interface {
	GetUsers(ctx context.Context) ([]*userv1.UserInfo, error)
	SearchUser(ctx context.Context,
		userId string) (string, int64, error)
	EditUser(ctx context.Context,
		userId string,
		newName string) (bool, error)
}

type serverAPI struct {
	userv1.UnimplementedUserServer
	user User
}

func Register(gRPC *grpc.Server) {
	userv1.RegisterUserServer(gRPC, &serverAPI{})
}

func (s *serverAPI) GetUsers(ctx context.Context, req *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	users, err := s.user.GetUsers(ctx)
	return &userv1.GetUsersResponse{Users: users}, err
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
