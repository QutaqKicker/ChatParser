package grpc

import (
	"backups/internal/exporter"
	"context"
	"fmt"
	backupv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/backup"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	commonv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/common"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Backup interface {
	ExportToDir(ctx context.Context, exportType backupv1.ExportType, messageFilter *commonv1.MessagesFilter) error
}

type serverAPI struct {
	backupv1.UnimplementedBackupServer
	backup Backup
}

func Register(gRPC *grpc.Server, chatClient chatv1.ChatClient, exportDir string) {
	backupv1.RegisterBackupServer(gRPC, &serverAPI{backup: exporter.NewExporter(chatClient, exportDir)})
}

func (s *serverAPI) ExportToDir(ctx context.Context, req *backupv1.ExportToDirRequest) (*backupv1.ExportToDirResponse, error) {
	err := s.backup.ExportToDir(ctx, req.Type, req.MessageFilter)
	return &backupv1.ExportToDirResponse{Ok: err == nil}, err
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	exportDir  string
}

func New(log *slog.Logger, exportDir string, port int, chatServicePort int) *App {
	cc, err := grpc.NewClient(fmt.Sprintf("localhost:%d", chatServicePort))
	if err != nil {
		log.Error("grpc server connection failed: %v", err)
		return nil
	}

	chatClient := chatv1.NewChatClient(cc)

	grpcServer := grpc.NewServer()
	Register(grpcServer, chatClient, exportDir)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
		exportDir:  exportDir,
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
