package cache

import (
	"context"
	"errors"
	"sync"

	"github.com/inter-hubly/pulse/internal/domain/aggregation"
	"github.com/inter-hubly/pulse/internal/domain/entity"
)

type WhatsApp interface {
	GetAllMessage(ctx context.Context, id string) ([]entity.Conversation, error)
	SaveMessage(ctx context.Context, id string, message entity.Conversation)
	SaveAllMessageInCache(ctx context.Context, id string, message []entity.Conversation)
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppCache
)

type whatsAppCache struct {
	conversations aggregation.MessageGroups
}

func NewWhatsApp() *whatsAppCache {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppCache{
			conversations: aggregation.NewMessageGroups(),
		}
	})
	return whatsApp
}

func (w *whatsAppCache) GetAllMessage(ctx context.Context, id string) ([]entity.Conversation, error) {

	if value, ok := w.conversations.GetConversations(ctx)[id]; ok {
		return value, nil
	}
	return nil, errors.New("not found")
}

func (w *whatsAppCache) SaveMessage(ctx context.Context, id string, message entity.Conversation) {
	if value, ok := w.conversations.GetConversations(ctx)[id]; ok {
		value = append(value, message)
	}
	w.conversations.AddConversation(ctx, id, message)
}

func (w *whatsAppCache) SaveAllMessageInCache(ctx context.Context, id string, message []entity.Conversation) {
	w.conversations.GetConversations(ctx)[id] = message
}
