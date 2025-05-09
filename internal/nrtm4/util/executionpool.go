package util

type ExecutionPool struct {
	executors chan executor
}

type executor struct{}

func NewExecutionPool(limit int) *ExecutionPool {
	pool := ExecutionPool{}
	pool.executors = make(chan executor, limit)
	for range limit {
		pool.executors <- executor{}
	}
	return &pool
}

func (pool *ExecutionPool) Acquire() executor {
	return <-pool.executors
}

func (pool *ExecutionPool) Release(p executor) {
	pool.executors <- p
}

func (pool *ExecutionPool) Close() {
	close(pool.executors)
}
