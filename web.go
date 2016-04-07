package main

import (
	"bufio"
	"io"
	"net/http"
	"sync"
)

// A custom header previously used to name the stream(s) to prepend to the line
// data. This isn't very useful yet
const (
	customHeader = "X-Omnilog-Stream"
	methodPost   = "POST"
)

var (
	wg sync.WaitGroup // ensure that our goroutines finish before shut down
)

type webHandler struct {
	io.Writer
}

// scans the body of the POST request and writes each line. Currently prepends
// the stream name (from the request header) to each line. This feature is less
// useful each passing minute.
func handleWeb(out io.Writer) http.HandlerFunc {
	return parseCustomHeader(func(rw http.ResponseWriter, req *http.Request) {
		// if we get here, don't let the program goroutine die before the goroutine finishes
		wg.Add(1)
		defer wg.Done() // cover the short-circuit returns
		parseRequest(rw, req)
	})
}

func checkShutdown(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request){
		// graceful shutdown, reject new requests
		if isShutdownMode() {
			s := http.StatusServiceUnavailable
			http.Error(rw, (newResponse(s, 0)).Json(), s)
			return
		}
		fn(rw, req)
	}
}

// ensure a POST
func ensurePost(fn http.HandlerFunc) http.HandlerFunc {
	return checkShutdown(func(rw http.ResponseWriter, req *http.Request){
		s := http.StatusMethodNotAllowed
		// ensure a POST
		if req.Method != methodPost {
			http.Error(rw, (newResponse(s, 0)).Json(), s)
			return
		}
		fn(rw, req)
	})
}

// check the reqeust headers for the Authorization header (e.g. 'Authorization: Bearer this-is-a-string')
func checkAuth(fn http.HandlerFunc) http.HandlerFunc {
	return ensurePost(func(rw http.ResponseWriter, req *http.Request){
		// if no password was given to the server, leave the doors wide open
		if pswd == "" {
			fn(rw, req)
			return
		}

		s := http.StatusForbidden

		// make sure the header exists
		a, ok := req.Header["Authorization"]
		if !ok {
			http.Error(rw, (newResponse(s, 0)).Json(), s)
			return
		}

		// go returns a map, so loop over it
		for _, v := range a {
			// Bearer
			if len(v) < 7 {
				continue
			}
			if v[7:] == pswd {
				fn(rw, req)
				return
			}
		}

		http.Error(rw, (newResponse(s, 0)).Json(), s)
		return
	})
}

func parseCustomHeader(fn http.HandlerFunc) http.HandlerFunc {
	return checkAuth(func(rw http.ResponseWriter, req *http.Request){
		// must have custom header (@TODO future stream separation?)
		if _, ok := req.Header[customHeader]; !ok {
			s := http.StatusBadRequest
			http.Error(rw, (newResponse(s, 0)).Json(), s)
			return
		}
		fn(rw, req)
	})
}

func parseRequest(rw http.ResponseWriter, req *http.Request){
	defer req.Body.Close()
	s := http.StatusOK
	// read the request body
	rn := 0
	scanner := bufio.NewScanner(req.Body)
	for scanner.Scan() {

		if isBrokenPipe() {
			out = getDest(ioFile)
		}

		n, _ := out.Write(scanner.Bytes())
		rn += n

		if err := scanner.Err(); err != nil {
			break
		}
	}

	http.Error(rw, (newResponse(s, rn)).Json(), s)
}
