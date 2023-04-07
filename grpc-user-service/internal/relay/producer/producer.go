package producer

import (
	"context"
	"github.com/gammazero/workerpool"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/relay/producer/kafka"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/repo"
	"github.com/pkg/errors"
	"log"
	"sync"
	"time"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration

	sender kafka.EventSender
	repo   repo.EventRepo
	events <-chan model.UserEvent

	workerPool *workerpool.WorkerPool

	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewProducer(
	n uint64,
	timeout time.Duration,
	sender kafka.EventSender,
	repo repo.EventRepo,
	events <-chan model.UserEvent,
	workerPool *workerpool.WorkerPool) Producer {

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	return &producer{
		n:          n,
		timeout:    timeout,
		sender:     sender,
		repo:       repo,
		events:     events,
		workerPool: workerPool,
		wg:         wg,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (p *producer) Start() {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events:
					log.Println("producer received event:", event)
					if err := p.sender.Send(&event); err != nil {
						log.Println("error Send(&event):", err)
						p.workerPool.Submit(func() {
							err = p.repo.UnlockEvents(p.ctx, event.ID)
							if err != nil {
								log.Println(errors.Wrap(err, "repo.UnlockEvents()"))
							}

						})
					} else {
						p.workerPool.Submit(func() {
							err = p.repo.RemoveEvents(p.ctx, event.ID)
							if err != nil {
								log.Println(errors.Wrap(err, "repo.UnlockEvents()"))
							}
						})
					}
				case <-p.ctx.Done():
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	p.cancel()
	p.wg.Wait()
}
