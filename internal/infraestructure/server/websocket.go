package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pulse/internal/app/mediator"
)

var (
	webSocketOnce   sync.Once
	webSocketServer *webSocket
)

type WebSocket interface {
	HandleWebSocket(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

type webSocket struct {
	connection       map[string]map[string]WebSocketConversation
	whatsAppMediator mediator.WhatsApp
	serverConv       WebSocketServer
}

func NewWebSocket() *webSocket {
	webSocketOnce.Do(func() {
		webSocketServer = &webSocket{
			whatsAppMediator: mediator.NewWhatsApp(),
			serverConv:       NewServerConversation(),
		}
	})

	return webSocketServer
}

func (s *webSocket) HandleWebSocket(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	toNumber := r.URL.Query().Get("toNumber")
	if toNumber == "" {
		return
	}

	connection, err := s.getConnection(ctx, w, r, toNumber)
	if err != nil {
		hlog.Error("webSocket.HandleWebSocket", fmt.Sprintf("Error: %s", err))
	}

	connection.ReadMessage(ctx)
}

func (s *webSocket) getConnection(ctx context.Context, w http.ResponseWriter, r *http.Request, numberId string) (WebSocketConversation, error) {
	connection, err := s.serverConv.GetConnection(ctx, numberId)

	// caso a conexão não exista
	if err != nil {
		hlog.Error("webSocket.HandleWebSocket", fmt.Sprintf("Error: %s", err))
		newConversation, err := NewConversation(ctx, w, r, nil)
		if err != nil {
			return nil, err
		}
		if err = s.serverConv.SaveConnection(ctx, numberId, newConversation); err != nil {
			return nil, err
		}
		return newConversation, nil
	}
	return connection, nil
}
