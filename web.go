package main

import (
	"bufio"
	"encoding/json"
	"io"
	// "log"
	"net/http"
)

// A custom header previously used to name the stream(s) to prepend to the line
// data. This isn't very useful yet
const (
	customHeader = "X-Omnilog-Stream"
	methodPost   = "POST"
	// authHeader   = "Authorization: Bearer "
)

// run a small web server
func web(out io.Writer, port string) {
	http.HandleFunc("/", handleWeb(out))
	http.ListenAndServe(":"+port, nil)
}

// scans the body of the POST request and writes each line. Currently prepends
// the stream name (from the request header) to each line. This feature is less
// useful each passing minute.
func handleWeb(out io.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		enc := json.NewEncoder(w)

		// graceful shutdown, reject new requests
		if isShutdownMode() {
			w.WriteHeader(http.StatusServiceUnavailable)
			enc.Encode(&response{
				errShutdown, 0,
			})
			// log.Println(errShutdown)
			return
		}

		// if we get here, don't let the program goroutine die before the goroutine finishes
		wg.Add(1)
		defer wg.Done() // cover the short-circuit returns

		// ensure a POST
		if req.Method != methodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			enc.Encode(&response{
				errMethodNotAllowed, 0,
			})
			// log.Println(errMethodNotAllowed)
			return
		}

		// must have custom header (@TODO future validation?)
		if _, ok := req.Header[customHeader]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(&response{
				errBadRequest, 0,
			})
			// log.Println(errBadRequest)
			return
		}

		// check if the Authorization header matches the provided password
		if !checkAuth(req.Header) {
			w.WriteHeader(http.StatusForbidden)
			enc.Encode(&response{
				errForbidden, 0,
			})
			// log.Println(errForbidden)
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

		w.WriteHeader(http.StatusOK)
		enc.Encode(&response{
			success, rn,
		})
		req.Body.Close()
		return
	}
}

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
