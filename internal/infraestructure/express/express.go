package express

import (
	rabbitmq "github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/server"
)

var exchangeBroker = "linker"

func Start() {

	rabbitmq.NewRabbitMQ(exchangeBroker, "topic")

	elasticsearch.NewConn(
		elasticsearch.WithUrl([]string{server.GetElasticSearch().Host}),
		elasticsearch.WithUsernameAndPassword(
			server.GetElasticSearch().Username,
			server.GetElasticSearch().Password,
		),
	)

	if err := rabbitmq.GetConnection().
		QueueBind(
			rabbitmq.NewQueueBinding("whatsapp.send", "whatsapp.send", exchangeBroker),
		); err != nil {
		panic(err)
	}

	NewPulseControllers()
}
