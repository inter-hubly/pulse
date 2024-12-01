package service

import (
	"context"
	"sync"

	"github.com/inter-hubly/pulse/dto"
	"github.com/inter-hubly/pulse/internal/mediator"
)

type WhatsApp interface {
	GetAllMessage(ctx context.Context) ([]dto.Message, error)
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppService
)

type whatsAppService struct {
	whatsAppMediator mediator.WhatsApp
}

func NewWhatsApp() *whatsAppService {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppService{
			whatsAppMediator: mediator.NewWhatsApp(),
		}
	})
	return whatsApp
}

func (s *whatsAppService) GetAllMessage(ctx context.Context) ([]dto.Message, error) {
	return s.whatsAppMediator.GetAllMessage(ctx, "515719138282305")
}
