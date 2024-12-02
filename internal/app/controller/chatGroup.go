package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

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
	channel   chan dto.Message
	websocket server.WebSocket
}

func NewChatGroup() *chatGroup {
	chatGroupOnce.Do(func() {
		chatGroupController = &chatGroup{
			channel:   make(chan dto.Message),
			websocket: server.NewWebSocket(),
		}
	})
	return chatGroupController
}

func (c *chatGroup) Handle(w http.ResponseWriter, r *http.Request) {
	c.websocket.HandleWebSocket(w, r, c.channel)
}

func (c *chatGroup) HandleStaticFiles(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("."))
	fs.ServeHTTP(w, r)
}

func (c *chatGroup) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	marshal, err := json.Marshal([]dto.Message{{Username: "test", Message: "Its ok"}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
	w.Header().Add("Content-Type", "application/json")
}

func (c *chatGroup) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	go func() {
		var message dto.Message
		if err = json.Unmarshal(body, &message); err != nil {
			http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
			return
		}
		log.Print("Mandando mensagem ", message)
		chatGroupController.channel <- message
	}()

	return
}
