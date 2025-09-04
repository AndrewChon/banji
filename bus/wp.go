package bus

import (
	"sync"
)

// A workerPool manages a bounded number of workers that execute tasks in a concurrent fashion.
type workerPool struct {
	sema chan struct{}
	wg   sync.WaitGroup
}

func newWorkerPool(n int) *workerPool {
	return &workerPool{
		sema: make(chan struct{}, n),
	}
}

// post submits a task for execution.
func (p *workerPool) post(task func()) {
	p.wg.Add(1)
	p.sema <- struct{}{}

	go func() {
		task()
		<-p.sema
		p.wg.Done()
	}()
}

// wait blocks until all tasks posted to the workerPool have been completed.
func (p *workerPool) wait() {
	p.wg.Wait()
}
