package main

import (
	"errors"
	"strings"
)

type Parser[T any] func(state) (T, state, error)

type Empty struct{}

var (
	ErrNoMatch         = errors.New("no match")
	ErrUnconsumedInput = errors.New("unconsumed input")
)

func (p Parser[T]) Parse(data string) (T, error) {
	initial := state{data: data, offset: 0}
	result, final, err := p(initial)
	if err != nil {
		var zero T
		return zero, err
	}
	if final.offset < len(final.data) {
		var zero T
		return zero, ErrUnconsumedInput
	}
	return result, err
}

// Fail[T] is a parser which always fails.
func Fail[T any](initial state) (T, state, error) {
	var zero T
	return zero, initial, ErrNoMatch
}

// Succeed[T] generates a Parser[T] which always succeeds and returns value.
func Succeed[T any](value T) Parser[T] {
	return func(initial state) (T, state, error) {
		return value, initial, nil
	}
}

func Map[T any, A any](parser Parser[T], mapper func(T) A) Parser[A] {
	return func(initial state) (A, state, error) {
		t, next, err := parser(initial)
		if err != nil {
			var zero A
			return zero, initial, err
		}
		return mapper(t), next, nil
	}
}

func AndThen[T any, U any](parser Parser[T], handler func(T) Parser[U]) Parser[U] {
	return func(initial state) (U, state, error) {
		t, next, err := parser(initial)
		if err != nil {
			var zero U
			return zero, initial, err
		}
		nextParser := handler(t)
		return nextParser(next)
	}
}

// OneOf[T] returns a Parser[T] which will try each parser in parsers in order.
// The value of the first successful parser is returned.
// If no parser succeeds, ErrNoMatch is returned.
func OneOf[T any](parsers ...Parser[T]) Parser[T] {
	return func(initial state) (T, state, error) {
		err := ErrNoMatch
		for _, parser := range parsers {
			result, next, err := parser(initial)
			if err == nil {
				return result, next, nil
			}
		}
		var zero T
		return zero, initial, err
	}
}

func ConsumeIf(condition func(rune) bool) Parser[Empty] {
	return func(initial state) (Empty, state, error) {
		r, next := initial.nextRune()
		if !condition(r) {
			return Empty{}, initial, ErrNoMatch
		}
		return Empty{}, next, nil
	}
}

func ConsumeWhile(condition func(rune) bool) Parser[Empty] {
	return func(initial state) (Empty, state, error) {
		current := initial
		for {
			r, next := current.nextRune()
			if !condition(r) {
				return Empty{}, current, nil
			}
			current = next
		}
	}
}

func ConsumeSome(condition func(rune) bool) Parser[Empty] {
	s := ConsumeIf(condition)
	return Skip(s, ConsumeWhile(condition))
}

func Exactly(token string) Parser[Empty] {
	return func(initial state) (Empty, state, error) {
		if strings.HasPrefix(initial.remaining(), token) {
			next := initial.consume(len(token))
			return Empty{}, next, nil
		}
		return Empty{}, initial, ErrNoMatch
	}
}

func (p Parser[T]) GetString() Parser[string] {
	return func(initial state) (string, state, error) {
		start := initial.offset
		_, next, err := p(initial)
		if err != nil {
			return "", initial, err
		}
		end := next.offset
		return next.data[start:end], next, nil
	}
}
