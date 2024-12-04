package cache

import (
	"context"
	"errors"
	"sync"

	"github.com/inter-hubly/pulse/internal/domain/aggregation"
	"github.com/inter-hubly/pulse/internal/domain/entity"
)

type WhatsApp interface {
	GetConversationById(ctx context.Context, ownerId, toId string) (entity.Conversation, error)
	SaveConversation(ctx context.Context, toId string, message entity.Conversation)
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppCache
)

type whatsAppCache struct {
	cachedMessageGroups map[string]aggregation.MessageGroups
}

func NewWhatsApp() *whatsAppCache {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppCache{
			cachedMessageGroups: make(map[string]aggregation.MessageGroups),
		}

	})
	return whatsApp
}

func (w *whatsAppCache) SaveConversation(ctx context.Context, toId string, message entity.Conversation) {
	owner := message.GetId(ctx)
	groups := aggregation.NewMessageGroups(owner)
	groups.AddConversation(ctx, toId, message)
	w.cachedMessageGroups[owner] = groups
}

func (w *whatsAppCache) GetConversationById(ctx context.Context, ownerId, toId string) (entity.Conversation, error) {
	if messageGroup, ok := w.cachedMessageGroups[ownerId]; ok {
		return messageGroup.GetConversationById(ctx, toId)
	}
	return nil, errors.New("not found")
}
