package aggregation

import (
	"context"
	"errors"

	"github.com/inter-hubly/pulse/internal/domain/entity"
)

type MessageGroups interface {
	GetConversationById(ctx context.Context, toId string) (entity.Conversation, error)
	AddConversation(ctx context.Context, toId string, conversation entity.Conversation)
}

type messageGroups struct {
	ownerId       string
	conversations map[string]entity.Conversation
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
