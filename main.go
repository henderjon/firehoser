package main

import (
	"bytes"
	"flag"
	"github.com/henderjon/omnilogger/counter"
	"github.com/henderjon/omnilogger/shutdown"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// the uint64 representation of a Kilobyte, Megabyte, Gigabyte as well as some program defaults
const (
	Kilobyte        = 1024
	Megabyte        = Kilobyte * Kilobyte
	Gigabyte        = Kilobyte * Megabyte
	defaultInterval = 10 * time.Minute
	defaultPrefix   = "log-omnilogs-"
	defaultPerms    = 0644
)

var (
	port          string                      // the port on which to listen
	pswd          string                      // a simple means of authentication
	requestBuffer int                         // the size of the incoming request buffer (channel)
	size          int                         // how many lines per log file
	scale         bool                        // how many bytes per log file
	numWorkers    int                         // how many bytes per log file
	logDir        string                      // the dir for the log file(s)
	help          bool                        // I forgot my options
	wg            sync.WaitGroup              // ensure that our goroutines finish before shut down
	closeInterval = defaultInterval           // how often to close our file and open a new one
	byteCount     = counter.NewCounter()      // keep track of how many bytes total have been received
	hitCount      = counter.NewCounter()      // keep track of how many bytes total have been received
	helpLogger    = log.New(os.Stderr, "", 0) // log to stderr without the timestamps
)

func init() {
	flag.StringVar(&port, "port", "8080", "The port used for the server.")
	flag.StringVar(&pswd, "auth", "", "If not empty, this is matched against the Authorization header (e.g. Authorization: Bearer my-password).")
	flag.IntVar(&requestBuffer, "buf", 500, "The size of the incoming request buffer. A zero (0) will disable buffering.")
	flag.IntVar(&numWorkers, "workers", 2, "The number of workers/buffers... most likely doesn't need changing and should probably not exceed the number of CPUs on the machine.")
	flag.IntVar(&size, "size", 64, "The size (in kilobytes) at which to split the log file(s).")
	flag.BoolVar(&scale, "m", false, "If set, -size will be in megabytes.")
	flag.StringVar(&logDir, "dir", "", "A dir to use for log files. The first arg after '--' is used as a filename prefix. (e.g. '% omnilogger -dir /path/to/log-dir -- file-prefix-')")
	flag.BoolVar(&help, "h", false, "Show this message.")
	flag.Parse()

	if help {
		helpLogger.Println("")
		helpLogger.Println("Omnilogger is an HTTP server that ingests log data from multiple sources to a common destination.")
		helpLogger.Println("")
		flag.PrintDefaults()
		helpLogger.Println("")
		os.Exit(0)
	}
}

func main() {
	inbound := make(chan []byte, requestBuffer)
	shutdownCh := make(shutdown.ShutdownChan)

	capacity := size * Kilobyte
	if scale {
		capacity = size * Megabyte
	}
	// if capacity == 0 -> stdout
	for t := 0; t < numWorkers; t += 1 {
		go newWorker(capacity).coalesce(inbound, &nameWriter{logDir, flag.Arg(0)}, shutdownCh)
	}

	go shutdown.Watch(shutdownCh, destructor)

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", http.StripPrefix("/", fs))
	http.Handle("/log", Adapt(parseRequest(inbound, &wg), parseCustomHeader, checkAuth(pswd), ensurePost(), checkShutdown(shutdownCh)))
	if e := http.ListenAndServe(":"+port, nil); e != nil {
		log.Fatal(e)
	}
}

// destructor is the func that gets called if we catch a shutdown signal. It waits
// for all goroutines to finish and then prints a final status message
func destructor() {
	wg.Wait()
	helpLogger.Printf(".collected %dm from %d hits in %s", byteCount.Current(Megabyte), hitCount.Current(0), byteCount.Since())
}

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
func (w *worker) coalesce(inbound chan []byte, fw io.Writer, shutdownCh chan struct{}) {
	wg.Add(1)
	for {
		select {
		case b := <-inbound: // pull data out of the channel
			if w.Len() >= w.cutoff {
				go fw.Write(w.Bytes())
				w.Reset()
			}

			n, _ := w.Write(b)

			byteCount.IncrBy(uint64(n))
			hitCount.IncrBy(uint64(1))

		case <-time.After(closeInterval): // after 10 minutes of inactivity close the file
			go fw.Write(w.Bytes())
			w.Reset()

		case <-shutdownCh: // if we're shutting down, make sure we flush to disk first
			fw.Write(w.Bytes())
			wg.Done()
			return
		}
	}
}
