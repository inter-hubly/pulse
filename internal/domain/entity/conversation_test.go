package entity

import (
	"context"
	"testing"
	"time"

	"github.com/inter-hubly/pulse/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestConversation(t *testing.T) {
	ctx := context.Background()
	id := "1234567"
	conversationEntity := NewConversation(ctx, id, ConversationTypeWhatsApp)
	t.Run("Conversation Test", func(t *testing.T) {
		messages := []*valueobject.Message{
			valueobject.NewMessage(true, "this is first message", time.Now().Unix()),
			valueobject.NewMessage(true, "this is second message", time.Now().Add(-2*time.Minute).Unix()),
		}
		conversationEntity.PushConversation(ctx, messages)

		getConversation := conversationEntity.GetConversation(ctx)
		assert.Equal(t, messages[0].GetMessage(), getConversation[0].GetMessage())
	})

	t.Run("Push Message", func(t *testing.T) {
		messages := []*valueobject.Message{
			valueobject.NewMessage(true, "this is first message", time.Now().Unix()),
			valueobject.NewMessage(true, "this is second message", time.Now().Add(-2*time.Minute).Unix()),
		}
		conversationEntity.PushConversation(ctx, messages)
		lastMessage := valueobject.NewMessage(true, "this is third message", time.Now().Unix())
		conversationEntity.PushMessage(ctx, lastMessage)

		allConv := conversationEntity.GetConversation(ctx)
		assert.Equal(t, allConv[2].GetMessage(), lastMessage.GetMessage())
	})
}
