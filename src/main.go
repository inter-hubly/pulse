package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	rabbitmq "github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/server"
	"github.com/inter-hubly/pulse/dto"
	"github.com/inter-hubly/pulse/internal/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	broadcast = make(chan dto.Message)
	prevHash  string
)

func main() {
	server.FillConfigEnvironment()

	elasticsearch.NewConn(
		elasticsearch.WithUrl([]string{server.GetElasticSearch().Host}),
		elasticsearch.WithUsernameAndPassword(
			server.GetElasticSearch().Username,
			server.GetElasticSearch().Password,
		),
	)
	// Rota para arquivos estáticos
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	go getAllCon()
	log.Println("HTTP server started on :8082")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	// Protege acesso ao map de clientes
	clientsMu.Lock()
	clients[ws] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, ws)
		clientsMu.Unlock()
	}()

	for {
		var msg dto.Message

		err = ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			break
		}
		type senderDtoStruct struct {
			SenderAndReceiver struct {
				OwnerNumberId string `json:"OwnerNumberId"`
				From          string `json:"from"`
				To            string `json:"to"`
			} `json:"senderAndReceiver"`
			Message string `json:"message"`
		}
		senderDto := senderDtoStruct{}
		senderDto.SenderAndReceiver.OwnerNumberId = ""
		senderDto.SenderAndReceiver.From = ""
		senderDto.SenderAndReceiver.To = ""
		senderDto.Message = msg.Message
		marshal, err := json.Marshal(&senderDto)
		if err != nil {
			log.Printf("Error marshalling message: %v", err)
		}
		rabbitmq.GetConnection().Publish("whatsapp.message", marshal)

		// Envia a mensagem recebida para o canal de broadcast
		broadcast <- msg
	}
}

func getAllCon() {
	app := service.NewWhatsApp()

	for {
		time.Sleep(10 * time.Second)

		messages, err := app.GetAllMessage(context.Background())
		if err != nil {
			log.Printf("Error getting messages: %v", err)
			continue
		}

		// Verifica se a lista mudou
		if hasListChanged(messages) {
			for _, msg := range messages {
				log.Printf("Message received: %v", msg)
				broadcast <- msg
			}
		}
	}
}

func handleMessages() {
	for msg := range broadcast {
		clientsMu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
	}
}

func hasListChanged(messages []dto.Message) bool {
	// Serializa a lista de mensagens para calcular o hash
	data, err := json.Marshal(messages)
	if err != nil {
		log.Printf("Error serializing messages: %v", err)
		return false
	}

	// Calcula o hash atual
	currentHash := sha256.Sum256(data)
	currentHashStr := string(currentHash[:])

	// Compara com o hash anterior
	if currentHashStr != prevHash {
		prevHash = currentHashStr
		return true
	}

	return false
}
