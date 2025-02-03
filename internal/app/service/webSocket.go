package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pulse/internal/app/mediator"
	"github.com/inter-hubly/pulse/internal/domain/dto"
)

type WebSocket interface {
	ReceiveMessageFromClient(ctx context.Context)
	SendMessageToClient(ctx context.Context, d *dto.Message) error
	CreateConnection(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteConnection(ctx context.Context)
}

var (
	webSocketOnce sync.Once
	webSocket     *webSocketService
)

type webSocketService struct {
	whatsAppMediator mediator.WhatsApp
	connections      map[string]*websocket.Conn
	channelError     chan string
}

func NewWebSocket() *webSocketService {
	webSocketOnce.Do(func() {
		webSocket = &webSocketService{
			whatsAppMediator: mediator.NewWhatsApp(),
			connections:      make(map[string]*websocket.Conn),
			channelError:     make(chan string),
		}
	})
	return webSocket
}

func (s *webSocketService) SendMessageToClient(ctx context.Context, dtoMessage *dto.Message) error {
	tenant := hctx.Tenant.Get(ctx)
	conn := s.connections[tenant]
	if conn == nil {
		hlog.Info(ctx, "webSocketService.SendMessageToClient", fmt.Sprintf("error when message sending to client %v", dtoMessage))
		return fmt.Errorf("connection for tenant %s not found", tenant)
	}
	marshal, err := json.Marshal(dtoMessage)
	if err != nil {
		hlog.Info(ctx, "webSocketService.SendMessageToClient", fmt.Sprintf("error when message marshal %v: %v", err, dtoMessage))
		return err
	}
	if err = conn.WriteMessage(websocket.TextMessage, marshal); err != nil {
		hlog.Info(ctx, "webSocketService.SendMessageToClient", fmt.Sprintf("error when write message %v: %v", err, dtoMessage))
		return err
	}
	return nil
}

func (s *webSocketService) ReceiveMessageFromClient(ctx context.Context) {
	tenant := hctx.Tenant.Get(ctx)
	conn := s.connections[tenant]

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.channelError <- tenant
			break
		}
		var valueObj base.SendTextDto
		if err = json.Unmarshal(message, &valueObj); err != nil {
			hlog.Info(ctx, "webSocketService.ReceiveMessageFromClient", fmt.Sprintf("error when message unmarshal %v", err))
		}
		if err = s.whatsAppMediator.SendMessage(ctx, &valueObj); err != nil {
			hlog.Info(ctx, "webSocketService.ReceiveMessageFromClient", fmt.Sprintf("error when send message %v", err))
		}
		hlog.Info(ctx, "webSocketService.ReceiveMessageFromClient", fmt.Sprintf("Message received from client: %s", string(message)))
	}
}

func (s *webSocketService) CreateConnection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	upgrade, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		hlog.Error(ctx, "webSocketService.CreateConnection", fmt.Sprint("Error when create connection", err))
		return err
	}
	tenant := hctx.Tenant.Get(ctx)
	s.connections[tenant] = upgrade
	return nil
}

func (s *webSocketService) DeleteConnection(ctx context.Context) {
	for tenant := range s.channelError {
		hlog.Info(ctx, "webSocketService.DeleteConnection", fmt.Sprint("Deleting connection", tenant))
		delete(s.connections, tenant)
	}
}
