package mediator

import (
	"context"
	"sync"

	"github.com/inter-hubly/pulse/internal/domain/dto"
)

type WebSocket interface {
	Handler(context.Context, Messages) error
}

type Messages interface {
	GetMessages(ctx context.Context) []dto.Message
}

var (
	webSocketOnce sync.Once
	webSocket     *webSocketService
)

type webSocketService struct {
}

func NewWebSocket() *webSocketService {
	webSocketOnce.Do(func() {
		webSocket = &webSocketService{}
	})
	return webSocket
}

func (s *webSocketService) Handler(ctx context.Context, msg Messages) error {
	return nil
}
