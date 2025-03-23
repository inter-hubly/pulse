package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
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
	ctx := context.Background()
	ctx, err := ws.getUserId(r)
	if err != nil {
		hlog.Error(ctx, "webSocket.StartTemplate", fmt.Sprintf("Error when get user id: %s", err.Error()))
		http.Error(w, "error get user id", http.StatusInternalServerError)
		return
	}

	if err := ws.webSocketService.CreateConnection(ctx, w, r); err != nil {
		hlog.Error(ctx, "webSocket.Handle", fmt.Sprintf("error create connection: %s", err.Error()))
		http.Error(w, "error get user id", http.StatusInternalServerError)
		return
	}

	go ws.webSocketService.ReceiveMessageFromClient(ctx)
}

func (ws *webSocket) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	ctx, err := ws.getUserId(r)
	if err != nil {
		hlog.Error(ctx, "webSocket.ReceiveMessage", fmt.Sprintf("Error when get user id: %s", err.Error()))
		http.Error(w, "error get user id", http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		hlog.Error(context.Background(), "webSocket.ReceiveMessage", fmt.Sprintf("Error when unmarshal message %s : %s", err.Error(), body))
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var message dto.Message
	if err = json.Unmarshal(body, &message); err != nil {
		hlog.Error(context.Background(), "webSocket.ReceiveMessage", fmt.Sprintf("Error when unmarshal message %s : %s", err.Error(), body))
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}

	if err = ws.webSocketService.SendMessageToClient(ctx, &message); err != nil {
		hlog.Error(context.Background(), "webSocket.ReceiveMessage", fmt.Sprintf("Error sending message %s", err.Error()))
		http.Error(w, "error send message to client", http.StatusInternalServerError)
		return
	}
}

func (ws *webSocket) StartTemplate(w http.ResponseWriter, r *http.Request) {
	ctx, err := ws.getUserId(r)
	if err != nil {
		hlog.Error(ctx, "webSocket.StartTemplate", fmt.Sprintf("Error when get user id: %s", err.Error()))
		http.Error(w, "error get user id", http.StatusInternalServerError)
		return
	}

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

func (ws *webSocket) getUserId(r *http.Request) (context.Context, error) {
	ownerId := r.URL.Query().Get("user")
	if ownerId == "" {
		hlog.Error(context.Background(), "webSocket.ReceiveMessage", "Error when get user")
		return nil, errors.New("error get userId")
	}
	hlog.Info(context.Background(), "webSocket.getUserId", fmt.Sprintf("User id is: %s", ownerId))
	return hctx.Tenant.New(ownerId), nil
}
