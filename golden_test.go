package golden

import (
	"bytes"
	"encoding/xml"
	"testing"
	"time"
)

func TestErr(t *testing.T) {
	errs := map[*Data]error{
		File(""): ErrNoPath,
		File("hello.json", func(d *Data) {
			d.Marsh = nil
		}): ErrNoMarshaller,
		File("world.txt"): ErrUnsupportedType,
	}
	for d, e := range errs {
		_, err := d.Eq(struct{}{})
		if e != err {
			t.Fatalf("expected %s but got %s", e, err)
		}
	}
}

type meta struct {
	Timestamp time.Time `json:"timestamp"`
	ID        int       `json:"id"`
}

type greeter struct {
	Greeting string `json:"greeting"`
	Meta     *meta  `json:"meta"`
}

func TestAsJSON(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2019-12-01T17:06:14Z")
	if err != nil {
		t.Fatal(err)
	}
	g := &greeter{
		Greeting: "hello",
		Meta: &meta{
			Timestamp: timestamp,
			ID:        123456789,
		},
	}
	gf := File("testdata/hello.json", JSON)
	gf.Equals(t, g)
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

type catalog struct {
	XMLName xml.Name `xml:"catalog"`
	Text    string   `xml:",chardata"`
	Books   []*book  `xml:"book"`
}

func TestXML(t *testing.T) {
	c := catalog{
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
	gf := File("testdata/catalog.xml", XML, IgnoreWhitespace)
	if diff, err := gf.Eq(c); err != nil {
		t.Fatal(err, diff)
	}
}

const content = "Hello, World!"

func TestReadWrite(t *testing.T) {
	gf := File("testdata/world.txt")

	c := []byte(content)
	w, err := gf.Write(c)
	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, w)
	r, err := gf.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if w != r || !bytes.Equal(c, b) {
		t.Fatal("read/write not in sync")
	}
}

func TestStrip(t *testing.T) {
	values := map[string]string{
		" ":               "",
		"\n\t\v\r\f":      "",
		"abc\n\t123":      "abc123",
		"\vhello world\r": "helloworld",
	}
	for str, exp := range values {
		s := []byte(str)
		e := []byte(exp)
		if !bytes.Equal(strip(s), e) {
			t.Fatal("not equal", str, exp)
		}
	}
}
