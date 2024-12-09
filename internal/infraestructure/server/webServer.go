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
			conversations: make(map[string]WebSocketConversation),
		}
	})
	return serverGroup
}

type WebSocketServer interface {
	GetConnection(ctx context.Context) (WebSocketConversation, error)
	SaveConnection(ctx context.Context, connection WebSocketConversation) error
	DeleteConnection(ctx context.Context)
}

type conversationGroups struct {
	conversations map[string]WebSocketConversation
}

func (s *conversationGroups) GetConnection(ctx context.Context) (WebSocketConversation, error) {
	tenantId := hctx.Tenant.Get(ctx)
	// se existe a conexão aberta
	if conv, ok := s.conversations[tenantId]; ok {
		return conv, nil
	}
	return nil, errors.New("connection not found")
}

func (s *conversationGroups) SaveConnection(ctx context.Context, connection WebSocketConversation) error {
	tenantId := hctx.Tenant.Get(ctx)

	s.conversations[tenantId] = connection

	return nil
}

func (s *conversationGroups) DeleteConnection(ctx context.Context) {
	tenantId := hctx.Tenant.Get(ctx)
	delete(s.conversations, tenantId)
}
