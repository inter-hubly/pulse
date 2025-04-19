//go:build e2e

package controller

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pulse/internal/domain/dto"
)

func TestWebSocket(t *testing.T) {
	socket := NewWebSocket()
	http.HandleFunc("/ws", socket.Handle)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		if err := http.ListenAndServe(":8082", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	t.Run("need handle message in front", func(t *testing.T) {
		ctx := hctx.Tenant.New("559153210606318")

		for i := 0; i < 10000; i++ {
			if err := socket.webSocketService.SendMessageToClient(ctx, generateRandomMessage()); err != nil {
				t.Errorf("failed to send message %d: %v", i+1, err)
				time.Sleep(10 * time.Second)
			}
		}

	})
	wg.Wait()
}

func generateRandomMessage() *dto.Message {
	return &dto.Message{
		Message:     fmt.Sprintf("test-%d", rand.Intn(1000)),
		ProfileName: fmt.Sprintf("profile-%d", rand.Intn(100)),
		ToPhone:     fmt.Sprintf("phone-%d", rand.Intn(100)),
		IsOwner:     rand.Intn(2) == 1,
	}
}
