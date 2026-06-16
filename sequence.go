package main

// Seq[T,U] is used to keep values for parser sequences built by Parser[T].Keep[U].
// It is expected that the first parser in the sequence will keep its value.
type Seq[T any, U any] struct {
	first  T
	second U
}

func Keep[T any, U any](parserT Parser[T], parserU Parser[U]) Parser[Seq[T, U]] {
	return func(initial state) (Seq[T, U], state, error) {
		t, next, err := parserT(initial)
		if err != nil {
			var zero Seq[T, U]
			return zero, initial, err
		}
		u, final, err := parserU(next)
		if err != nil {
			var zero Seq[T, U]
			return zero, initial, err
		}
		return Seq[T, U]{first: t, second: u}, final, nil
	}
}

func Skip[T any, U any](parserT Parser[T], parserU Parser[U]) Parser[T] {
	return func(initial state) (T, state, error) {
		t, next, err := parserT(initial)
		if err != nil {
			var zero T
			return zero, initial, err
		}
		_, final, err := parserU(next)
		if err != nil {
			var zero T
			return zero, initial, err
		}
		return t, final, nil
	}
}

func StartKeeping[T any](parser Parser[T]) Parser[Seq[Empty, T]] {
	return Map(parser, func(t T) Seq[Empty, T] {
		return Seq[Empty, T]{first: Empty{}, second: t}
	})
}

func StartSkipping[T any](parser Parser[T]) Parser[Empty] {
	return Map(parser, func(T) Empty { return Empty{} })
}

func Apply[T any, A any](parser Parser[Seq[Empty, T]], mapper func(T) A) Parser[A] {
	return func(initial state) (A, state, error) {
		seq, next, err := parser(initial)
		if err != nil {
			var zero A
			return zero, initial, err
		}
		return mapper(seq.second), next, nil
	}
}

func Apply2[T any, U any, A any](parser Parser[Seq[Seq[Empty, T], U]], mapper func(T, U) A) Parser[A] {
	return func(initial state) (A, state, error) {
		seq, next, err := parser(initial)
		if err != nil {
			var zero A
			return zero, initial, err
		}
		return mapper(seq.first.second, seq.second), next, nil
	}
}

func Apply3[T any, U any, V any, A any](parser Parser[Seq[Seq[Seq[Empty, T], U], V]], mapper func(T, U, V) A) Parser[A] {
	return func(initial state) (A, state, error) {
		seq, next, err := parser(initial)
		if err != nil {
			var zero A
			return zero, initial, err
		}
		return mapper(seq.first.first.second, seq.first.second, seq.second), next, nil
	}
}
