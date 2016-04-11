package main

import (
	"flag"
	ws "github.com/henderjon/omnilogger/writesplitter"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	port           string                      // the port on which to listen
	pswd           string                      // a simple means of authentication
	forceStdout    bool                        // skip disk io and allow output redirection
	reqBuffer      int                         // the size of the incoming request buffer (channel)
	splitLineCount int                         // how many lines per log file
	splitByteCount int                         // how many bytes per log file
	splitPrefix    string                      // the prefix for the log file(s) name
	help           bool                        // I forgot my options
	helpLogger     = log.New(os.Stderr, "", 0) // log to stderr without the timestamps
	totalBytes     uint64                      // 9223372036854775806
	closeInterval  = 10 * time.Minute          // how often to close our file and open a new one
)

func init() {
	flag.StringVar(&port, "port", "8080", "The port used for the server.")
	flag.StringVar(&pswd, "auth", "", "If not empty, this is matched against the Authorization header (e.g. Authorization: Bearer my-password).")
	flag.BoolVar(&forceStdout, "c", false, "Write to stdout; disregard -l, -b, and -prefix.")
	flag.IntVar(&reqBuffer, "buf", 0, "The size of the incoming request buffer. A zero (0) will disable buffering.")
	flag.IntVar(&splitLineCount, "l", 5000, "The number of lines at which to split the log files. A zero (0) will disable splitting by lines.")
	flag.IntVar(&splitByteCount, "b", 0, "The number of bytes at which to split the log files. A zero (0) will disable splitting by bytes.")
	flag.StringVar(&splitPrefix, "prefix", "", "A custom prefix to use for log files.")
	flag.BoolVar(&help, "h", false, "Show this message.")
	flag.Parse()

	if help {
		helpLogger.Println("\nOmnilogger is an HTTP server that coalesces log data (line by line) from multiple sources to a common destination. This defaults to consecutively named log files of ~5000 lines. After 10 min of inactivity, the current log file is closed and a new one is opened upon the next write.\n")
		flag.PrintDefaults()
		helpLogger.Println("\n")
		os.Exit(0)
	}

	if err := ws.TestFileIO(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	data := make(chan []byte, reqBuffer)
	wc := newWriteCloser(nil)
	shutdown := make(chan struct{}, 0)

	go watchShutdown(shutdown) // catch system signals and shutdown gracefully
	go watchStatus()           // periodically echo the total number of bytes collects and the duration of the program
	go coalesce(data, wc)      // send all our request data to a single WriteCloser

	// adapters are closures and therefore executed in reverse order
	http.Handle("/", Adapt(parseRequest(data), parseCustomHeader, checkAuth(), ensurePost(), checkShutdown(shutdown)))
	http.ListenAndServe(":"+port, nil)

}

// coalesce runs in it's own goroutine and ranges over a channel and writing the
// data to an io.Writer. All our goroutines send data on this channel and this
// func coalesces them in to one stream.
func coalesce(data chan []byte, wc io.WriteCloser) {
	for {
		select {
		case b := <-data:
			n, _ := wc.Write(append(b, '\n'))
			totalBytes += uint64(n)
		case <-time.After(closeInterval): // after 10 minutes of inactivity close the file
			wc = newWriteCloser(wc)
		}
	}
}

// newWriteCloser is a factory for various io.WriteClosers.
func newWriteCloser(wc io.WriteCloser) io.WriteCloser {
	if forceStdout {
		return os.Stdout // don't ever worry about recycling stdout
	}

	// recycle the old io.WriteCloser with a new WriteSplitter, an explicit Close is better than waiting for GC
	if wc != nil {
		if e := wc.Close(); e != ws.ErrNotAFile {
			log.Fatal(e)
		}
	}

	if splitByteCount > 0 {
		wc = ws.ByteSplitter(splitByteCount, splitPrefix)
	} else {
		wc = ws.LineSplitter(splitLineCount, splitPrefix)
	}
	return wc
}
