package util

type executionPool struct {
	executors chan executor
}

type executor struct{}

func NewExecutionPool(limit int) *executionPool {
	pool := executionPool{}
	pool.executors = make(chan executor, limit)
	for range limit {
		pool.executors <- executor{}
	}
	return &pool
}

func (pool *executionPool) Acquire() executor {
	return <-pool.executors
}

func (pool *executionPool) Release(p executor) {
	pool.executors <- p
}

func (pool *executionPool) Close() {
	close(pool.executors)
}
