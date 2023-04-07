package config

import (
	"fmt"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/relay"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/relay/producer/kafka"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/repo"
	"time"
)

var _ relay.Config = &Relay{}

type Relay struct {
	Brokers     []string `yaml:"brokers"`
	ChannelSize uint64   `yaml:"channelSize"`

	ConsumerCount  uint64        `yaml:"consumerCount"`
	ConsumeSize    uint64        `yaml:"consumeSize"`
	ConsumeTimeout time.Duration `yaml:"consumeTimeout"`

	ProducerCount uint64 `yaml:"producerCount"`
	WorkerCount   int    `yaml:"workerCount"`

	Repo   repo.EventRepo
	Sender kafka.EventSender
}

func (r *Relay) GetBrokers() []string {
	return r.Brokers
}

func (r *Relay) GetChannelSize() uint64 {
	return r.ChannelSize
}

func (r *Relay) GetConsumerCount() uint64 {
	return r.ConsumerCount
}

func (r *Relay) GetConsumeSize() uint64 {
	return r.ConsumeSize
}

func (r *Relay) GetConsumeTimeout() time.Duration {
	return r.ConsumeTimeout
}

func (r *Relay) GetProducerCount() uint64 {
	return r.ProducerCount
}

func (r *Relay) GetWorkerCount() int {
	return r.WorkerCount
}

func (r *Relay) SetEventRepo(repo repo.EventRepo) {
	r.Repo = repo
	return
}

func (r *Relay) SetEventSender(sender kafka.EventSender) {
	r.Sender = sender
	return
}

func (r *Relay) GetEventRepo() (repo.EventRepo, error) {
	if r.Repo == nil {
		return nil, fmt.Errorf("event repo is not set")
	}
	return r.Repo, nil
}

func (r *Relay) GetEventSender() (kafka.EventSender, error) {
	if r.Sender == nil {
		return nil, fmt.Errorf("event sender is not set")
	}
	return r.Sender, nil
}
