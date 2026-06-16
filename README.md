# Generic Parser Combinator
This parser combinator was inspired by https://www.youtube.com/watch?v=x5p_SJNRB4U. 
With the generic methods proposal having been accepted (and hopefully in Go 1.27), this takes advantage of the
generic methods to design a parser combinator. Now and before generic methods, to accomplish a more functional
style one would typically use a sequential application of functions.
```go
	s := parser.StartKeeping(nameParser)
	s1 := parser.AppendSkipping(s, whitespaceParser)
	s2 := parser.AppendSkipping(s1, Exactly("="))
	s3 := parser.AppendSkipping(s2, whitespaceParser)
	s4 := parser.AppendKeeping(s3, valueParser)
  bindingParser := parser.Apply2(s4,
    func(name string, value BindingValue) Binding {
	 		return Binding{Name: name, Value: value}
	 	})
```
By explicitly accepting a generic type as the first function argument, we can still use generics. However, with
generic methods we can reduce the verbosity for sequential function applications on generic types.
```go
   bindingParser := nameParser
	 	.Skip(whitespaceParser)
	 	.Skip(Exactly("="))
	 	.Skip(whitespaceParser)
	 	.Append(valueParser)
	 	.Apply2(func(name string, value BindingValue) Binding {
	 		return Binding{Name: name, Value: value}
	 	})
```
I for one welcome this change to developer ergonmics and wait for its arrival to Go.

# Structure
All files ending in *.methodgen are written in the style of Go generic methods. For now, they are not valid Go code 
(and maybe with generic methods there are still some problems with correctness). However, once generic methods are
in Go it should be as easy as updating the go version in go.mod and overwriting all *.go files by their respective 
*.methodgen if it exists.
