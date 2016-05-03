package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"bytes"
	"log"
	"sync"
)

const (
	Kilobyte = 1024
	Megabyte = Kilobyte * Kilobyte
	Gigabyte = Kilobyte * Megabyte

	Capacity = Megabyte * 64
)

var (
	wg sync.WaitGroup // ensure that our goroutines finish before shut down
)

type worker struct {
	bytes.Buffer
}

func NewWorker(id int) (w *worker) {
	return &worker{}
}

func main() {
	workers := make([]*worker, 4)
	inbound := make(chan []byte, 5000)
	shutdown := make(chan struct{})

	for t := 0; t <= 3; t += 1 {
		workers[t] = &worker{}
		go workers[t].coalesce(inbound, shutdown)
	}

	go monitorStatus(shutdown, &wg)

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", http.StripPrefix("/", fs))
	http.Handle("/log", parseRequest(inbound))
	http.ListenAndServe(":8080", nil)
}

// coalesce runs in it's own goroutine and ranges over a channel and writing the
// data to an io.Writer. All our goroutines send data on this channel and this
// func coalesces them in to one stream.
func (w *worker) coalesce(inbound chan []byte, shutdown chan struct{}) {
	// wg.Add(1)
	for {
		select {
		case b := <-inbound:
			if w.Len() > Capacity {
				go Save(w.Bytes())
				w.Reset()
			}
			w.Write(b)
		case <-shutdown :
			// log.Println(w.Len())
			Save(w.Bytes())
			// wg.Done()
			return
		}
	}
}

// parseRequest returns a handler that reads the body of the request and sends
// it on the channel to be coalesced
func parseRequest(data chan []byte) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		s := http.StatusOK
		// read the request body
		a, _ := ioutil.ReadAll(req.Body)
		data <- a
		// log.Println(len(a))
		http.Error(rw, "success", s)
	})
}

func Save(payload []byte) {
	if len(payload) == 0 {
		return
	}
	filename := filepath.Join("TL1-" + time.Now().Format(time.RFC3339Nano))
	f, e := os.Create(filename)
	if e != nil {
		log.Fatal(e)
	}
	defer f.Close()
	f.Write(payload)
	// log.Println(string(payload))
}
