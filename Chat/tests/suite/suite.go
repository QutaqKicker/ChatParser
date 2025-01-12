package suite

import (
	"chat/internal/config"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	"testing"
)

type Suite struct {
	*testing.T
	cfg        *config.Config
	ChatClient chatv1.ChatClient
}
