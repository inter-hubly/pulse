package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pulse/internal/domain/dto"
)

type WebSocket interface {
	HandleWebSocket(w http.ResponseWriter, r *http.Request, sendChan <-chan dto.Message)
}

var (
	webSocketOnce   sync.Once
	webSocketServer *webSocket
)

type webSocket struct {
	connection *websocket.Conn
	upGrader   websocket.Upgrader
	broker     broker.Connection
}

func NewWebSocket() *webSocket {
	webSocketOnce.Do(func() {
		webSocketServer = &webSocket{
			broker: broker.GetConnection(),
			upGrader: websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin: func(r *http.Request) bool {
					hlog.Info("webSocketServer.HandleWebSocket.Connection", "Opening Connection")
					return true
				},
			},
		}
	})

	return webSocketServer
}

func (s *webSocket) HandleWebSocket(w http.ResponseWriter, r *http.Request, sendChan <-chan dto.Message) {
	var err error
	s.connection, err = s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		hlog.Error("webSocketServer.HandleWebSocket.conn", fmt.Sprintf("err: %v", err))
		return
	}
	defer s.connection.Close()

	go func() {
		for msg := range sendChan {
			if err = s.connection.WriteJSON(msg); err != nil {
				hlog.Error("webSocketServer.HandleWebSocket.WriteJSON", fmt.Sprintf("websocket error sending message: %v", err))
				break
			}
		}
	}()

	for {
		messageType, message, err := s.connection.ReadMessage()
		if err != nil && websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			hlog.Info("webSocketServer.HandleWebSocket.Connection", "Closing Connection")
			break
		}
		if err = s.broker.Publish("whatsapp.send", message); err != nil {
			hlog.Error("webSocketServer.HandleWebSocket.Publish", fmt.Sprintf("websocket error : %v", err))
		} else {
			if err = s.connection.WriteMessage(messageType, message); err != nil {
				hlog.Error("webSocketServer.HandleWebSocket.WriteMessage", fmt.Sprintf("websocket error : %v", err))
			}
		}
	}
}
