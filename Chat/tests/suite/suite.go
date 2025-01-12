package suite

import (
	"chat/internal/config"
	"context"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	"testing"
)

type Suite struct {
	*testing.T
	cfg        *config.Config
	ChatClient chatv1.ChatClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()
}
