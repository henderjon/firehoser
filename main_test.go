package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// a protected Write Closer allows Close to be called on os.Stdout without the danger
// of Stdout actually being closed
type pwc struct {
	io.Writer
}

// satisfies the io.Close interface in such a way as to avoid closing anything
func (pwc) Close() error {
	return nil
}

func TestCoalesce(t *testing.T) {
	in := make(chan *payload, 7) // buffered channels are 0-based and we're sending 8 lines ...
	out := make(chan int, 7)     // buffered channels are 0-based and we're sending 8 lines ...

	go coalesce(in, out, func(dir, prefix string) writeCloserRecycler {
		return func(io.WriteCloser) io.WriteCloser {
			return &pwc{&bytes.Buffer{}}
		}
	})

	homeHandle := Adapt(parseRequest(in), func(h http.Handler) http.Handler {
		return h
	})

	mockD := bytes.NewBufferString(`Lorem ipsum dolor sit amet consectetur adipiscing elit
Cras in lacinia eros Aliquam aliquet sapien a
Ut mauris orci varius et cursus sed blandit
Mauris iaculis ac magna non tincidunt In rhoncus
Pellentesque quis erat quis ex aliquam porttitor Vestibulum
Pellentesque nec mollis nibh interdum eleifend nisl Donec
id commodo urna sed tempus mi Vestibulum facilisis
imperdiet dolor sed sollicitudin Proin in lectus sed`)

	req, _ := http.NewRequest("POST", "", mockD)
	req.Header.Add(customHeader, "test")
	pswd = ""

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)

	var total int
Loop:
	for {
		select {
		case x := <-out:
			total += x
		case <-time.After(5 * time.Second):
			break Loop // how do you test an infinite loop, kill it with fire
		}
	}

	expected := 416 // 415 + the last newline added by coalesce()
	if total != expected {
		t.Error("Coalesce error: \nexpected\n", expected, "\nactual\n", total)
	}

}

func TestNewWriteCloser(t *testing.T) {
	var ok bool
	var O writeCloserRecycler

	O = writeCloser("", "")
	if _, ok = O(nil).(io.WriteCloser); !ok {
		t.Error("Interface error: WriteSplitter (lines)")
	}

	byBytes = true
	O = writeCloser("", "")
	if _, ok = O(nil).(io.WriteCloser); !ok {
		t.Error("Interface error: WriteSplitter (bytes)")
	}

	forceStdout = true
	O = writeCloser("", "")
	if _, ok = O(nil).(io.WriteCloser); !ok {
		t.Error("Interface error: Stdout")
	}
}
