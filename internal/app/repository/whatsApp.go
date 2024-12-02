package repository

import (
	"context"
	"sync"
	"time"

	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pulse/internal/domain/valueobject"
)

type WhatsApp interface {
	GetAllMessage(ctx context.Context, id string) ([]valueobject.Message, error)
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

func (w *whatsAppRepository) GetAllMessage(ctx context.Context, id string) ([]valueobject.Message, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"ownerId": id,
						},
					},
				},
			},
		},
	}

	all, err := w.elastic.FindAll(ctx, elasticIndex, query)
	if err != nil {
		return nil, err
	}
	sms := make([]valueobject.Message, 0)
	for _, v := range all.Hits.Hits {
		response := v.(map[string]interface{})["_source"].(map[string]interface{})
		var dto *valueobject.Message
		var sender string
		if s, ok := response["sender"].(string); ok {
			sender = s
		}
		if response["type"] == "template" {
			dto = valueobject.NewMessage(sender, response["templateName"].(string), time.Now().Unix())
		} else {
			dto = valueobject.NewMessage(sender, response["message"].(string), time.Now().Unix())
		}
		sms = append(sms, *dto)
	}
	return sms, nil
}
