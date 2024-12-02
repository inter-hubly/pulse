package cache

import (
	"context"
	"errors"
	"sync"

	"github.com/inter-hubly/pulse/internal/domain/dto"
)

var messages map[string][]dto.Message

type WhatsApp interface {
	GetAllMessage(ctx context.Context, id string) ([]dto.Message, error)
	SaveMessage(ctx context.Context, id string, message dto.Message)
	SaveAllMessageInCache(ctx context.Context, id string, message []dto.Message)
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppCache
)

type whatsAppCache struct {
}

func NewWhatsApp() *whatsAppCache {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppCache{}
		messages = make(map[string][]dto.Message)
	})
	return whatsApp
}

func (w *whatsAppCache) GetAllMessage(ctx context.Context, id string) ([]dto.Message, error) {
	if value, ok := messages[id]; ok {
		return value, nil
	}
	return nil, errors.New("not found")
}

func (w *whatsAppCache) SaveMessage(ctx context.Context, id string, message dto.Message) {
	if value, ok := messages[id]; ok {
		value = append(value, message)
	}
	messages[id] = []dto.Message{message}
}

func (w *whatsAppCache) SaveAllMessageInCache(ctx context.Context, id string, message []dto.Message) {
	messages[id] = message
}
