package engine

type Engine interface {
	Execute(t Task)
}

type Task interface {
	Execute()
}

type Work[E Engine, T any] func(e E) (T, error)

// Wrapper for work response
type ResponseWrapper[T any] struct {
	data T
	err  error
}

// Wrapper for work
type Wrapper[E Engine, T any] struct {
	response chan ResponseWrapper[T]
	engine   E
	work     Work[E, T]
}

// Execute the work
func (wrapper Wrapper[E, T]) Execute() {
	data, err := wrapper.work(wrapper.engine)
	wrapper.response <- ResponseWrapper[T]{data, err}
}

// Wait to receive the response
func (wrapper *Wrapper[E, T]) Wait() (T, error) {
	res := <-wrapper.response
	return res.data, res.err
}
