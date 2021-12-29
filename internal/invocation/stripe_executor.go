package invocation

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
)

var (
	// Default values differ from java impl. Also queue size is calculated differently.
	// Java Client: queueSize per worker = defaultEventQueueCapacity / defaultEventWorkerCount
	// Go Client: queueSize per worker = defaultEventQueueCapacity
	defaultEventQueueCapacity = int32(10000)
	defaultEventWorkerCount   = int32(runtime.NumCPU())
)

// executor represents the function that will run on workers of stripeExecutor.
type executor func(queue chan func(), quit chan struct{}, wg *sync.WaitGroup)

// stripeExecutor executes given "tasks" preserving the order among the ones that are given with the same key.
type stripeExecutor struct {
	quit       chan struct{}
	wg         *sync.WaitGroup
	execFn     executor
	taskQueues []chan func()
	queueCount int32
}

// newStripeExecutor returns a new stripeExecutor with default configuration.
func newStripeExecutor() stripeExecutor {
	// ignore error, default values do not raise error.
	ex, _ := newStripeExecutorWithConf(defaultEventWorkerCount, defaultEventQueueCapacity)
	return ex
}

// newStripeExecutor returns a new stripeExecutor with configured queueCount and queueSize.
func newStripeExecutorWithConf(queueCount, queueSize int32) (stripeExecutor, error) {
	if queueCount <= 0 {
		return stripeExecutor{}, fmt.Errorf("queueCount must be greater than 0")
	}
	if queueSize <= 0 {
		return stripeExecutor{}, fmt.Errorf("queueSize must be greater than 0")
	}
	se := stripeExecutor{
		taskQueues: make([]chan func(), queueCount),
		queueCount: queueCount,
	}
	for i := range se.taskQueues {
		se.taskQueues[i] = make(chan func(), queueSize)
	}
	se.quit = make(chan struct{})
	se.wg = &sync.WaitGroup{}
	se.execFn = defaultExecFn
	return se, nil
}

// start fires up the workers for each queue.
func (se stripeExecutor) start() {
	se.wg.Add(int(se.queueCount))
	for i := range se.taskQueues {
		go se.execFn(se.taskQueues[i], se.quit, se.wg)
	}
}

// dispatch sends the handler "task" to one of the appropriate taskQueues, "tasks" with the same key end up on the same queue.
func (se stripeExecutor) dispatch(key int32, task func()) {
	if key < 0 {
		// dispatch random
		key = rand.Int31n(se.queueCount)
	}
	se.taskQueues[key%se.queueCount] <- task
}

// stop blocks until all workers are stopped.
func (se stripeExecutor) stop() {
	close(se.quit)
	se.wg.Wait()
}

func (se stripeExecutor) setExecutorFnc(f executor) {
	se.execFn = f
}

func defaultExecFn(queue chan func(), quit chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case task := <-queue:
			task()
		case <-quit:
			return
		}
	}
}
