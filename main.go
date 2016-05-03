package main

import (
	"net/http"
	"bytes"
	"log"
	"sync"
	"github.com/henderjon/omnilogger/shutdown"
)

const (
	Kilobyte = 1024
	Megabyte = Kilobyte * Kilobyte
	Gigabyte = Kilobyte * Megabyte

	Capacity = Kilobyte * 64
)

var (
	wg sync.WaitGroup // ensure that our goroutines finish before shut down
)

type worker struct {
	bytes.Buffer
}

func main() {
	workers    := make([]*worker, 4)
	inbound    := make(chan []byte, 5000)
	shutdownCh := make(shutdown.SignalChan)

	for t := 0; t <= 3; t += 1 {
		workers[t] = &worker{}
		go workers[t].coalesce(inbound, shutdownCh)
	}

	go shutdown.Watch(shutdownCh, destructor)

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", http.StripPrefix("/", fs))
	http.Handle("/log", Adapt(parseRequest(inbound, &wg), parseCustomHeader, checkAuth(), ensurePost(), checkShutdown(shutdownCh)))
	if e := http.ListenAndServe(":8080", nil); e != nil {
		log.Fatal(e)
	}
}

// coalesce runs in it's own goroutine and ranges over a channel and writing the
// data to an io.Writer. All our goroutines send data on this channel and this
// func coalesces them in to one stream.
func (w *worker) coalesce(inbound chan []byte, shutdownCh chan struct{}) {
	wg.Add(1)
	for {
		select {
		case b := <-inbound:
			if w.Len() > Capacity {
				go fwrite(w.Bytes())
				w.Reset()
			}
			w.Write(b)
		case <-shutdownCh :
			fwrite(w.Bytes())
			wg.Done()
			return
		}
	}
}

func destructor(){
	wg.Wait()
}
