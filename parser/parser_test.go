package parser_test

import (
	"fmt"
	"testing"

	"github.com/ngdinhtoan/hari/parser"
)

func TestParse(t *testing.T) {
	data := []byte(`
{
    "name": "Toan",
    "age": 30,
    "active": true,
    "children": [{
        "name": "Hachi",
        "age": 3
    },
    {
        "name": "Yuri",
        "age": 2
    }],
	"group": ["family", "work"],
	"contact": {
		"phone": "9282882822",
		"email": "mail@gmail.com"
	}
}
`)

	rs := make(chan *parser.Struct)
	errs := make(chan error)
	done := make(chan bool)

	go parser.Parse("Person", data, rs, errs, done)

	for {
		select {
		case <-done:
			close(rs)
			close(errs)
			close(done)
			return
		case s := <-rs:
			w := &stringWriter{}
			s.WriteTo(w)
			fmt.Println(w.Data)
		}
	}
}

type stringWriter struct {
	Data string
}

func (sw *stringWriter) Write(p []byte) (n int, err error) {
	sw.Data += string(p)
	return len(p), nil
}
