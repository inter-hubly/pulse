package service

import (
	"context"
	"net/http"
	"sync"

	"github.com/inter-hubly/pulse/internal/app/mediator"
	"github.com/inter-hubly/pulse/internal/domain/valueobject"
)

type WebSocket interface {
	SendMessageToClient(ctx context.Context, ownerId, toId string) ([]*valueobject.Message, error)
	ReceiveMessageFromClient(ctx context.Context, ownerId, toId string) ([]*valueobject.Message, error)
}

var (
	webSocketOnce sync.Once
	webSocket     *webSocketService
)

type webSocketService struct {
	whatsAppMediator mediator.WhatsApp
}

func NewWebSocket(w http.ResponseWriter, r *http.Request) *webSocketService {
	webSocketOnce.Do(func() {
		webSocket = &webSocketService{
			whatsAppMediator: mediator.NewWhatsApp(),
		}
	})
	return webSocket
}

func (s *webSocketService) SendMessageToClient(ctx context.Context, ownerId, toId string) ([]*valueobject.Message, error) {

}

func (s *webSocketService) ReceiveMessageFromClient(ctx context.Context, ownerId, toId string) ([]*valueobject.Message, error) {

}
