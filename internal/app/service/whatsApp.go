package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pulse/internal/app/mediator"
	"github.com/inter-hubly/pulse/internal/domain/dto"
	"github.com/inter-hubly/pulse/internal/domain/valueobject"
	"github.com/inter-hubly/pulse/internal/infraestructure/server"
)

type WhatsApp interface {
	GetAllMessage(ctx context.Context, ownerId, toId string) ([]*valueobject.Message, error)
	SendMessage(ctx context.Context, message *dto.Message) error
	StartTemplate(ctx context.Context, ownerId string, message *dto.Template) error
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppService
)

type whatsAppService struct {
	whatsAppMediator mediator.WhatsApp
	whatsAppBroker   broker.Connection
	webSocketServer  server.WebSocket
}

func NewWhatsApp() *whatsAppService {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppService{
			whatsAppMediator: mediator.NewWhatsApp(),
			whatsAppBroker:   broker.GetConnection(),
			webSocketServer:  server.NewWebSocket(),
		}
	})
	return whatsApp
}

func (s *whatsAppService) GetAllMessage(ctx context.Context, ownerId, toId string) ([]*valueobject.Message, error) {
	conversation, err := s.whatsAppMediator.GetConversation(ctx, ownerId, toId)
	if err != nil {
		return nil, err
	}
	return conversation.GetConversation(ctx).GetMessages(ctx), nil
}

func (s *whatsAppService) SendMessage(ctx context.Context, message *dto.Message) error {
	type senderDtoStruct struct {
		SenderAndReceiver struct {
			OwnerNumberId string `json:"OwnerNumberId"`
			From          string `json:"from"`
			To            string `json:"to"`
		} `json:"senderAndReceiver"`
		Message string `json:"message"`
	}
	senderDto := senderDtoStruct{}
	senderDto.SenderAndReceiver.OwnerNumberId = "515719138282305"
	// senderDto.SenderAndReceiver.From = "15551817023"
	senderDto.SenderAndReceiver.To = "+5548991784586"
	senderDto.Message = message.Message
	marshal, err := json.Marshal(&senderDto)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
	}
	if err = s.whatsAppBroker.Publish("whatsapp.send", marshal); err != nil {
		hlog.Error("whatsAppService.SendMessage", fmt.Sprintf("error sending message :%s", err))
	}
	return nil
}

func (s *whatsAppService) StartTemplate(ctx context.Context, ownerId string, template *dto.Template) error {

	// if err := s.whatsAppMediator.StartTemplate(ctx, ownerId, template); err != nil {
	// 	return err
	// }

	connection, err := s.webSocketServer.GetConnection(ctx)
	if err != nil {
		return err
	}
	connection.WriteMessage(ctx, valueobject.Message{
		ToNumber: template.ToNumber,
		Message:  template.Name,
		IsOwner:  true,
	})
	return nil
}
