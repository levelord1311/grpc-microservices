package relay

import (
	"github.com/gammazero/workerpool"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/relay/consumer"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/relay/producer"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/relay/producer/kafka"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/repo"
	"log"
	"time"
)

type Config interface {
	GetBrokers() []string
	GetChannelSize() uint64
	GetConsumerCount() uint64
	GetConsumeSize() uint64
	GetConsumeTimeout() time.Duration

	GetProducerCount() uint64
	GetWorkerCount() int

	GetEventRepo() (repo.EventRepo, error)
	GetEventSender() (kafka.EventSender, error)
}

type Relay struct {
	events     chan model.UserEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
}

func NewRelay(cfg Config) (*Relay, error) {
	events := make(chan model.UserEvent, cfg.GetChannelSize())
	workerPool := workerpool.New(cfg.GetWorkerCount())

	eventRepo, err := cfg.GetEventRepo()
	if err != nil {
		return nil, err
	}
	eventSender, err := cfg.GetEventSender()
	if err != nil {
		return nil, err
	}

	c := consumer.NewDbConsumer(
		cfg.GetConsumerCount(),
		cfg.GetConsumeSize(),
		cfg.GetConsumeTimeout(),
		eventRepo,
		events)

	p := producer.NewProducer(
		cfg.GetProducerCount(),
		cfg.GetConsumeTimeout(),
		eventSender,
		eventRepo,
		events,
		workerPool)

	return &Relay{
		events:     events,
		consumer:   c,
		producer:   p,
		workerPool: workerPool,
	}, nil

}

func (r *Relay) Start() {
	log.Println("relay started")

	r.producer.Start()
	r.consumer.Start()
}

func (r *Relay) Close() {
	r.consumer.Close()
	r.producer.Close()
	r.workerPool.StopWait()
}
