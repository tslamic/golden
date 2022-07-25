[![CircleCI](https://circleci.com/gh/tslamic/golden.svg?style=svg)](https://circleci.com/gh/tslamic/golden) [![Go Report Card](https://goreportcard.com/badge/github.com/tslamic/golden)](https://goreportcard.com/report/github.com/tslamic/golden)

# :large_orange_diamond: Golden

Excruciatingly simple golden file handling. 
If you're unsure what golden files are, check [this video](https://youtu.be/8hQG7QlcLBk?t=735). 

## How to use it?

```  
go get -u github.com/tslamic/golden/v2
```

If your files are JSON, XML or plain text, you're good to go:

```go
func TestJSON(t *testing.T) {
    expected := newGreeter()
    gf := golden.File("testdata/hello.json")
    gf.Equals(t, expected)
}

func TestXML(t *testing.T) {
    expected := newCatalog()
    gf := golden.File("testdata/catalog.xml")
    gf.Equals(t, expected)
}

func TestText(t *testing.T) {
    expected := "Hello World!"
    gf := golden.File("testdata/world.txt")
    gf.Equals(t, expected)
}
```

For custom marshalling or attributes, you can do the following:

```go
gf := File(name, func(attrs *Attrs) {
    attrs.Marshaller = func(v interface{}) ([]byte, error) {
        b := new(bytes.Buffer)
        err := gob.NewEncoder(b).Encode(v)
        return b.Bytes(), err
    }
})
```

and if you want to tweak the underlying `[]byte` slices that will be compared, apply some transformers:

```go
func TestJSONWhitespace(t *testing.T) {
    expected := newGreeter()
    gf := File("testdata/hello_whitespace.json").Apply(StripWhitespace)
    gf.Equals(t, expected)
}
```

Make sure to run the tests with the `-update` command line argument to populate the files, 
then drop the flag for all subsequent tests. 

## License

The MIT License. 
