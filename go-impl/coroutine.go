package main

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// TaskType represents the type of task (coding or analyze)
type TaskType string

const (
	TaskTypeCoding  TaskType = "coding"
	TaskTypeAnalyze TaskType = "analyze"
)

// Task represents a unit of work
type Task struct {
	ID      string
	Type    TaskType
	Execute func(ctx context.Context) error
}

// TaskResult represents the result of a task execution
type TaskResult struct {
	TaskID string
	Error  error
}

// TaskRunner manages async task execution using goroutines and channels
type TaskRunner struct {
	ctx       context.Context
	cancel    context.CancelFunc
	taskChan  chan Task
	resultChan chan TaskResult
	wg        sync.WaitGroup
	concurrency int
}

// NewTaskRunner creates a new task runner
func NewTaskRunner(concurrency int) *TaskRunner {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskRunner{
		ctx:          ctx,
		cancel:       cancel,
		taskChan:     make(chan Task, concurrency*2),
		resultChan:   make(chan TaskResult, concurrency*2),
		concurrency: concurrency,
	}
}

// Context returns the runner's context
func (r *TaskRunner) Context() context.Context {
	return r.ctx
}

// Start starts the task runner with worker pool
func (r *TaskRunner) Start() {
	for i := 0; i < r.concurrency; i++ {
		r.wg.Add(1)
		go r.worker(i)
	}
}

// Stop stops the task runner and cancels all running tasks
func (r *TaskRunner) Stop() {
	r.cancel()
	close(r.taskChan)
	r.wg.Wait()
	close(r.resultChan)
}

// Submit submits a task for execution
func (r *TaskRunner) Submit(task Task) error {
	select {
	case r.taskChan <- task:
		return nil
	case <-r.ctx.Done():
		return r.ctx.Err()
	}
}

// ResultChan returns the channel for task results
func (r *TaskRunner) ResultChan() <-chan TaskResult {
	return r.resultChan
}

// worker is a single worker goroutine
func (r *TaskRunner) worker(id int) {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		case task, ok := <-r.taskChan:
			if !ok {
				return
			}
			r.executeTask(task)
		}
	}
}

// executeTask executes a single task
func (r *TaskRunner) executeTask(task Task) {
	err := task.Execute(r.ctx)
	result := TaskResult{
		TaskID: task.ID,
		Error:  err,
	}

	select {
	case r.resultChan <- result:
	case <-r.ctx.Done():
	}
}

// AsyncTask creates an async task from a function
func AsyncTask(id string, taskType TaskType, fn func(ctx context.Context) error) Task {
	return Task{
		ID:      id,
		Type:    taskType,
		Execute: fn,
	}
}

// GoAsync executes a function asynchronously using goroutine
func GoAsync(ctx context.Context, fn func() error) <-chan error {
	resultChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		resultChan <- fn()
	}()

	// Return a channel that also respects context cancellation
	go func() {
		select {
		case <-ctx.Done():
			// Context cancelled, but we can't stop the goroutine
			// The goroutine will complete on its own
		case <-resultChan:
		}
	}()

	return resultChan
}

// Pool manages a pool of goroutines for concurrent execution
type Pool struct {
	workers   int
	taskQueue chan func()
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewPool creates a new worker pool
func NewPool(workers int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		workers:   workers,
		taskQueue: make(chan func(), workers*2),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Context returns the pool's context
func (p *Pool) Context() context.Context {
	return p.ctx
}

// Start starts the worker pool
func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// Stop stops the worker pool
func (p *Pool) Stop() {
	p.cancel()
	close(p.taskQueue)
	p.wg.Wait()
}

// Submit submits a task to the pool
func (p *Pool) Submit(task func()) error {
	select {
	case p.taskQueue <- task:
		return nil
	case <-p.ctx.Done():
		return p.ctx.Err()
	}
}

// worker is a single pool worker
func (p *Pool) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.taskQueue:
			if !ok {
				return
			}
			task()
		}
	}
}

// Mutex provides a simple mutex wrapper for goroutine safety
type Mutex struct {
	mu sync.Mutex
}

// Lock acquires the mutex
func (m *Mutex) Lock() {
	m.mu.Lock()
}

// Unlock releases the mutex
func (m *Mutex) Unlock() {
	m.mu.Unlock()
}

// TryLock attempts to acquire the mutex without blocking
func (m *Mutex) TryLock() bool {
	return m.mu.TryLock()
}

// WithLock executes fn with mutex locked
func (m *Mutex) WithLock(fn func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	fn()
}

// Once provides a way to execute something exactly once
type Once struct {
	once sync.Once
}

// Do executes fn only once
func (o *Once) Do(fn func()) {
	o.once.Do(fn)
}

// Sleep is a context-aware sleep that respects cancellation
func Sleep(ctx context.Context, duration time.Duration) error {
	select {
	case <-time.After(duration):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// RateLimiter limits the rate of operations
type RateLimiter struct {
	tokens    int64
	maxTokens int64
	refillRate time.Duration
	lastRefill int64
	mu        sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:    int64(maxTokens),
		maxTokens: int64(maxTokens),
		refillRate: refillRate,
		lastRefill: time.Now().UnixNano(),
	}
}

// Allow checks if an operation is allowed
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UnixNano()
	elapsed := now - r.lastRefill

	if elapsed >= r.refillRate.Nanoseconds() {
		refills := elapsed / r.refillRate.Nanoseconds()
		r.tokens = minInt64(r.maxTokens, r.tokens+refills)
		r.lastRefill = now
	}

	if r.tokens > 0 {
		r.tokens--
		return true
	}
	return false
}

// WaitForToken waits until a token is available
func (r *RateLimiter) WaitForToken(ctx context.Context) error {
	for {
		if r.Allow() {
			return nil
		}
		if err := Sleep(ctx, r.refillRate/2); err != nil {
			return err
		}
	}
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// AtomicCounter provides thread-safe counter
type AtomicCounter struct {
	value int64
}

// Add adds delta to the counter
func (c *AtomicCounter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

// Get gets the current value
func (c *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Reset resets the counter to zero
func (c *AtomicCounter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}