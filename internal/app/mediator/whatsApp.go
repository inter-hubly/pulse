package mediator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pulse/internal/app/cache"
	"github.com/inter-hubly/pulse/internal/app/repository"
	"github.com/inter-hubly/pulse/internal/domain/dto"
	"github.com/inter-hubly/pulse/internal/domain/entity"
	"github.com/inter-hubly/pulse/internal/domain/valueobject"
)

type WhatsApp interface {
	GetConversation(ctx context.Context, ownerId, ToId string) (entity.Conversation, error)
	SendMessage(ctx context.Context, ownerId string, msg *valueobject.Message) error
	ReceiveMessage(ctx context.Context, ownerId string, msg dto.Message) error
	StartTemplate(ctx context.Context, ownerId string, msg *dto.Template) error
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppMediator
)

type whatsAppMediator struct {
	whatsAppCache      cache.WhatsApp
	whatsAppRepository repository.WhatsApp
	broker             broker.Connection
}

func NewWhatsApp() *whatsAppMediator {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppMediator{
			broker:             broker.GetConnection(),
			whatsAppCache:      cache.NewWhatsApp(),
			whatsAppRepository: repository.NewWhatsApp(),
		}
	})
	return whatsApp
}

func (m *whatsAppMediator) GetConversation(ctx context.Context, ownerId, ToId string) (entity.Conversation, error) {
	var conversation entity.Conversation
	var err error
	conversation, err = m.whatsAppCache.GetConversationById(ctx, ownerId, ToId)
	if err != nil {
		hlog.Debug("whatsAppMediator.GetConversation", "Get info in repository")
		msg, err := m.whatsAppRepository.GetAllMessage(ctx, ownerId, ToId)
		if err != nil {
			return nil, err
		}
		conversation = entity.NewConversation(ctx, ownerId, entity.ConversationTypeWhatsApp)
		conversation.PushConversation(ctx, msg)
		m.whatsAppCache.SaveConversation(ctx, ToId, conversation)
	}
	return conversation, nil
}

func (m *whatsAppMediator) SendMessage(ctx context.Context, ownerId string, msg *valueobject.Message) error {
	msg.IsOwner = true
	sendMessage := senderMessage{
		Message: msg.Message,
		SenderAndReceiver: senderAndReceiverDto{
			OwnerId: ownerId,
			To:      msg.ToPhone,
		},
	}

	marshal, err := json.Marshal(sendMessage)
	if err != nil {
		hlog.Debug("whatsAppMediator.SendMessage", fmt.Sprintf("Marshal json fail %s", err))
		return err
	}
	if err = m.broker.Publish("whatsapp.send", marshal); err != nil {
		hlog.Debug("whatsAppMediator.SendMessage", fmt.Sprintf("Publish message fail %s", err))
		return err
	}
	conversation, err := m.GetConversation(ctx, ownerId, msg.ToPhone)
	if err != nil {
		return err
	}
	conversation.PushMessage(ctx, msg)

	return nil
}

func (m *whatsAppMediator) ReceiveMessage(ctx context.Context, ownerId string, msg dto.Message) error {
	conversation, err := m.GetConversation(ctx, ownerId, msg.ToPhone)
	if err != nil {
		return err
	}
	message := valueobject.NewMessage(
		valueobject.WithToPhone(msg.ToPhone),
		valueobject.WithMessage(msg.Message),
		valueobject.WithIsOwner(false),
	)
	conversation.PushMessage(ctx, message)

	return nil
}
func (m *whatsAppMediator) StartTemplate(ctx context.Context, ownerId string, template *dto.Template) error {

	type senderDtoStruct struct {
		SenderAndReceiver struct {
			OwnerId string `json:"ownerId"`
			To      string `json:"to"`
		} `json:"senderAndReceiver"`
		Name     string `json:"name"`
		Language string `json:"language"`
	}

	senderDto := senderDtoStruct{}
	senderDto.SenderAndReceiver.OwnerId = ownerId
	senderDto.SenderAndReceiver.To = template.ToPhone
	senderDto.Name = template.Name
	senderDto.Language = template.Language

	marshal, err := json.Marshal(&senderDto)
	if err != nil {
		log.Printf("Error marshalling template: %v", err)
	}

	if err = m.broker.Publish("whatsapp.start", marshal); err != nil {
		hlog.Error("whatsAppService.SendMessage", fmt.Sprintf("error sending template :%s", err))
	}
	conversation, err := m.GetConversation(ctx, ownerId, template.ToPhone)
	if err != nil {
		return err
	}
	message := valueobject.NewMessage(
		valueobject.WithToPhone(template.ToPhone),
		valueobject.WithProfileName(template.Name),
		valueobject.WithIsOwner(true),
	)
	conversation.PushMessage(ctx, message)
	return nil
}

func (m *whatsAppMediator) Notify(ctx context.Context, msg valueobject.Message) error {
	tenantId := hctx.Tenant.Get(ctx)
	return m.SendMessage(ctx, tenantId, &msg)
}

type senderMessage struct {
	SenderAndReceiver senderAndReceiverDto `json:"senderAndReceiver"`
	Message           string               `json:"message"`
}
type senderAndReceiverDto struct {
	OwnerId string `json:"OwnerId"`
	To      string `json:"to"`
}
