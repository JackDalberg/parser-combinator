package main

import "strconv"

type Binding struct {
	Name  string
	Value BindingValue
}

type BindingValue interface {
	IsBindingValue()
}

type BindingInt int

func (BindingInt) IsBindingValue() {}

type BindingBool bool

func (BindingBool) IsBindingValue() {}

type ConfigParsers struct {
	trueParser          Parser[bool]
	falseParser         Parser[bool]
	boolParser          Parser[bool]
	intParser           Parser[int]
	valueParser         Parser[BindingValue]
	nameParser          Parser[string]
	bindingParser       Parser[Binding]
	whitespaceParser    Parser[Empty]
	bindingsParser      Parser[[]Binding]
	ConfigurationParser Parser[[]Binding]
}

func isAsciiLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func isDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlphaNum(r rune) bool {
	return isAsciiLetter(r) || isDecimalDigit(r)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

func NewConfigParser() ConfigParsers {
	var p ConfigParsers

	p.trueParser = Map(Exactly("true"), func(Empty) bool { return true })

	p.falseParser = Map(Exactly("false"), func(Empty) bool { return false })

	p.boolParser = OneOf(p.trueParser, p.falseParser)

	p.intParser = AndThen(
		ConsumeSome(isDecimalDigit).GetString(),
		func(digits string) Parser[int] {
			if len(digits) > 1 && digits[0] == '0' {
				return Fail[int]
			}
			v, err := strconv.Atoi(digits)
			if err != nil {
				return Fail[int]
			}
			return Succeed(v)
		},
	)

	p.valueParser = OneOf(
		Map(p.boolParser, func(v bool) BindingValue { return BindingBool(v) }),
		Map(p.intParser, func(i int) BindingValue { return BindingInt(i) }),
	)

	p.nameParser = AndThen(
		ConsumeIf(isAsciiLetter),
		func(Empty) Parser[Empty] { return ConsumeWhile(isAlphaNum) },
	).GetString()

	p.whitespaceParser = ConsumeWhile(isWhitespace)

	{ // bindingParser
		s1 := StartKeeping(p.nameParser)
		s2 := Skip(s1, p.whitespaceParser)
		s3 := Skip(s2, Exactly("="))
		s4 := Skip(s3, p.whitespaceParser)
		s5 := Keep(s4, p.valueParser)
		p.bindingParser = Apply2(s5, func(name string, value BindingValue) Binding {
			return Binding{Name: name, Value: value}
		})
	} // bindingParser

	type BindingNode struct {
		binding Binding
		next    *BindingNode
	}

	{ // bindingsParser
		p.bindingsParser = Loop(nil,
			func(bindNode *BindingNode) Parser[Step[*BindingNode, []Binding]] {
				if bindNode == nil {
					return Map(p.bindingParser, func(binding Binding) Step[*BindingNode, []Binding] {
						return Step[*BindingNode, []Binding]{
							Accum: &BindingNode{binding: binding},
							Done:  false,
						}
					})
				}
				s1 := StartSkipping(p.whitespaceParser)
				s2 := Skip(s1, Exactly(","))
				s3 := Skip(s2, p.whitespaceParser)
				s4 := Keep(s3, p.bindingParser)
				extend := Apply(s4, func(b Binding) Step[*BindingNode, []Binding] {
					return Step[*BindingNode, []Binding]{
						Accum: &BindingNode{binding: b, next: bindNode},
						Done:  false,
					}
				})
				var bindingSlice []Binding
				b := bindNode
				for b != nil {
					bindingSlice = append(bindingSlice, b.binding)
					b = b.next
				}
				return OneOf(
					extend,
					Succeed(Step[*BindingNode, []Binding]{Value: bindingSlice, Done: true}),
				)
			},
		)
	} // bindingsParser

	{ // ConfigurationParser
		s1 := StartSkipping(Exactly("["))
		s2 := Skip(s1, p.whitespaceParser)
		s3 := Keep(s2, p.bindingsParser)
		s4 := Skip(s3, p.whitespaceParser)
		s5 := Skip(s4, Exactly("]"))
		p.ConfigurationParser = Apply(s5, func(b []Binding) []Binding { return b })
	} // ConfigurationParser

	return p
}
