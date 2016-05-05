package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var body = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam id turpis sit amet nibh tempus fringilla. Vivamus lacinia metus et neque dignissim egestas eu non sem. Phasellus pretium augue ultrices, tristique dui vel, euismod est. Maecenas egestas mauris quis diam maximus laoreet. Curabitur mattis, diam sed mollis posuere, felis ipsum rhoncus nulla, non gravida metus ipsum lobortis orci. Mauris quis tellus et enim elementum fermentum.
`

func Test_OK(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan []byte, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch, &wg), checkHeader(customHeader), checkAuth(""), checkMethod(methodPost), checkShutdown(nil))

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(body))
	req.Header.Add(customHeader, "test")

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Status error: expected", http.StatusOK, "actual", w.Code)
	}

	expected := `{"Status":200,"Message":"OK","Received":442}` + "\n"
	if w.Body.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", w.Body.String())
	}
}

func Test_OKAuth(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan []byte, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch, &wg), checkHeader(customHeader), checkAuth("PASSWORD"), checkMethod(methodPost), checkShutdown(nil))

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(body))
	req.Header.Add(customHeader, "test")
	req.Header.Add("Authorization", "Bearer PASSWORD")

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Status error: expected", http.StatusOK, "actual", w.Code)
	}

	expected := `{"Status":200,"Message":"OK","Received":442}` + "\n"
	if w.Body.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", w.Body.String())
	}
}

func Test_BadRequest(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan []byte, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch, &wg), checkHeader(customHeader))

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(body))
	// req.Header.Add(customHeader, "test")

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Error("Status error: expected", http.StatusBadRequest, "actual", w.Code)
	}

	expected := `{"Status":400,"Message":"Bad Request","Received":0}` + "\n"
	if w.Body.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", w.Body.String())
	}
}

func Test_MethodNotAllowed(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan []byte, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch, &wg), checkMethod(methodPost))

	req, _ := http.NewRequest("GET", "", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Error("Status error: expected", http.StatusMethodNotAllowed, "actual", w.Code)
	}

	expected := `{"Status":405,"Message":"Method Not Allowed","Received":0}` + "\n"
	if w.Body.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", w.Body.String())
	}
}

func Test_Forbidden(t *testing.T) { // no password
	var wg sync.WaitGroup
	ch := make(chan []byte, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch, &wg), checkAuth("PASSWORD"))

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Error("Status error: expected", http.StatusForbidden, "actual", w.Code)
	}

	expected := `{"Status":403,"Message":"Forbidden","Received":0}` + "\n"
	if w.Body.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", w.Body.String())
	}
}

func Test_Forbidden2(t *testing.T) { // wrong password
	var wg sync.WaitGroup
	ch := make(chan []byte, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	homeHandle := Adapt(parseRequest(ch, &wg), checkAuth("PASSWORD"))

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(body))
	req.Header.Add("Authorization", "PSWD")

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Error("Status error: expected", http.StatusForbidden, "actual", w.Code)
	}

	expected := `{"Status":403,"Message":"Forbidden","Received":0}` + "\n"
	if w.Body.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", w.Body.String())
	}
}

func Test_Shutdown(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan []byte, 9) // buffer the chan to avoid blocking since we're not reading OUT of the channel
	sh := make(chan struct{}, 1)
	homeHandle := Adapt(parseRequest(ch, &wg), checkShutdown(sh))

	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(body))
	req.Header.Add(customHeader, "test")

	close(sh)

	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Error("Status error: expected", http.StatusServiceUnavailable, "actual", w.Code)
	}

	expected := `{"Status":503,"Message":"Service Unavailable","Received":0}` + "\n"
	if w.Body.String() != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", w.Body.String())
	}
}
