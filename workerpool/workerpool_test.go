package workerpool

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestWorkerPool_Run_Success(t *testing.T) {
	ctx := context.Background()
	pool := New(ctx, 3)

	var counter int32
	task := func(ctx context.Context) error {
		atomic.AddInt32(&counter, 1)
		return nil
	}

	for i := 0; i < 10; i++ {
		pool.Run(task)
	}

	if err := pool.Wait(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if counter != 10 {
		t.Errorf("expected 10 tasks to run, but got %d", counter)
	}
}

func TestWorkerPool_Run_Empty(t *testing.T) {
	ctx := context.Background()
	pool := New(ctx, 3)

	err := pool.Wait()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorkerPool_Run_MultipleTasks(t *testing.T) {
	ctx := context.Background()
	pool := New(ctx, 3)

	var counter1, counter2 int32
	task1 := func(ctx context.Context) error {
		atomic.AddInt32(&counter1, 1)
		return nil
	}
	task2 := func(ctx context.Context) error {
		atomic.AddInt32(&counter2, 1)
		return nil
	}

	for i := 0; i < 5; i++ {
		pool.Run(task1)
		pool.Run(task2)
	}

	if err := pool.Wait(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if counter1 != 5 {
		t.Errorf("expected 5 task1 to run, but got %d", counter1)
	}
	if counter2 != 5 {
		t.Errorf("expected 5 task2 to run, but got %d", counter2)
	}
}
