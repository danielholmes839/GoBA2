package engine

// Execute some work in an Engine
func Execute[E Engine, T any](engine E, work Work[E, T]) (T, error) {
	// create the wrapper
	wrapper := Wrapper[E, T]{
		response: make(chan ResponseWrapper[T], 1),
		engine:   engine,
		work:     work,
	}

	defer close(wrapper.response)
	engine.Execute(wrapper)
	return wrapper.Wait()
}
