[![CircleCI](https://circleci.com/gh/tslamic/golden.svg?style=svg)](https://circleci.com/gh/tslamic/golden) [![Go Report Card](https://goreportcard.com/badge/github.com/tslamic/golden)](https://goreportcard.com/report/github.com/tslamic/golden)

# :large_orange_diamond: Golden  
  
Excruciatingly simple golden file handling.

```  
go get -u github.com/tslamic/golden
```   

## How to use it?

```go
func TestJSON(t *testing.T) {
	greet := &greeter{Greeting: "Hello, World!"}

	gf := golden.File("testdata/hello.json", golden.JSON, golden.IgnoreWhitespace)
	gf.Equals(t, greet)
}
```

It's easy to provide custom attributes:

```go
gf := File("testdata/golden.file", func(d *golden.Data) {
	// apply custom attributes to d here.
})
```

You can add a custom `Marshaller`, `Differ`, and as many `Transformer` funcs as you'd like. To update the golden files, use the `-update` flag, or roll your own update mechanism.

## License

The MIT License. 
