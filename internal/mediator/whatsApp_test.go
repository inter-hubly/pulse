package mediator

import (
	"context"
	"os"
	"testing"

	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/server"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppMediator(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	server.MockStartEnv("../../")
	host := server.GetElasticSearch().Host

	elasticsearch.NewConn(elasticsearch.WithUrl([]string{host}))

	mediator := NewWhatsApp()

	t.Run("need to get all values", func(t *testing.T) {
		ctx := context.Background()
		message, err := mediator.GetAllMessage(ctx, "515719138282305")
		assert.Nil(t, err)
		assert.NotNil(t, message)

		allMessage, err := mediator.whatsAppCache.GetAllMessage(ctx, "515719138282305")
		assert.Nil(t, err)
		assert.NotNil(t, allMessage)
	})
}
