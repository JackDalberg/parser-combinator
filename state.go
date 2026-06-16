package main

import "unicode/utf8"

// Internal state for the parser
type state struct {
	data   string
	offset int
}

func (s state) remaining() string {
	return s.data[s.offset:]
}

func (s state) consume(n int) state {
	s.offset += n
	return s
}

func (s state) nextRune() (rune, state) {
	r, w := utf8.DecodeRuneInString(s.remaining())
	return r, s.consume(w)
}
