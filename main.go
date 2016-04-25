package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	port          string                      // the port on which to listen
	pswd          string                      // a simple means of authentication
	forceStdout   bool                        // skip disk io and allow output redirection
	reqBuffer     int                         // the size of the incoming request buffer (channel)
	limit         int                         // how many lines per log file
	byBytes       bool                        // how many bytes per log file
	splitDir      string                      // the dir for the log file(s)
	help          bool                        // I forgot my options
	helpLogger    = log.New(os.Stderr, "", 0) // log to stderr without the timestamps
	totalBytes    uint64                      // 9223372036854775806
	closeInterval = 10 * time.Minute          // how often to close our file and open a new one
	wg            sync.WaitGroup              // ensure that our goroutines finish before shut down
)

// payload wraps the data and the intended stream to transport via channel
type payload struct {
	stream string
	data   []byte
}

func init() {
	flag.StringVar(&port, "port", "8080", "The port used for the server.")
	flag.StringVar(&pswd, "auth", "", "If not empty, this is matched against the Authorization header (e.g. Authorization: Bearer my-password).")
	flag.BoolVar(&forceStdout, "c", false, "Write to stdout; disregard -l, -b, and -prefix. **Not** recommended for production use.")
	flag.IntVar(&reqBuffer, "buf", 0, "The size of the incoming request buffer. A zero (0) will disable buffering.")
	flag.IntVar(&limit, "l", 5000, "The limit at which to split the log files. Assumes a line count. A zero (0) will disable splitting.")
	flag.BoolVar(&byBytes, "bytes", false, "Split according to byte count, not line count.")
	flag.StringVar(&splitDir, "dir", "", "A dir to use for log files. The first arg after '--' is used as a filename prefix. (e.g. '% omnilogger -dir /path/to/log-dir -- file-prefix-')")
	flag.BoolVar(&help, "h", false, "Show this message.")
	flag.Parse()

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
	inbound := make(chan *payload, reqBuffer)
	shutdown := make(chan struct{}, 0)
	byteCount := countBytes() // send the total bytes collected on a channel for periodic output

	go monitorStatus(shutdown)                   // catch system signals and shutdown gracefully
	go coalesce(inbound, byteCount, writeCloser) // send all our request data to a single WriteCloser

	// adapters are closures and therefore executed in reverse order
	http.Handle("/", Adapt(parseRequest(inbound), parseCustomHeader, checkAuth(), ensurePost(), checkShutdown(shutdown)))
	http.ListenAndServe(":"+port, nil)
}

// coalesce runs in it's own goroutine and ranges over a channel and writing the
// data to an io.Writer. All our goroutines send data on this channel and this
// func coalesces them in to one stream.
func coalesce(inbound chan *payload, byteCount chan int, factory writeCloserRecyclerFactory) {
	wcMap := make(map[string]io.WriteCloser, 0)
	for {
		select {
		case b := <-inbound:
			wg.Add(1)
			if _, ok := wcMap[b.stream]; !ok {
				wcMap[b.stream] = factory(splitDir, b.stream)(nil)
			}

			// go func(wr io.WriteCloser) {
			n, _ := wcMap[b.stream].Write(append(b.data, '\n'))
			byteCount <- n
			wg.Done()
			// }(wcMap[b.stream])
			// case <-time.After(closeInterval): // after 10 minutes of inactivity close the file
			// wcMap[b.stream] = wcr(wc)
		}
	}
}
