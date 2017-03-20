package workerpool

import "sync"

//Task is a unit for work for the worker pool to execute
type Task interface {
	Execute()
}

//Pool is the worker pool type
type Pool struct {
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

/*
NewPool is the convenience function to get a new pool with a designated number of workers and
the buffer (backlog) specified
*/
func NewPool(size int, buffer uint64) *Pool {
	pool := &Pool{
		tasks: make(chan Task, buffer),
		kill:  make(chan struct{}),
	}
	pool.Resize(size)
	return pool
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.Execute()
		case <-p.kill:
			return
		}
	}
}

//Resize resizes the pool to the number of workers.
func (p *Pool) Resize(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}
	for p.size > n {
		p.size--
		p.kill <- struct{}{}
	}
}

//Close closes the worker pool, use this after you have sent a batch of tasks
func (p *Pool) Close() {
	close(p.tasks)
}

//Wait will block until the worker pool has finished it's tasks
func (p *Pool) Wait() {
	p.wg.Wait()
}

//Exec adds a task to be executed, either directly or queued depending on current work load.
func (p *Pool) Exec(task Task) {
	p.tasks <- task
}
