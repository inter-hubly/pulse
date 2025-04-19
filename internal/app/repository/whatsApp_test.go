package repository

import (
	"testing"
)

// mockgen -source=/home/saimon/Documents/Hubly/pulse/internal/app/repository/whatsApp.go  -destination=/home/saimon/Documents/Hubly/pulse/internal/app/repository/whatsApp_mock.go -package=repository

func TestWhatsApp(t *testing.T) {
	// os.Setenv("ENVIRONMENT", "test")
	// server.MockStartEnv(context.Background(), "../../")
	// host := server.GetElasticSearch().Host
	//
	// elasticsearch.NewConn(elasticsearch.WithUrl([]string{host}))
	//
	// repository := NewWhatsApp()
	//
	// t.Run("Get all value", func(t *testing.T) {
	// 	ctx := context.Background()
	// 	message, err := repository.GetAllMessage(ctx, "515719138282305", "")
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, message)
	// })
}
