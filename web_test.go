package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOK(t *testing.T) {
	ch := make(chan *payload, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch), parseCustomHeader, checkAuth(), ensurePost(), checkShutdown(nil))

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
	if w.Code != http.StatusOK {
		t.Error("Status error: expected", http.StatusOK, "actual", w.Code)
	}

	b := &bytes.Buffer{}
	b.ReadFrom(w.Body)

	expected := `{"Status":200,"Message":"OK","Received":408}` + "\n"
	if b.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", b.String())
	}
}

func TestOKAuth(t *testing.T) {
	ch := make(chan *payload, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch), parseCustomHeader, checkAuth(), ensurePost(), checkShutdown(nil))

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
	req.Header.Add("Authorization", "Bearer PASSWORD")

	pswd = "PASSWORD"

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Status error: expected", http.StatusOK, "actual", w.Code)
	}

	b := &bytes.Buffer{}
	b.ReadFrom(w.Body)

	expected := `{"Status":200,"Message":"OK","Received":408}` + "\n"
	if b.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", b.String())
	}
}

func TestBadRequest(t *testing.T) {
	ch := make(chan *payload, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch), parseCustomHeader)

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString("..."))
	// req.Header.Add(customHeader, "test")

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Error("Status error: expected", http.StatusBadRequest, "actual", w.Code)
	}

	b := &bytes.Buffer{}
	b.ReadFrom(w.Body)

	expected := `{"Status":400,"Message":"Bad Request","Received":0}` + "\n"
	if b.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", b.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	ch := make(chan *payload, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch), ensurePost())

	req, _ := http.NewRequest("GET", "", bytes.NewBufferString("..."))

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Error("Status error: expected", http.StatusMethodNotAllowed, "actual", w.Code)
	}

	b := &bytes.Buffer{}
	b.ReadFrom(w.Body)

	expected := `{"Status":405,"Message":"Method Not Allowed","Received":0}` + "\n"
	if b.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", b.String())
	}
}

func TestForbidden(t *testing.T) {
	ch := make(chan *payload, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch), checkAuth())

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString("..."))

	pswd = "PASSWORD"

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Error("Status error: expected", http.StatusForbidden, "actual", w.Code)
	}

	b := &bytes.Buffer{}
	b.ReadFrom(w.Body)

	expected := `{"Status":403,"Message":"Forbidden","Received":0}` + "\n"
	if b.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", b.String())
	}
}

func TestForbidden2(t *testing.T) {
	ch := make(chan *payload, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch), checkAuth())

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString("..."))
	req.Header.Add("Authorization", "PSWD")

	pswd = "PASSWORD"

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Error("Status error: expected", http.StatusForbidden, "actual", w.Code)
	}

	b := &bytes.Buffer{}
	b.ReadFrom(w.Body)

	expected := `{"Status":403,"Message":"Forbidden","Received":0}` + "\n"
	if b.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", b.String())
	}
}

func TestShutdown(t *testing.T) {
	ch := make(chan *payload, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	sh := make(chan struct{}, 1)
	homeHandle := Adapt(parseRequest(ch), checkShutdown(sh))

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString("..."))
	req.Header.Add(customHeader, "test")

	pswd = ""

	close(sh)

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Error("Status error: expected", http.StatusServiceUnavailable, "actual", w.Code)
	}

	b := &bytes.Buffer{}
	b.ReadFrom(w.Body)

	expected := `{"Status":503,"Message":"Service Unavailable","Received":0}` + "\n"
	if b.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", b.String())
	}
}
