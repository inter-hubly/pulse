package mediator

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/inter-hubly/pilot/broker"
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
			OwnerNumberId: ownerId,
			To:            msg.ToId,
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
	conversation, err := m.GetConversation(ctx, ownerId, msg.ToId)
	if err != nil {
		return err
	}
	conversation.PushMessage(ctx, msg)

	return nil
}

func (m *whatsAppMediator) ReceiveMessage(ctx context.Context, ownerId string, msg dto.Message) error {
	conversation, err := m.GetConversation(ctx, ownerId, msg.ToId)
	if err != nil {
		return err
	}
	message := valueobject.NewMessage(msg.ToId, msg.Message, time.Now().Unix(), false)
	conversation.PushMessage(ctx, message)

	return nil
}

type senderMessage struct {
	SenderAndReceiver senderAndReceiverDto `json:"senderAndReceiver"`
	Message           string               `json:"message"`
}
type senderAndReceiverDto struct {
	OwnerNumberId string `json:"OwnerNumberId"`
	To            string `json:"to"`
}
