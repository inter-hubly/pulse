package controller

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/inter-hubly/pilot/hctx"
	"net/http"
	"sync"

	"github.com/inter-hubly/pulse/internal/app/service"
	"github.com/inter-hubly/pulse/internal/infraestructure/server"
)

type WebSocket interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

var (
	webSocketOnce       sync.Once
	webSocketController *webSocket
)

type webSocket struct {
	webSocketService service.WebSocket
	connections      map[string]websocket.Conn
}

func NewWebSocket() *webSocket {
	webSocketOnce.Do(func() {
		webSocketController = &webSocket{
			webSocketService: server.NewWebSocket(),
		}
	})
	return webSocketController
}

func (ws *webSocket) Handle(w http.ResponseWriter, r *http.Request) {
	ownerId := r.URL.Query().Get("user")
	if ownerId == "" {
		return
	}

	toPhone := r.Header.Get("to-phone")
	if toPhone == "" {
		return
	}
	ctx := hctx.Tenant.New(ownerId)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Permitir qualquer origem
			return true
		},
	}
	if conn, exists := ws.connections[ownerId]; exists {

	} else {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Erro ao atualizar para WebSocket:", err)
			return
		}
	}

	ws.webSocketService.ReceiveMessageFromClient(ctx, ownerId, toPhone)
}
