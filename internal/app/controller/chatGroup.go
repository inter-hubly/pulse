package controller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/inter-hubly/pulse/internal/app/service"
	"github.com/inter-hubly/pulse/internal/domain/dto"
	"github.com/inter-hubly/pulse/internal/infraestructure/server"
)

type ChatGroup interface {
	Handle(w http.ResponseWriter, r *http.Request)
	GetAllMessages(w http.ResponseWriter, r *http.Request)
	ReceiveMessage(w http.ResponseWriter, r *http.Request)
	HandleStaticFiles(w http.ResponseWriter, r *http.Request)
}

var (
	chatGroupOnce       sync.Once
	chatGroupController *chatGroup
)

type chatGroup struct {
	channel         chan dto.Message
	websocket       server.WebSocket
	whatsAppService service.WhatsApp
}

func NewChatGroup() *chatGroup {
	chatGroupOnce.Do(func() {
		chatGroupController = &chatGroup{
			channel:         make(chan dto.Message),
			websocket:       server.NewWebSocket(),
			whatsAppService: service.NewWhatsApp(),
		}
	})
	return chatGroupController
}

func (c *chatGroup) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	c.websocket.HandleWebSocket(ctx, w, r, c.channel)
}

func (c *chatGroup) HandleStaticFiles(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("."))
	fs.ServeHTTP(w, r)
}

func (c *chatGroup) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	OwnerId := r.URL.Query().Get("user")
	if OwnerId == "" {
		return
	}

	message, err := c.whatsAppService.GetAllMessage(context.Background(), OwnerId, "+5548991784586")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	resp := make([]dto.Message, 0)
	for _, v := range message {
		resp = append(resp, dto.Message{
			Username: OwnerId,
			Message:  v.Message,
		})
	}
	marshal, err := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
	w.Header().Add("Content-Type", "application/json")
}

func (c *chatGroup) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	go func() {
		var message dto.Message
		if err = json.Unmarshal(body, &message); err != nil {
			return
		}
		chatGroupController.channel <- message
	}()

	return
}
