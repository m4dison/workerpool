package workerpool

import (
	"context"
	"sync"
)

type WorkerPool struct {
	workerCount int
	taskQueue   chan func(ctx context.Context) error
	wg          sync.WaitGroup
	errOnce     sync.Once
	err         error
	cancel      context.CancelFunc
}

func New(ctx context.Context, workerCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	pool := &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan func(ctx context.Context) error),
		cancel:      cancel,
	}
	pool.start(ctx)
	return pool
}

func (p *WorkerPool) start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for task := range p.taskQueue {
				if err := task(ctx); err != nil {
					p.errOnce.Do(func() {
						p.err = err
						p.cancel()
					})
					break
				}
			}
		}()
	}
}

func (p *WorkerPool) Run(task func(ctx context.Context) error) {
	select {
	case p.taskQueue <- task:
	case <-context.Background().Done():
	}
}

func (p *WorkerPool) Wait() error {
	close(p.taskQueue)
	p.wg.Wait()
	return p.err
}
