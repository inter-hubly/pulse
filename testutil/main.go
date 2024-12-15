package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Configuração do WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Permitir qualquer origem
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Atualiza a conexão HTTP para WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao atualizar para WebSocket:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Cliente conectado via WebSocket")

	// Canal para fechar o loop de envio
	stop := make(chan bool)

	// Goroutine para enviar mensagens automáticas a cada 30 segundos
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				type Message struct {
					MessageId   string `json:"messageId"`
					ToPhone     string `json:"toPhone,omitempty"`
					ProfileName string `json:"profileName,omitempty"`
					Message     string `json:"message"`
					Time        int64  `json:"time"`
					IsOwner     bool   `json:"isOwner"`
				}
				rand.Seed(time.Now().UnixNano())
				i := rand.Int()

				var number string
				if i%2 == 0 {
					number = "+5548991784586"
				} else {
					number = "+5548996711701"
				}
				message := Message{
					MessageId:   uuid.New().String(),
					ToPhone:     number,
					ProfileName: "Saimon",
					Message:     fmt.Sprintf("O numero sorteado foi: %d", i),
				}

				marshal, err := json.Marshal(message)
				if err != nil {
					fmt.Println("Erro ao atualizar para JSON:", err)
				}

				if err = conn.WriteMessage(websocket.TextMessage, marshal); err != nil {
					fmt.Println("Erro ao enviar mensagem:", err)
					stop <- true
					return
				}
				fmt.Println("Mensagem enviada ao cliente:", message)
			case <-stop:
				fmt.Println("Parando envio de mensagens automáticas")
				return
			}
		}
	}()

	// Loop para receber mensagens do cliente
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Conexão encerrada ou erro ao ler mensagem:", err)
			stop <- true // Para o ticker ao desconectar
			break
		}
		fmt.Printf("Mensagem recebida do cliente: %s\n", string(message))
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/receive", receive)

	fmt.Println("Servidor WebSocket rodando na porta 8082...")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}

func receive(w http.ResponseWriter, r *http.Request) {
	
}
