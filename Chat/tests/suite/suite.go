package suite

import (
	"chat/internal/domain/models"
	"context"
	"github.com/QutaqKicker/ChatParser/Common/config"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	cfg        *models.Config
	ChatClient chatv1.ChatClient
}

// New creates new test suite.
//
// TODO: for pipeline tests we need to wait for app is ready
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath[models.Config](configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // Используем insecure-коннект для тестов
	if err != nil {
		t.Fatalf("router server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/local_tests.yaml"
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
