package server

import (
	"context"
	"errors"
	"sync"

	"github.com/inter-hubly/pilot/hctx"
)

var (
	serverConversationOnce sync.Once
	serverGroup            *conversationGroups
)

func NewServerConversation() *conversationGroups {
	serverConversationOnce.Do(func() {
		serverGroup = &conversationGroups{
			conversations: make(map[string]map[string]WebSocketConversation),
		}
	})
	return serverGroup
}

type WebSocketServer interface {
	GetConnection(ctx context.Context, numberId string) (WebSocketConversation, error)
	SaveConnection(ctx context.Context, numberId string, connection WebSocketConversation) error
}

type conversationGroups struct {
	conversations map[string]map[string]WebSocketConversation
}

func (s *conversationGroups) GetConnection(ctx context.Context, numberId string) (WebSocketConversation, error) {
	tenantId := hctx.Tenant.Get(ctx)
	// se existe a conexão aberta
	if conv, ok := s.conversations[tenantId]; ok {

		// se a conexão com o número já está aberta
		if singleConv, convOk := conv[numberId]; convOk {
			return singleConv, nil
		}
	}
	return nil, errors.New("connection not found")
}

func (s *conversationGroups) SaveConnection(ctx context.Context, numberId string, connection WebSocketConversation) error {
	tenantId := hctx.Tenant.Get(ctx)
	if conv, ok := s.conversations[tenantId]; ok {
		conv[numberId] = connection
		return nil
	}
	s.conversations[tenantId] = map[string]WebSocketConversation{
		numberId: connection,
	}
	return nil
}
