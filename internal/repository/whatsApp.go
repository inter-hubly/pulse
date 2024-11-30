package repository

import (
	"context"
	"log"
	"sync"

	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pulse/dto"
)

type WhatsApp interface {
	GetMessage() ([]dto.Message, error)
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppRepository
	elasticIndex = "whatsapp.ready"
)

type whatsAppRepository struct {
	elastic elasticsearch.ElasticConn
}

func NewWhatsApp() *whatsAppRepository {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppRepository{
			elastic: elasticsearch.GetConnection(),
		}
	})
	return whatsApp
}

func (w *whatsAppRepository) GetMessage() ([]dto.Message, error) {
	ctx := context.Background()
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	all, err := w.elastic.FindAll(ctx, "whatsapp.ready", query)
	if err != nil {
		return nil, err
	}
	sms := make([]dto.Message, 0)
	for _, v := range all.Hits.Hits {
		log.Print(v)
		sms = append(sms, dto.Message{
			Username: "ok",
		})
	}
	return sms, nil
}
