package main

import "fmt"

func main() {
	// Want to make a parser like https://www.youtube.com/watch?v=x5p_SJNRB4U&t=1429s
	// Take advantage of type parameters on methods to finally have
	//
	// func (p Parser[T]) Keep[U any](Parser[U]) Parser[...T... U...]
	//
	// bindingParser := nameParser
	// 	.Skip(whitespaceParser)
	// 	.Skip(Exactly("="))
	// 	.Skip(whitespaceParser)
	// 	.Append(s3, valueParser)
	// 	.Apply2(func(name string, value BindingValue) Binding {
	// 		return Binding{Name: name, Value: value}
	// 	})
	//
	// This avoids the older way using sequences
	//
	// s := parser.StartKeeping(nameParser)
	// s1 := parser.AppendSkipping(s, whitespaceParser)
	// s2 := parser.AppendSkipping(s1, Exactly("="))
	// s3 := parser.AppendSkipping(s2, whitespaceParser)
	// s4 := parser.AppendKeeping(s3, valueParser)
	// bindingParser := parser.Apply2(s4,
	// 	func(name string, value BindingValue) Binding {
	// 		return Binding{Name: name, Value: value}
	// 	})
	p := NewConfigParser()
	parsed, err := p.ConfigurationParser.Parse("[Name=false]")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", parsed)
}
