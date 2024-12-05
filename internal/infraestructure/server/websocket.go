package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pulse/internal/app/mediator"
	"github.com/inter-hubly/pulse/internal/domain/dto"
	"github.com/inter-hubly/pulse/internal/domain/valueobject"
)

type WebSocket interface {
	HandleWebSocket(ctx context.Context, w http.ResponseWriter, r *http.Request, sendChan <-chan dto.Message)
}

var (
	webSocketOnce   sync.Once
	webSocketServer *webSocket
)

type webSocket struct {
	connection       *websocket.Conn
	upGrader         websocket.Upgrader
	whatsAppMediator mediator.WhatsApp
}

func NewWebSocket() *webSocket {
	webSocketOnce.Do(func() {
		webSocketServer = &webSocket{
			whatsAppMediator: mediator.NewWhatsApp(),
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

func (s *webSocket) HandleWebSocket(ctx context.Context, w http.ResponseWriter, r *http.Request, sendChan <-chan dto.Message) {
	var err error
	s.connection, err = s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		hlog.Error("webSocketServer.HandleWebSocket.conn", fmt.Sprintf("err: %v", err))
		return
	}
	defer s.connection.Close()

	OwnerId, err := getUserId(r)
	if err != nil {
		hlog.Error("webSocketServer.HandleWebSocket.getUserId", fmt.Sprintf("err: %v", err))
	}

	go func() {
		for msg := range sendChan {
			err = s.whatsAppMediator.ReceiveMessage(ctx, OwnerId, msg)
			if err != nil {
				hlog.Error("webSocketServer.HandleWebSocket.receiveMessage", fmt.Sprintf("error receiving message: %v", err))
			}
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
		var msg valueobject.Message
		if err = json.Unmarshal(message, &msg); err != nil {
			hlog.Error("webSocketServer.HandleWebSocket.Unmarshal", fmt.Sprintf("err: %v", err))
		}

		if err = s.whatsAppMediator.SendMessage(ctx, OwnerId, &msg); err != nil {
			hlog.Error("webSocketServer.HandleWebSocket.Publish", fmt.Sprintf("websocket error : %v", err))
		} else {
			if err = s.connection.WriteMessage(messageType, message); err != nil {
				hlog.Error("webSocketServer.HandleWebSocket.WriteMessage", fmt.Sprintf("websocket error : %v", err))
			}
		}
	}
}

func getUserId(r *http.Request) (string, error) {
	OwnerId := r.URL.Query().Get("user")
	if OwnerId == "" {
		return "", errors.New("no user id")
	}
	return OwnerId, nil
}
