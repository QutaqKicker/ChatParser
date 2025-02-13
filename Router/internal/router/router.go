package router

import (
	"context"
	"fmt"
	backupv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/backup"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	userv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
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

func NewRouter(logger *slog.Logger, routerPort, chatPort, userPort, backupPort string) http.Handler {
	mux := http.NewServeMux()

	cc, err := grpc.NewClient(fmt.Sprintf("localhost:%d", chatPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*20)))
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	chatClient := chatv1.NewChatClient(cc)

	uc, err := grpc.NewClient(fmt.Sprintf("localhost:%d", userPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions())
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	userClient := userv1.NewUserClient(uc)

	bc, err := grpc.NewClient(fmt.Sprintf("localhost:%d", backupPort),
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
	mux.Handle("/chat/messages/", getMessagesHandler(logger, chatClient))
	mux.Handle("/chat/messages/", handleTenantsGet(logger, tenantsStore))
	mux.Handle("/oauth2/", handleOAuth2Proxy(logger, authProxy))
	mux.HandleFunc("/healthz", handleHealthzPlease(logger))
	mux.Handle("/", http.NotFoundHandler())
}

func getMessagesHandler(logger *slog.Logger, chatClient *chatv1.ChatClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)
		if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
			r2 := new(Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			r2.URL.RawPath = rp
			h.ServeHTTP(w, r2)
		} else {
			NotFound(w, r)
		}
	})
}

func (s *serverAPI) ParseHtml(ctx context.Context, req *routerv1.ParseHtmlRequest) (*routerv1.ParseHtmlResponse, error) {
	if req.DirPath == "" {
		return nil, status.Error(codes.InvalidArgument, "dirPath is empty")
	}
	isSuccess, err := s.router.ParseHtml(ctx, req.DirPath)
	return &routerv1.ParseHtmlResponse{IsSuccess: isSuccess}, err
}
