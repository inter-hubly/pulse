package service

import (
	"context"
	"github.com/gorilla/websocket"
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
	connections      map[string]websocket.Conn
}

func NewWebSocket(w http.ResponseWriter, r *http.Request) *webSocketService {
	webSocketOnce.Do(func() {
		webSocket = &webSocketService{
			whatsAppMediator: mediator.NewWhatsApp(),
			connections:      make(map[string]websocket.Conn),
		}
	})
	return webSocket
}

func (s *webSocketService) SendMessageToClient(ctx context.Context, ownerId, toId string) ([]*valueobject.Message, error) {

}

func (s *webSocketService) ReceiveMessageFromClient(ctx context.Context, ownerId, toId string) ([]*valueobject.Message, error) {

}
