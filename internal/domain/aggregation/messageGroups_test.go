package aggregation

import (
	"context"
	"testing"
	"time"

	"github.com/inter-hubly/pulse/internal/domain/entity"
	"github.com/inter-hubly/pulse/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestMessageGroups(t *testing.T) {
	groups := NewMessageGroups()
	ctx := context.Background()
	t.Run("Create a message group", func(t *testing.T) {
		conversation := entity.NewConversation(ctx, "1", entity.ConversationTypeWhatsApp)

		conversation.PushMessage(ctx, valueobject.NewMessage(true, "mesasgeId", time.Now().Unix()))
		groups.AddConversation(ctx, "user1", conversation)

		conversations := groups.GetConversations(ctx)

		conv := conversations["user1"]
		assert.NotNil(t, conv)
		assert.Equal(t, conversation.GetConversation(ctx)[0].GetMessage(), conv[0].GetConversation(ctx)[0].GetMessage())
	})
}
