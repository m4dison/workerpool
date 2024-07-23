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
	ctx         context.Context
	done        chan struct{}
}

func New(ctx context.Context, workerCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	pool := &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan func(ctx context.Context) error),
		cancel:      cancel,
		ctx:         ctx,
		done:        make(chan struct{}),
	}
	pool.start(ctx)
	return pool
}

func (p *WorkerPool) start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case task, ok := <-p.taskQueue:
					if !ok {
						return
					}
					if err := task(ctx); err != nil {
						p.errOnce.Do(func() {
							p.err = err
							p.cancel()
						})
						return
					}
				case <-p.done:
					return
				}
			}
		}()
	}
}

func (p *WorkerPool) Run(task func(ctx context.Context) error) {
	select {
	case p.taskQueue <- task:
	case <-p.ctx.Done():
	}
}

func (p *WorkerPool) Wait() error {
	close(p.taskQueue)
	p.wg.Wait()
	close(p.done)
	return p.err
}
