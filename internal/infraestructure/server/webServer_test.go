package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/inter-hubly/pilot/hctx"
)

func TestNewConversations(t *testing.T) {
	newServerConversation := NewServerConversation()
	for _, v := range []struct {
		tenantName string
		server     WebSocketServer
	}{
		{
			tenantName: "tenant1",
			server:     newServerConversation,
		},
		{
			tenantName: "tenant2",
			server:     newServerConversation,
		},
	} {
		t.Run(v.tenantName, func(t *testing.T) {
			ctx := hctx.Tenant.New(v.tenantName)

			server := httptest.NewServer(http.HandlerFunc(websocketHandler))
			defer server.Close()

			req := httptest.NewRequest("GET", "/test?toNumber=554899178", nil)

			rr := httptest.NewRecorder()
			NewWebSocket().HandleWebSocket(ctx, rr, req)

		})
	}

}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to websocket", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Aqui você pode adicionar lógica para manipular mensagens WebSocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
