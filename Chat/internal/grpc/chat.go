package grpc

import (
	"chat/internal/domain/filters"
	"chat/internal/services"
	"context"
	"database/sql"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/dbHelper"
	"github.com/QutaqKicker/ChatParser/Common/myKafka"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func New(log *slog.Logger, db *sql.DB, port int, userMessageCounterProducer *myKafka.UserMessageCounterProducer) *App {
	grpcServer := grpc.NewServer()
	Register(grpcServer, log, db, userMessageCounterProducer)

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

func Register(gRPC *grpc.Server, log *slog.Logger, db *sql.DB, userMessageCounterProducer *myKafka.UserMessageCounterProducer) {
	chatv1.RegisterChatServer(gRPC, &serverAPI{chat: services.NewChatService(log, db, userMessageCounterProducer)})
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

func (s *serverAPI) GetMessages(ctx context.Context, request *chatv1.SearchMessagesRequest) (*chatv1.GetMessagesResponse, error) {
	messages, err := s.chat.SearchMessages(ctx, searchMessagesRequestConvertGrpcToDbHelper(request))
	if err != nil {
		return nil, err
	}

	response := &chatv1.GetMessagesResponse{Messages: make([]*chatv1.ChatMessage, len(messages))}
	for i := 0; i < len(messages); i++ {
		response.Messages[i] = &chatv1.ChatMessage{
			Id:               messages[i].Id,
			ChatId:           messages[i].ChatId,
			ChatName:         messages[i].ChatName,
			UserId:           messages[i].UserId,
			UserName:         messages[i].UserName,
			ReplyToMessageId: messages[i].ReplyToMessageId,
			Text:             messages[i].Text,
			Created:          timestamppb.New(messages[i].Created),
		}
	}
	return response, err
}

func (s *serverAPI) DeleteMessages(ctx context.Context, request *chatv1.SearchMessagesRequest) (*chatv1.DeleteMessageResponse, error) {
	err := s.chat.DeleteMessages(ctx, searchMessagesRequestConvertGrpcToDbHelper(request))
	return &chatv1.DeleteMessageResponse{Ok: err == nil}, err
}

func searchMessagesRequestConvertGrpcToDbHelper(req *chatv1.SearchMessagesRequest) *dbHelper.SelectBuildRequest {
	var minCreatedDate, maxCreatedDate time.Time
	if req.Filter.MinCreatedDate != nil {
		minCreatedDate = req.Filter.MinCreatedDate.AsTime()
	}

	if req.Filter.MinCreatedDate != nil {
		maxCreatedDate = req.Filter.MaxCreatedDate.AsTime()
	}

	return &dbHelper.SelectBuildRequest{
		Filter: filters.MessageFilter{
			Id:             req.Filter.Id,
			MinCreatedDate: minCreatedDate,
			MaxCreatedDate: maxCreatedDate,
			SubText:        req.Filter.SubText,
			UserId:         req.Filter.UserId,
			UserIds:        req.Filter.UserIds,
			ChatIds:        req.Filter.ChatIds,
		},
		Sorts: []dbHelper.SortField{{FieldName: "created", Direction: dbHelper.Desc}},
		Take:  int(req.Take),
		Skip:  int(req.Skip),
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
