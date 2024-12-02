package repository

import (
	"context"
	"sync"

	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pulse/internal/domain/dto"
)

type WhatsApp interface {
	GetAllMessage(ctx context.Context, id string) ([]dto.Message, error)
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

func (w *whatsAppRepository) GetAllMessage(ctx context.Context, id string) ([]dto.Message, error) {
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
	sms := make([]dto.Message, 0)
	for _, v := range all.Hits.Hits {
		response := v.(map[string]interface{})["_source"].(map[string]interface{})
		var dto dto.Message
		if response["type"] == "template" {
			dto.Message = response["templateName"].(string)
		} else {
			dto.Message = response["message"].(string)
		}
		dto.Username = response["toPhoneId"].(string)
		sms = append(sms, dto)
	}
	return sms, nil
}
