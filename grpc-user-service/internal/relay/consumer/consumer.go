package consumer

import (
	"context"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/repo"
	"log"
	"sync"
	"time"
)

type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.UserEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.UserEvent) Consumer {

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	return &consumer{
		n:         n,
		events:    events,
		repo:      repo,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		ctx:       ctx,
		cancel:    cancel,
		wg:        wg,
	}
}

func (c *consumer) Start() {
	log.Printf("starting %d consumer(s), tick: %s\n", c.n, c.timeout)
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.LockEvents(c.ctx, c.batchSize)
					if err != nil {
						log.Println("LockEvents() error:", err)
						continue
					}
					for _, event := range events {
						c.events <- event
					}
				case <-c.ctx.Done():
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	c.cancel()
	c.wg.Wait()
}
