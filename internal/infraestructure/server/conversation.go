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
	handleMessage(ctx context.Context)
	ReadMessage(ctx context.Context)
	WriteMessage(ctx context.Context, msg valueobject.Message)
}

type conversation struct {
	connection     *websocket.Conn
	messageChannel chan valueobject.Message
}

func NewConversation(ctx context.Context, w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*conversation, error) {
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
	}
	go conv.handleMessage(ctx)
	return conv, nil
}

func (c *conversation) handleMessage(ctx context.Context) {
	for msg := range c.messageChannel {
		c.connection.WriteJSON(msg)
	}
}

func (c *conversation) ReadMessage(ctx context.Context) {
	for {
		hlog.Info("webSocketServer.HandleWebSocket.Connection", "Reading Message")
		_, message, err := c.connection.ReadMessage()
		if err != nil && websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
			hlog.Info("webSocketServer.HandleWebSocket.Connection", "Closing Connection")
			return
		}
		var msg valueobject.Message
		if err = json.Unmarshal(message, &msg); err != nil {
			hlog.Error("webSocketServer.HandleWebSocket.Unmarshal", fmt.Sprintf("err: %v", err))
		}
		c.messageChannel <- msg
	}
}

func (c *conversation) WriteMessage(ctx context.Context, msg valueobject.Message) {
	c.messageChannel <- msg
}
