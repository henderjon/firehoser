package main

import (
	"bytes"
	"flag"
	"github.com/henderjon/omnilogger/counter"
	"github.com/henderjon/omnilogger/shutdown"
	"log"
	"net/http"
	"os"
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
	size          int                      // how many lines per log file
	scale         bool                        // how many bytes per log file
	splitDir      string                      // the dir for the log file(s)
	help          bool                        // I forgot my options
	wg            sync.WaitGroup              // ensure that our goroutines finish before shut down
	capacity      = Kilobyte * 64             // the default capacity of the buffers
	helpLogger    = log.New(os.Stderr, "", 0) // log to stderr without the timestamps
	closeInterval = 10 * time.Minute          // how often to close our file and open a new one
	byteCount     = counter.NewCounter()      // keep track of how many bytes total have been received
)

type worker struct {
	bytes.Buffer
}

func init() {
	flag.StringVar(&port, "port", "8080", "The port used for the server.")
	flag.StringVar(&pswd, "auth", "", "If not empty, this is matched against the Authorization header (e.g. Authorization: Bearer my-password).")
	flag.IntVar(&requestBuffer, "buf", 500, "The size of the incoming request buffer. A zero (0) will disable buffering.")
	flag.IntVar(&size, "size", 64, "The size (in kilobytes) at which to split the log file(s).")
	flag.BoolVar(&scale, "M", false, "If set, -size will be in megabytes.")
	flag.StringVar(&splitDir, "dir", "", "A dir to use for log files. The first arg after '--' is used as a filename prefix. (e.g. '% omnilogger -dir /path/to/log-dir -- file-prefix-')")
	flag.BoolVar(&help, "h", false, "Show this message.")
	flag.Parse()

	capacity = size * Kilobyte
	if scale {
		capacity = size * Megabyte
	}
	// if capacity == 0 -> stdout

	if help {
		helpLogger.Println("")
		helpLogger.Println("Omnilogger is an HTTP server that coalesces log data (line by line) from multiple sources to a common destination. This defaults to consecutively named log files of ~5000 lines.")
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

	for t := 0; t <= 3; t += 1 {
		workers[t] = &worker{}
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
		case b := <-inbound:
			if w.Len() > capacity {
				go fwrite(w.Bytes())
				w.Reset()
			}
			n, _ := w.Write(b)
			byteCount.IncrBy(uint64(n))
		case <-time.After(closeInterval): // after 10 minutes of inactivity close the file
			go fwrite(w.Bytes())
			w.Reset()
		case <-shutdownCh:
			fwrite(w.Bytes())
			wg.Done()
			return
		}
	}
}

func destructor() {
	wg.Wait()
	helpLogger.Printf(".collected %dk in %s", byteCount.Current(Kilobyte), byteCount.Since())
}
