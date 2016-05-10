package raw

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
	// "syscall"
)

func main() {
	inbound := make(chan []byte, 5000)

	// var rLimit syscall.Rlimit
	// rLimit.Cur = 4096
	// syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	// if capacity == 0 -> stdout?
	for t := 0; t < 2048; t += 1 {
		go func() {
			for {
				b := <-inbound // pull data out of the channel
				name := filepath.Join("/mnt/logs", time.Now().Format(time.RFC3339Nano))
				ioutil.WriteFile(name, b, 0644)
			}
		}()
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", http.StripPrefix("/", fs))
	http.Handle("/log", parseRequest(inbound))
	if e := http.ListenAndServe(":80", nil); e != nil {
		log.Fatal(e)
	}
}

// parseRequest returns a handler that reads the body of the request and sends
// it on the channel to be coalesced
func parseRequest(data chan []byte) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		s := http.StatusOK

		b, e := ioutil.ReadAll(req.Body)
		if e != nil {
			log.Println(e)
			s = http.StatusBadRequest
			http.Error(rw, "no", s)
			return
		}

		http.Error(rw, "yup", s)

		data <- b
	})
}
