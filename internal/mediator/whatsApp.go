package mediator

import (
	"context"
	"sync"

	"github.com/inter-hubly/pulse/dto"
	"github.com/inter-hubly/pulse/internal/cache"
	"github.com/inter-hubly/pulse/internal/repository"
)

type WhatsApp interface {
	GetAllMessage(ctx context.Context, id string) ([]dto.Message, error)
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

func (m *whatsAppMediator) GetAllMessage(ctx context.Context, id string) (message []dto.Message, err error) {
	message, err = m.whatsAppCache.GetAllMessage(ctx, id)
	if err == nil {
		return message, nil
	}

	message, err = m.whatsAppRepository.GetAllMessage(ctx, id)
	if err != nil {
		return nil, err
	}
	m.whatsAppCache.SaveAllMessageInCache(ctx, id, message)
	return message, nil
}
