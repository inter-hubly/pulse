package controller

import (
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
	ws.webSocketService.
}
