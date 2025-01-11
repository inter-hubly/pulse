package express

import (
	"net/http"
	"sync"

	"github.com/inter-hubly/pulse/internal/app/controller"
)

var (
	controllersOnce      sync.Once
	chatGroupControllers *controllers
)

type controllers struct {
	webSocketController controller.WebSocket
}

func NewPulseControllers() {
	controllersOnce.Do(func() {
		chatGroupControllers = &controllers{
			webSocketController: controller.NewWebSocket(),
		}
	})
	chatGroupControllers.startControllers()
}
func (c *controllers) startControllers() {
	http.HandleFunc("/ws", withCors(c.webSocketController.Handle))
	http.HandleFunc("/receive", withCors(c.webSocketController.ReceiveMessage))
	http.HandleFunc("/start-template", withCors(c.webSocketController.StartTemplate))
}

func withCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
