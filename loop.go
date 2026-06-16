package main

type Step[A any, T any] struct {
	Done  bool
	Accum A
	Value T
}

func Loop[A any, T any](startAccum A, stepper func(A) Parser[Step[A, T]]) Parser[T] {
	return func(initial state) (T, state, error) {
		accum := startAccum
		currentState := initial
		for {
			parser := stepper(accum)
			step, nextState, err := parser(currentState)
			if err != nil {
				var zero T
				return zero, initial, err
			}
			if step.Done {
				return step.Value, nextState, nil
			}
			accum = step.Accum
			currentState = nextState
		}
	}
}
