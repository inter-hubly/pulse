package aggregation

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
	"github.com/inter-hubly/pulse/internal/domain/valueobject"

	"github.com/inter-hubly/pulse/internal/domain/entity"
)

type MessageGroups interface {
	GetConversationById(ctx context.Context, toId string) (entity.Conversation, error)
	AddConversation(ctx context.Context, toId string, conversation entity.Conversation)
}

type messageGroups struct {
	conversations map[string]entity.Conversation
	connection    *websocket.Conn
	ownerId       string
}

func NewMessageGroups(ownerId string) *messageGroups {
	return &messageGroups{
		ownerId:       ownerId,
		conversations: make(map[string]entity.Conversation),
	}
}

func (w *messageGroups) AddConversation(ctx context.Context, toId string, conversation entity.Conversation) {
	w.conversations[toId] = conversation
}

func (w *messageGroups) GetConversationById(ctx context.Context, toId string) (entity.Conversation, error) {
	if conversation, ok := w.conversations[toId]; ok {
		return conversation, nil
	}
	return nil, errors.New("conversation not found")
}

func (w *messageGroups) SendMessageToClient(ctx context.Context, msg valueobject.Message) error {
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return w.connection.WriteMessage(websocket.TextMessage, marshal)
}
