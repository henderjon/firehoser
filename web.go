package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
)

// A custom header previously used to name the stream(s) to prepend to the line
// data. This isn't very useful yet
const HeaderStream = "X-Omnilog-Stream"
const MethodPost = "POST"

// run a small web server
func web(out *log.Logger, port string) {
	http.HandleFunc("/", HandleWeb(out))
	http.ListenAndServe(":"+port, nil)
}

// scans the body of the POST request and writes each line. Currently prepends
// the stream name (from the request header) to each line. This feature is less
// useful each passing minute.
func HandleWeb(out *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		enc := json.NewEncoder(w)
		if req.Method != MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			enc.Encode(&response{
				ErrMethodNotAllowed, 0,
			})
			log.Println(ErrMethodNotAllowed)
			return
		}

		if _, ok := req.Header[HeaderStream]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(&response{
				ErrBadRequest, 0,
			})
			log.Println(ErrBadRequest)
			return
		}

		rn := 0
		scanner := bufio.NewScanner(req.Body)
		for scanner.Scan() {
			out.Println(scanner.Text())
			rn += len(scanner.Text())

			if err := scanner.Err(); err != nil {
				break
			}
		}

		enc.Encode(&response{
			Success, rn,
		})
	}
}
