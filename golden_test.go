package golden_test

import (
	"bytes"
	"encoding/gob"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/tslamic/golden/v2"
)

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

func TestTextByte(t *testing.T) {
	expected := []byte("Hello World!")
	gf := golden.File("testdata/world.txt")
	gf.Equals(t, expected)
}

func TestJSONWhitespace(t *testing.T) {
	expected := newGreeter()
	gf := golden.File("testdata/hello_whitespace.json").Apply(golden.StripWhitespace)
	gf.Equals(t, expected)
}

func TestXMLWhitespace(t *testing.T) {
	expected := newCatalog()
	gf := golden.File("testdata/catalog_whitespace.xml").Apply(golden.StripWhitespace)
	gf.Equals(t, expected)
}

func TestTextErr(t *testing.T) {
	expected := "Oh noes!"
	gf := golden.File("testdata/world.txt")
	_, err := gf.Eq(expected)
	if !errors.Is(err, golden.ErrNotEqual) {
		t.Fatal("unexpected match")
	}
}

func TestMissingFile(t *testing.T) {
	gf := golden.File("testdata/missing.txt")
	if _, err := gf.Eq(struct{}{}); err == nil {
		t.Fatal("unexpected match")
	}
}

func TestCustomMarshaller(t *testing.T) {
	v := &vector{X: 0, Y: 1, Z: 2}
	f, err := ioutil.TempFile("", "*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	err = gob.NewEncoder(f).Encode(v)
	if err != nil {
		t.Fatal(err)
	}
	gf := golden.File(f.Name(), func(attrs *golden.Attrs) {
		attrs.Marshaller = func(v interface{}) ([]byte, error) {
			b := new(bytes.Buffer)
			err := gob.NewEncoder(b).Encode(v)
			return b.Bytes(), err
		}
	})
	gf.Equals(t, v)
}

type greeter struct {
	Greeting string `json:"greeting"`
}

func newGreeter() *greeter {
	return &greeter{Greeting: "Hello, World!"}
}

type catalog struct {
	XMLName xml.Name `xml:"catalog"`
	Text    string   `xml:",chardata"`
	Books   []*book  `xml:"book"`
}

type book struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id,attr"`
	Author      string `xml:"author"`
	Title       string `xml:"title"`
	Genre       string `xml:"genre"`
	Price       string `xml:"price"`
	PublishDate string `xml:"publish_date"`
	Description string `xml:"description"`
}

func newCatalog() *catalog {
	return &catalog{
		Books: []*book{
			{
				ID:          "bk101",
				Author:      "Gambardella, Matthew",
				Title:       "XML Developer Guide",
				Genre:       "Computer",
				Price:       "44.95",
				PublishDate: "2000-10-01",
				Description: "An in-depth look at creating applications with XML.",
			},
			{
				ID:          "bk102",
				Author:      "Ralls, Kim",
				Title:       "Midnight Rain",
				Genre:       "Fantasy",
				Price:       "5.95",
				PublishDate: "2000-12-16",
				Description: "A former architect battles corporate zombies, an evil sorceress, " +
					"and her own childhood to become queen of the world.",
			},
		},
	}
}

type vector struct {
	X, Y, Z int
}
