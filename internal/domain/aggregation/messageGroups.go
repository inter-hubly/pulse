package aggregation

import (
	"context"

	"github.com/inter-hubly/pulse/internal/domain/entity"
)

type MessageGroups interface {
	GetConversations(ctx context.Context) map[string][]entity.Conversation
	AddConversation(ctx context.Context, id string, conversation entity.Conversation)
}

type messageGroups struct {
	conversations map[string][]entity.Conversation
}

func NewMessageGroups() *messageGroups {
	return &messageGroups{
		conversations: make(map[string][]entity.Conversation),
	}
}

func (w *messageGroups) GetConversations(ctx context.Context) map[string][]entity.Conversation {
	return w.conversations
}

func (w *messageGroups) AddConversation(ctx context.Context, id string, conversation entity.Conversation) {
	if _, ok := w.conversations[id]; ok {
		w.conversations[id] = append(w.conversations[id], conversation)
	}
	w.conversations[id] = []entity.Conversation{conversation}
}
