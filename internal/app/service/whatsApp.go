package service

import (
	"context"
	"sync"

	"github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pulse/internal/app/mediator"
	"github.com/inter-hubly/pulse/internal/domain/dto"
)

type WhatsApp interface {
	StartTemplate(ctx context.Context, message *dto.Template) error
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppService
)

type whatsAppService struct {
	whatsAppMediator mediator.WhatsApp
	whatsAppBroker   broker.Connection
}

func NewWhatsApp() *whatsAppService {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppService{
			whatsAppMediator: mediator.NewWhatsApp(),
			whatsAppBroker:   broker.GetConnection(),
		}
	})
	return whatsApp
}

func (s *whatsAppService) StartTemplate(ctx context.Context, template *dto.Template) error {
	tenant := hctx.Tenant.Get(ctx)
	if err := s.whatsAppMediator.StartTemplate(ctx, tenant, template); err != nil {
		return err
	}

	return nil
}
