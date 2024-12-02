package entity

import (
	"context"

	"github.com/inter-hubly/pulse/internal/domain/valueobject"
)

type ConversationType string

const (
	ConversationTypeWhatsApp ConversationType = "whatsApp"
)

type Conversation interface {
	GetConversation(ctx context.Context) []*valueobject.Message
	PushConversation(ctx context.Context, v []*valueobject.Message)
	PushMessage(ctx context.Context, v *valueobject.Message)
}

type conversation struct {
	id               string                 `json:"id"`
	conversationType ConversationType       `json:"type"`
	messages         []*valueobject.Message `json:"messages"`
}

func NewConversation(ctx context.Context, id string, convType ConversationType) *conversation {
	return &conversation{
		id:               id,
		conversationType: convType,
	}
}

// PushConversation need get all conversation
func (c *conversation) PushConversation(ctx context.Context, v []*valueobject.Message) {
	c.messages = append(c.messages, v...)
}

// GetConversation get the current conversation
func (c *conversation) GetConversation(ctx context.Context) []*valueobject.Message {
	return c.messages
}

func (c *conversation) PushMessage(ctx context.Context, v *valueobject.Message) {
	c.messages = append(c.messages, v)
}
