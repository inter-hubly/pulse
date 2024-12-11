package repository

import (
	"context"
	"sync"

	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pulse/internal/domain/valueobject"
)

type WhatsApp interface {
	GetAllMessage(ctx context.Context, id, toId string) ([]*valueobject.Message, error)
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

func (w *whatsAppRepository) GetAllMessage(ctx context.Context, id, toId string) ([]*valueobject.Message, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"ownerId": id,
						},
					},
					{
						"match": map[string]interface{}{
							"toPhone": toId,
						},
					},
				},
			},
		},
	}

	elasticResponse, err := w.elastic.FindAll(ctx, elasticIndex, query)
	if err != nil {
		return nil, err
	}
	sms := make([]*valueobject.Message, 0)
	if elasticResponse == nil {
		return []*valueobject.Message{}, nil
	}
	for _, v := range elasticResponse.Hits.Hits {
		response := v.(map[string]interface{})["_source"].(map[string]interface{})
		var dto *valueobject.Message
		var sender string
		var isOwner, ok bool

		if isOwner, ok = response["isOwner"].(bool); ok {
			var key string
			if isOwner {
				key = "ownerId"
			} else {
				key = "profileName"
			}

			if s, ok := response[key].(string); ok {
				sender = s
			} else {
				sender = response["toPhone"].(string)
			}
		}

		if response["type"] == "template" {
			dto = valueobject.NewMessage(
				valueobject.WithProfileName(sender),
				valueobject.WithMessage(response["templateName"].(string)),
				valueobject.WithIsOwner(isOwner),
				valueobject.WithToPhone(response["toPhone"].(string)),
			)
		} else {
			dto = valueobject.NewMessage(
				valueobject.WithProfileName(sender),
				valueobject.WithMessage(response["message"].(string)),
				valueobject.WithIsOwner(isOwner),
				valueobject.WithToPhone(response["toPhone"].(string)),
			)
		}
		sms = append(sms, dto)
	}
	return sms, nil
}
