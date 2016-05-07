package main

import (
	"bytes"
	"github.com/henderjon/omnilogger/shutdown"
	"io"
	"time"
)

// worker wraps a bytes.Buffer in order to attach the coalesce func
type worker struct {
	cutoff int
	bytes.Buffer
}

// newWorker creates a new worker with and allocates the memory necessary
func newWorker(max int) *worker {
	worker := &worker{cutoff: max}
	worker.Grow(max) // allocating memory here seemed to help performance
	return worker
}

// coalesce runs in it's own goroutine and ranges over a channel and writing the
// data to an io.Writer. All our goroutines send data on this channel and this
// func coalesces them in to one stream.
func (w *worker) coalesce(inbound chan []byte, fw io.Writer, sig shutdown.SignalChan) {
	wg.Add(1)
	for {
		select {
		case b := <-inbound: // pull data out of the channel
			n := len(b)
			go fw.Write(b)
			w.Reset()

			byteCount.IncrBy(uint64(n))
			hitCount.IncrBy(uint64(1))

		case <-time.After(closeInterval): // after 10 minutes of inactivity flush to disk
			go fw.Write(w.Bytes())
			w.Reset()

		case <-sig: // if we're shutting down, make sure we flush to disk first
			fw.Write(w.Bytes())
			wg.Done()
			return
		}
	}
}
