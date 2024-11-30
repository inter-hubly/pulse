package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/server"
	"github.com/inter-hubly/pulse/dto"
	"github.com/inter-hubly/pulse/internal/repository"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan dto.Message)

func main() {
	server.FillConfigEnvironment()

	elasticsearch.NewConn(
		elasticsearch.WithUrl([]string{server.GetElasticSearch().Host}),
		elasticsearch.WithUsernameAndPassword(
			server.GetElasticSearch().Username,
			server.GetElasticSearch().Password,
		),
	)

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	log.Println("HTTP server started on :8082")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		app := repository.NewWhatsApp()
		message, err := app.GetMessage()
		if err != nil {
			log.Println(err)
		}

		time.Sleep(1 * time.Second)
		for _, msg := range message {
			broadcast <- msg

		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
