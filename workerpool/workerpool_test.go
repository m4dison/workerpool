package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"
)

func createTask(duration time.Duration, resultErr error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		select {
		case <-time.After(duration):
			return resultErr
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func TestWorkerPool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	workerCount := 3
	pool := New(ctx, workerCount)

	taskDurations := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		3 * time.Second,
	}

	expectedErr := errors.New("task error")

	for _, duration := range taskDurations {
		task := createTask(duration, expectedErr)
		pool.Run(task)
	}

	err := pool.Wait()

	if err != expectedErr {
		t.Errorf("expected error: %v, got: %v", expectedErr, err)
	}
}

func TestWorkerPoolWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	workerCount := 3
	pool := New(ctx, workerCount)

	taskDurations := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
	}

	for _, duration := range taskDurations {
		task := createTask(duration, nil)
		pool.Run(task)
	}

	err := pool.Wait()

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestWorkerPoolWithCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workerCount := 3
	pool := New(ctx, workerCount)

	taskDurations := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		5 * time.Second,
	}

	for _, duration := range taskDurations {
		task := createTask(duration, nil)
		pool.Run(task)
	}

	time.Sleep(1 * time.Second)
	cancel()

	err := pool.Wait()

	if err != context.Canceled {
		t.Errorf("expected error: %v, got: %v", context.Canceled, err)
	}
}
