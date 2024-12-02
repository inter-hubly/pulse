package mediator

import (
	"context"
	"sync"

	"github.com/inter-hubly/pulse/internal/app/cache"
	"github.com/inter-hubly/pulse/internal/app/repository"
	"github.com/inter-hubly/pulse/internal/domain/entity"
)

type WhatsApp interface {
	GetAllMessage(ctx context.Context, id string) ([]entity.Conversation, error)
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppMediator
)

type whatsAppMediator struct {
	whatsAppCache      cache.WhatsApp
	whatsAppRepository repository.WhatsApp
}

func NewWhatsApp() *whatsAppMediator {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppMediator{
			whatsAppCache:      cache.NewWhatsApp(),
			whatsAppRepository: repository.NewWhatsApp(),
		}
	})
	return whatsApp
}

func (m *whatsAppMediator) GetAllMessage(ctx context.Context, id string) ([]entity.Conversation, error) {
	messages, err := m.whatsAppCache.GetAllMessage(ctx, id)
	if err == nil {
		return messages, nil
	}

	repositoryMessages, err := m.whatsAppRepository.GetAllMessage(ctx, id)
	for _, message := range repositoryMessages {
		entity.
			NewConversation(ctx, message.GetOwner(), entity.ConversationTypeWhatsApp).
			PushMessage(ctx, &message)
		sss
	}
	if err != nil {
		return nil, err
	}
	m.whatsAppCache.SaveAllMessageInCache(ctx, id, messages)
	return messages, nil
}
