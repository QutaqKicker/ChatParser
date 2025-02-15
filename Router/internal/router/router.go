package router

import (
	"Router/internal/handlers"
	"context"
	"fmt"
	backupv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/backup"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	userv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
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

func NewRouter(logger *slog.Logger, chatPort, userPort, backupPort string) http.Handler {
	mux := http.NewServeMux()

	cc, err := grpc.NewClient(fmt.Sprintf("localhost:%s", chatPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*20)))
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	chatClient := chatv1.NewChatClient(cc)

	uc, err := grpc.NewClient(fmt.Sprintf("localhost:%s", userPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions())
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	userClient := userv1.NewUserClient(uc)

	bc, err := grpc.NewClient(fmt.Sprintf("localhost:%s", backupPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions())
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	backupClient := backupv1.NewBackupClient(bc)

	addRoutes(
		mux,
		logger,
		&chatClient,
		&userClient,
		&backupClient,
	)

	return mux
}

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	chatClient *chatv1.ChatClient,
	userClient *userv1.UserClient,
	backupClient *backupv1.BackupClient,

) {
	mux.Handle("/chat/messages/search", handlers.SearchMessagesHandler(logger, chatClient))
	mux.Handle("/chat/messages/count", handlers.GetMessagesCountHandler(logger, chatClient))
	mux.Handle("/chat/parse-from-dir", handlers.ParseFromDirHandler(logger, chatClient))
	mux.Handle("/backup/export-to-dir", handlers.ExportToDirHandler(logger, backupClient))
	mux.Handle("/user/messages-count", handlers.GetUsersWithMessagesCountHandler(logger, userClient))
	mux.Handle("/", http.NotFoundHandler())
}
