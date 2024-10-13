package grpc

import (
	"context"
	"fmt"
	backupv1 "github.com/QutaqKicker/ChatParser/protos/gen/go/backup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
)

type Backup interface {
	ExportToFile(ctx context.Context,
		exportType backupv1.ExportType,
		exportDir string) (bool, error)
	ImportFromFile(ctx context.Context,
		exportDir string) (bool, error)
}

type serverAPI struct {
	backupv1.UnimplementedBackupServer
	backup Backup
}

func Register(gRPC *grpc.Server) {
	backupv1.RegisterBackupServer(gRPC, &serverAPI{})
}

func (s *serverAPI) ExportToFile(ctx context.Context, req *backupv1.ExportToFileRequest) (*backupv1.ExportToFileResponse, error) {
	if req.ExportDir == "" {
		return nil, status.Error(codes.InvalidArgument, "exportDir is empty")
	}
	isSuccess, err := s.backup.ExportToFile(ctx, req.Type, req.ExportDir)
	return &backupv1.ExportToFileResponse{IsSuccess: isSuccess}, err
}

func (s *serverAPI) ImportFromFile(ctx context.Context, req *backupv1.ImportFromFileRequest) (*backupv1.ImportFromFileResponse, error) {
	if req.ExportDir == "" { //TODO Поменять на нормальное название
		return nil, status.Error(codes.InvalidArgument, "exportDir is empty")
	}
	isSuccess, err := s.backup.ImportFromFile(ctx, req.ExportDir)
	return &backupv1.ImportFromFileResponse{IsSuccess: isSuccess}, err
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
	const op = "BackupApp.Run"

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
	const op = "BackupApp.Stop"
	a.gRPCServer.GracefulStop()
}
