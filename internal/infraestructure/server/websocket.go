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
	GetConnection(ctx context.Context) (WebSocketConversation, error)
}

type webSocket struct {
	connection       map[string]WebSocketConversation
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
	var connection WebSocketConversation
	connection, err := s.getConnection(ctx, w, r)
	if err != nil {
		hlog.Error("webSocket.HandleWebSocket", fmt.Sprintf("Error: %s", err))
	}

	connection.ReadMessage(ctx)

	s.serverConv.DeleteConnection(ctx)
}

func (s *webSocket) getConnection(ctx context.Context, w http.ResponseWriter, r *http.Request) (WebSocketConversation, error) {
	connection, err := s.serverConv.GetConnection(ctx)

	// caso a conexão não exista
	if err != nil {
		hlog.Error("webSocket.HandleWebSocket", fmt.Sprintf("Error: %s", err))
		newConversation, err := NewConversation(ctx, w, r, nil)
		if err != nil {
			return nil, err
		}
		if err = s.serverConv.SaveConnection(ctx, newConversation); err != nil {
			return nil, err
		}
		return newConversation, nil
	}
	return connection, nil
}

func (s *webSocket) GetConnection(ctx context.Context) (WebSocketConversation, error) {
	return s.serverConv.GetConnection(ctx)
}
