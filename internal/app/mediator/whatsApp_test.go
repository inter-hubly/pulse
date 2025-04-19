//go:build e2e

package mediator

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/server"
	"github.com/inter-hubly/pulse/internal/app/cache"
	"github.com/inter-hubly/pulse/internal/app/repository"
	"github.com/inter-hubly/pulse/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppMediator(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	server.MockStartEnv("../../../")

	ctx := context.Background()
	host := server.GetElasticSearch().Host
	elasticsearch.NewConn(elasticsearch.WithUrl([]string{host}))

	ctrl := gomock.NewController(t)
	repository := repository.NewMockWhatsApp(ctrl)

	repository.EXPECT().GetAllMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		[]*valueobject.Message{
			valueobject.NewMessage("123456", "first_ok", time.Now().Unix()),
			valueobject.NewMessage("123456", "second_ok", time.Now().Unix()),
			valueobject.NewMessage("123456", "third_ok", time.Now().Unix()),
			valueobject.NewMessage("515719138282305", "fourth_ok", time.Now().Unix()),
		}, nil,
	)

	mediator := &whatsAppMediator{
		whatsAppRepository: repository,
		whatsAppCache:      cache.NewWhatsApp(),
	}
	for _, v := range []struct {
		testName string
		useCache bool
	}{
		{
			testName: "Need get all message and no use cache",
		},
		{
			testName: "Need get all message and use cache",
		},
	} {
		t.Run(v.testName, func(t *testing.T) {
			message, err := mediator.GetConversation(ctx, "515719138282305", "123456")
			if v.useCache {
				message, err = mediator.GetConversation(ctx, "515719138282305", "123456")
			}

			assert.Nil(t, err)
			assert.NotNil(t, message)
		})
	}
}
