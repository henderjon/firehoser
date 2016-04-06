package main

import (
	"bufio"
	"io"
	"net/http"
)

// A custom header previously used to name the stream(s) to prepend to the line
// data. This isn't very useful yet
const (
	customHeader = "X-Omnilog-Stream"
	methodPost   = "POST"
)

// scans the body of the POST request and writes each line. Currently prepends
// the stream name (from the request header) to each line. This feature is less
// useful each passing minute.
func handleWeb(out io.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		s := http.StatusServiceUnavailable

		// graceful shutdown, reject new requests
		if isShutdownMode() {
			http.Error(w, (newResponse(s, 0)).Json(), s)
			return
		}

		// if we get here, don't let the program goroutine die before the goroutine finishes
		wg.Add(1)
		defer wg.Done() // cover the short-circuit returns

		s = http.StatusOK

		// ensure a POST
		if req.Method != methodPost {
			s = http.StatusMethodNotAllowed
		}

		// must have custom header (@TODO future validation?)
		if _, ok := req.Header[customHeader]; !ok {
			s = http.StatusBadRequest
		}

		// check if the Authorization header matches the provided password
		if !checkAuth(req.Header) {
			s = http.StatusForbidden
		}

		if s != http.StatusOK {
			http.Error(w, (newResponse(s, 0)).Json(), s)
			return
		}

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

		http.Error(w, (newResponse(s, rn)).Json(), s)
	}
}

// check the reqeust headers for the Authorization header (e.g. 'Authorization: Bearer this-is-a-string')
func checkAuth(h http.Header) bool {

	// if no password was given to the server, leave the doors wide open
	if pswd == "" {
		return true
	}

	// make sure the header exists
	a, ok := h["Authorization"]
	if !ok {
		return false
	}

	// go returns a map, so loop over it
	for _, v := range a {
		// Bearer
		if len(v) < 7 {
			continue
		}
		if v[7:] == pswd {
			return true
		}
	}

	return false
}
