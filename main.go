package main

import (
	"bytes"
	"flag"
	"github.com/henderjon/omnilogger/counter"
	"github.com/henderjon/omnilogger/shutdown"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	Kilobyte = 1024
	Megabyte = Kilobyte * Kilobyte
	Gigabyte = Kilobyte * Megabyte
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
	capacity      = Kilobyte * 64             // the default capacity of the buffers
	helpLogger    = log.New(os.Stderr, "", 0) // log to stderr without the timestamps
	closeInterval = 10 * time.Minute          // how often to close our file and open a new one
	byteCount     = counter.NewCounter()      // keep track of how many bytes total have been received
	hitCount      = counter.NewCounter()      // keep track of how many bytes total have been received
)

type worker struct {
	dir, prefix string
	bytes.Buffer
}

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

	capacity = size * Kilobyte
	if scale {
		capacity = size * Megabyte
	}
	// if capacity == 0 -> stdout

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
	workers := make([]*worker, 4)
	inbound := make(chan []byte, requestBuffer)
	shutdownCh := make(shutdown.SignalChan)

	for t := 0; t < numWorkers; t += 1 {
		workers[t] = &worker{
			dir: logDir, prefix: flag.Arg(0),
		}
		go workers[t].coalesce(inbound, shutdownCh)
	}

	go shutdown.Watch(shutdownCh, destructor)

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", http.StripPrefix("/", fs))
	http.Handle("/log", Adapt(parseRequest(inbound, &wg), parseCustomHeader, checkAuth(pswd), ensurePost(), checkShutdown(shutdownCh)))
	if e := http.ListenAndServe(":"+port, nil); e != nil {
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
		case b := <-inbound: // pull data out of the channel
			if w.Len() > capacity {
				go fwrite(w.name(), w.Bytes())
				w.Reset()
			}

			n, _ := w.Write(b)
			byteCount.IncrBy(uint64(n))
			hitCount.IncrBy(uint64(1))

		case <-time.After(closeInterval): // after 10 minutes of inactivity close the file
			go fwrite(w.name(), w.Bytes())
			w.Reset()

		case <-shutdownCh: // if we're shutting down, make sure we flush to disk first
			fwrite(w.name(), w.Bytes())
			wg.Done()
			return
		}
	}
}

func (w *worker) name() string {
	return filepath.Join(w.dir, w.prefix + time.Now().Format(time.RFC3339Nano))
}

func destructor() {
	wg.Wait()
	helpLogger.Printf(".collected %dk from %d hits in %s", byteCount.Current(Kilobyte), hitCount.Current(0), byteCount.Since())
}
