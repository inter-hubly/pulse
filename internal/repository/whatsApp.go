package repository

//
// import (
// 	"sync"
//
// 	"github.com/inter-hubly/pilot/database/elasticsearch"
// 	"github.com/inter-hubly/webhook/dto"
// )
//
// type WhatsApp interface {
// 	PersistMessage(ctx context.Context, message *entity.Chat) (string, error)
// 	SetStatusMessageById(ctx context.Context, messageId string, status dto.MessageStatus) error
// }
//
// var (
// 	whatsAppOnce sync.Once
// 	whatsApp     *whatsAppRepository
// 	elasticIndex = "whatsapp.ready"
// )
//
// type whatsAppRepository struct {
// 	elastic elasticsearch.ElasticConn
// }
//
// func NewWhatsApp() *whatsAppRepository {
// 	whatsAppOnce.Do(func() {
// 		whatsApp = &whatsAppRepository{
// 			elastic: elasticsearch.GetConnection(),
// 		}
// 	})
// 	return whatsApp
// }
