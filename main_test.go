package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
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
	req.Body.Close()

	wg.Wait()

	expected := 416 // 415 + the last newline added by coalesce()
	if b.Len() != expected {
		t.Error("Coalesce error: \nexpected\n", expected, "\nactual\n", b.Len())
	}

}
