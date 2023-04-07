package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
)

const (
	userEventsTopic = "user_events"
)

type EventSender interface {
	Send(userEvent *model.UserEvent) error
}

type sender struct {
	syncProd sarama.SyncProducer
}

func NewSender(brokers []string) (EventSender, error) {
	syncProd, err := NewSyncProducer(brokers)
	if err != nil {
		return nil, err
	}
	return &sender{
		syncProd: syncProd,
	}, nil
}

func NewSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)

	return producer, err
}

func (p *sender) Send(event *model.UserEvent) error {
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic:     userEventsTopic,
		Partition: -1,
		Value:     sarama.ByteEncoder(b),
	}
	_, _, err = p.syncProd.SendMessage(msg)
	return err
}
