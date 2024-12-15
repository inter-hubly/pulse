package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pulse/internal/domain/valueobject"
)

type WebSocketConversation interface {
	ReadMessage(ctx context.Context, chanError chan error)
	WriteMessage(ctx context.Context, msg valueobject.Message)
}

type WebSocketMediator interface {
	Notify(ctx context.Context, msg valueobject.Message) error
}

type conversation struct {
	connection     *websocket.Conn
	messageChannel chan valueobject.Message
	mediator       WebSocketMediator
}

func NewConversation(ctx context.Context,
	w http.ResponseWriter, r *http.Request,
	responseHeader http.Header,
	mediator WebSocketMediator,
) (*conversation, error) {
	updater := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			hlog.Info("webSocketServer.HandleWebSocket.Connection", "Opening Connection")
			return true
		},
	}
	upgrade, err := updater.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	conv := &conversation{
		connection:     upgrade,
		messageChannel: make(chan valueobject.Message),
		mediator:       mediator,
	}
	return conv, nil
}

func (c *conversation) ReadMessage(ctx context.Context, chanError chan error) {
	for {
		hlog.Info("webSocketServer.HandleWebSocket.Connection", "Reading Message")
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			hlog.Info("webSocketServer.HandleWebSocket.Connection", "Closing Connection")
			return
		}
		var msg valueobject.Message
		if err = json.Unmarshal(message, &msg); err != nil {
			chanError <- err
			hlog.Error("webSocketServer.HandleWebSocket.Unmarshal", fmt.Sprintf("err: %v", err))
		}

		c.mediator.Notify(ctx, msg)
	}
}

func (c *conversation) WriteMessage(ctx context.Context, msg valueobject.Message) {
	c.connection.WriteJSON(msg)
	hlog.Info("webSocketServer.handleMessage", fmt.Sprint(msg))
}
