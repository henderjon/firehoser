package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCoalesce(t *testing.T) {
	ch := make(chan []byte, 7) // buffered channels are 0-based and we're sending 8 lines ...

	b := &bytes.Buffer{}
	go coalesce(ch, b)

	homeHandle := Adapt(parseRequest(ch), parseCustomHeader, checkAuth(), ensurePost(), checkShutdown())

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

	time.Sleep(5 * time.Second)

	expected := 416 // 415 + the last newline added by coalesce()
	if b.Len() != expected {
		t.Error("Coalesce error: \nexpected\n", expected, "\nactual\n", b.Len())
	}

}

func TestNewWriteCloser(t *testing.T) {
	var ok bool

	if _, ok = newWriteCloser().(io.WriteCloser); !ok {
		t.Error("Interface error: WriteSplitter (lines)")
	}

	splitByteCount = 50
	if _, ok = newWriteCloser().(io.WriteCloser); !ok {
		t.Error("Interface error: WriteSplitter (bytes)")
	}

	forceStdout = true
	if _, ok = newWriteCloser().(io.WriteCloser); !ok {
		t.Error("Interface error: Stdout")
	}
}
