package repository

import (
	"context"
	"os"
	"testing"

	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/server"
	"github.com/stretchr/testify/assert"
)

func TestWhatsApp(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	server.MockStartEnv("../../")
	host := server.GetElasticSearch().Host

	elasticsearch.NewConn(elasticsearch.WithUrl([]string{host}))

	repository := NewWhatsApp()

	t.Run("Get all value", func(t *testing.T) {
		ctx := context.Background()
		message, err := repository.GetAllMessage(ctx, "515719138282305")
		assert.Nil(t, err)
		assert.NotNil(t, message)
	})

}
