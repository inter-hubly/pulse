package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pulse/internal/app/service"
	"github.com/inter-hubly/pulse/internal/domain/dto"
	"github.com/inter-hubly/pulse/internal/infraestructure/server"
)

type ChatGroup interface {
	Handle(w http.ResponseWriter, r *http.Request)
	GetAllMessages(w http.ResponseWriter, r *http.Request)
	ReceiveMessage(w http.ResponseWriter, r *http.Request)
	HandleStaticFiles(w http.ResponseWriter, r *http.Request)
	StartTemplate(w http.ResponseWriter, r *http.Request)
}

var (
	chatGroupOnce       sync.Once
	chatGroupController *chatGroup
)

type chatGroup struct {
	whatsAppService service.WhatsApp
	webSocket       server.WebSocket
}

func NewChatGroup() *chatGroup {
	chatGroupOnce.Do(func() {
		chatGroupController = &chatGroup{
			whatsAppService: service.NewWhatsApp(),
			webSocket:       server.NewWebSocket(),
		}
	})
	return chatGroupController
}

func (c *chatGroup) Handle(w http.ResponseWriter, r *http.Request) {
	ownerId := r.URL.Query().Get("user")
	if ownerId == "" {
		return
	}
	ctx := hctx.Tenant.New(ownerId)
	c.webSocket.HandleWebSocket(ctx, w, r)
}

func (c *chatGroup) HandleStaticFiles(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("."))
	fs.ServeHTTP(w, r)
}

func (c *chatGroup) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	ownerId := r.URL.Query().Get("user")
	if ownerId == "" {
		return
	}

	toPhone := r.Header.Get("to-phone")
	if toPhone == "" {
		return
	}
	ctx := hctx.Tenant.New(ownerId)

	message, err := c.whatsAppService.GetAllMessage(ctx, ownerId, toPhone)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	resp := make([]dto.Message, 0)
	for _, v := range message {
		resp = append(resp, dto.Message{
			Username: v.ProfileName,
			Message:  v.Message,
			IsOwner:  v.IsOwner,
			ToPhone:  v.ToPhone,
		})
	}
	marshal, err := json.Marshal(resp)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
}

func (c *chatGroup) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	ownerId := r.URL.Query().Get("user")
	if ownerId == "" {
		return
	}

	ctx := hctx.Tenant.New(ownerId)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	var message dto.Message
	if err = json.Unmarshal(body, &message); err != nil {
		return
	}
	c.whatsAppService.ReceiveMessage(ctx, &message)

	return
}

func (c *chatGroup) StartTemplate(w http.ResponseWriter, r *http.Request) {
	ownerId := r.URL.Query().Get("user")
	if ownerId == "" {
		return
	}

	ctx := hctx.Tenant.New(ownerId)

	toNumber := r.Header.Get("to-number")
	if toNumber == "" {
		return
	}
	template := dto.Template{
		OwnerId:  ownerId,
		ToPhone:  toNumber,
		Name:     "hello_world",
		Language: "en_US",
	}

	if err := c.whatsAppService.StartTemplate(ctx, ownerId, &template); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
