package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pulse/internal/domain/dto"

	"github.com/inter-hubly/pulse/internal/app/service"
)

type WebSocket interface {
	Handle(w http.ResponseWriter, r *http.Request)
	ReceiveMessage(w http.ResponseWriter, r *http.Request)
	StartTemplate(w http.ResponseWriter, r *http.Request)
}

var (
	webSocketOnce       sync.Once
	webSocketController *webSocket
)

type webSocket struct {
	webSocketService service.WebSocket
	whatsAppService  service.WhatsApp
}

func NewWebSocket() *webSocket {
	webSocketOnce.Do(func() {
		webSocketController = &webSocket{
			webSocketService: service.NewWebSocket(),
			whatsAppService:  service.NewWhatsApp(),
		}
	})
	return webSocketController
}

func (ws *webSocket) Handle(w http.ResponseWriter, r *http.Request) {
	ownerId := r.URL.Query().Get("user")
	if ownerId == "" {
		return
	}
	ctx := hctx.Tenant.New(ownerId)

	if err := ws.webSocketService.CreateConnection(ctx, w, r); err != nil {
		return
	}

	go ws.webSocketService.ReceiveMessageFromClient(ctx)
}

func (ws *webSocket) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
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

	if err = ws.webSocketService.SendMessageToClient(ctx, &message); err != nil {
		return
	}
}

func (ws *webSocket) StartTemplate(w http.ResponseWriter, r *http.Request) {
	ownerId := r.URL.Query().Get("user")
	if ownerId == "" {
		return
	}

	ctx := hctx.Tenant.New(ownerId)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var message dto.Template
	if err = json.Unmarshal(body, &message); err != nil {
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}
	if err = ws.whatsAppService.StartTemplate(ctx, &message); err != nil {
		http.Error(w, "Error sending template", http.StatusInternalServerError)
		return
	}
}
