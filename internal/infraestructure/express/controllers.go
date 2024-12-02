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
	chatGroupController controller.ChatGroup
}

func NewPulseControllers() {
	controllersOnce.Do(func() {
		chatGroupControllers = &controllers{
			chatGroupController: controller.NewChatGroup(),
		}
	})
	chatGroupControllers.startControllers()
}

func (c *controllers) startControllers() {
	http.HandleFunc("/", c.chatGroupController.HandleStaticFiles)
	http.HandleFunc("/ws", c.chatGroupController.Handle)
	http.HandleFunc("/receive", c.chatGroupController.ReceiveMessage)
	http.HandleFunc("/messages", c.chatGroupController.GetAllMessages)
}
