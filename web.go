package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// A custom header previously used to name the stream(s) to prepend to the line
// data. This isn't very useful yet
const headerStream = "X-Omnilog-Stream"
const methodPost = "POST"

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
			log.Println(errShutdown)
			return
		}

		// if we get here, don't let the program goroutine die before the goroutine finishes
		wg.Add(1)

		// ensure a POST
		if req.Method != methodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			enc.Encode(&response{
				errMethodNotAllowed, 0,
			})
			log.Println(errMethodNotAllowed)
			return
		}

		// must have custom header (@TODO future validation?)
		if _, ok := req.Header[headerStream]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(&response{
				errBadRequest, 0,
			})
			log.Println(errBadRequest)
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
		wg.Done()
		return
	}
}
