
# :large_orange_diamond: Golden  
  
```  
go get -u github.com/tslamic/golden
```  

Excruciatingly simple golden file handling. To update the files, use the `-update` flag, or roll your own update mechanism:

```go
gf := golden.File("testdata/hello.json", JSON, func(d *Data) {
  d.Update = true // ¯\_(ツ)_/¯
})
```
  
## Examples  

```go
// JSON  
hello := &Greeting{Greeting: "Hello, World!"}  
gf := golden.File("testdata/hello.json", JSON)  
gf.Equals(t, greeter) // t is of type *testing.T

// Ignore all whitespace
gf := golden.File("testdata/boring.xml", XML, IgnoreWhitespace)  
if diff, err := gf.Eq(something); err != nil {  
    t.Fatal(err, "difference: ", diff)
}
```

## License

The MIT License (MIT), Copyright (c) 2019 tslamic