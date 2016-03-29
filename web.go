package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const HEADER_STREAM = "X-Omnilog-Stream"

func web(out *log.Logger) {
	http.HandleFunc("/", HandleWeb(out))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}

func HandleWeb(out *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		streams, ok := req.Header[HEADER_STREAM]
		if !ok {
			log.Println("missing", HEADER_STREAM)
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
				// @todo do not assume tab delim
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
