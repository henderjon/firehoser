package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// A custom header previously used to name the stream(s) to prepend to the line
// data. This isn't very useful yet
const HEADER_STREAM = "X-Omnilog-Stream"

// run a small web server
func web(out *log.Logger) {
	http.HandleFunc("/", HandleWeb(out))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}

// scans the body of the POST request and writes each line. Currently prepends
// the stream name (from the request header) to each line. This feature is less
// useful each passing minute.
func HandleWeb(out *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		streams, ok := req.Header[HEADER_STREAM]
		if !ok {
			// log.Println("missing", HEADER_STREAM)
			json.NewEncoder(w).Encode(&response{
				"error: missing header", 0, 0,
			})
			return
		}

		var rn, wn int

		scanner := bufio.NewScanner(req.Body)
		for scanner.Scan() {
			for i := 0; i < len(streams); i += 1 {
				wn += len(streams[i]) + 1
				// @todo do not assume tab delim, get rid of stream names
				out.Print(streams[i], "\t", scanner.Text(), "\n")
			}

			rn += len(scanner.Text())
			wn += rn

			if err := scanner.Err(); err != nil {
				break
			}
		}

		json.NewEncoder(w).Encode(&response{
			"success", rn, wn,
		})
	}
}
