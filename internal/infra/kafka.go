package infra

import (
	"github.com/segmentio/kafka-go"
	"github.com/test/loan-service/internal/consts"
	"go.uber.org/dig"
	"time"
)

// KafkaCfg menyimpan konfigurasi Kafka
type KafkaCfg struct {
	BrokerAddress   string        `envconfig:"BROKER_ADDRESS" required:"true" default:"localhost:9092"`
	Topic           string        `envconfig:"TOPIC" required:"true" default:"loan-topic"`
	ConsumerGroup   string        `envconfig:"CONSUMER_GROUP" required:"true" default:"loan-consumer-group"`
	ProducerRetries int           `envconfig:"PRODUCER_RETRIES" default:"3"`
	Timeout         time.Duration `envconfig:"TIMEOUT" default:"30s"`
}

// KafkaClients menyimpan klien Kafka yang digunakan untuk Producer dan Consumer
type KafkaClients struct {
	dig.Out
	Producer *kafka.Writer
	Consumer *kafka.Reader
}

// KafkaCfgs adalah struktur untuk menerima konfigurasi Kafka
type KafkaCfgs struct {
	dig.In
	Kafka *KafkaCfg
}

// NewKafkaClients menginisialisasi Kafka clients (Producer dan Consumer)
func NewKafkaClients(cfgs KafkaCfgs) KafkaClients {
	// Inisialisasi Kafka Producer
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{cfgs.Kafka.BrokerAddress},
		Topic:        "",
		MaxAttempts:  3,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: cfgs.Kafka.Timeout,
	})

	// Inisialisasi Kafka Consumer
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfgs.Kafka.BrokerAddress},
		GroupTopics: []string{string(consts.ApprovalLoanTopic), string(consts.LoanDisburseTopic), string(consts.FundingProcessTopic)},
		GroupID:     cfgs.Kafka.ConsumerGroup,
		StartOffset: kafka.FirstOffset,
		// Mengatur timeout sesuai dengan kebutuhan
		SessionTimeout: cfgs.Kafka.Timeout,
	})

	// Mengembalikan klien Kafka (Producer dan Consumer)
	return KafkaClients{
		Producer: producer,
		Consumer: consumer,
	}
}
