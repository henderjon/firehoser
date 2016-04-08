package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	ws "github.com/henderjon/omnilogger/writesplitter"
)

var (
	port           string      // the port on which to listen
	pswd           string      // a simple means of authentication
	forceStdout    bool        // skip disk io and allow output redirection
	splitLineCount int         // how many lines per log file
	splitByteCount int         // how many bytes per log file
	splitPrefix    string      // the prefix for the log file(s) name
	help           bool        // I forgot my options
	bareLog        *log.Logger // log to stderr without the timestamps
)

func init() {

	flag.StringVar(&port, "port", "8080", "The port used for the server.")
	flag.StringVar(&pswd, "auth", "", "If not empty, this is matched against the Authorization header (e.g. Authorization: Bearer my-password).")
	flag.BoolVar(&forceStdout, "c", false, "Send output to stdout and not disk also disregards -l, -b, and -prefix.")
	flag.IntVar(&splitLineCount, "l", 5000, "The number of lines at which to split the log files. A zero (0) will disable splitting by lines.")
	flag.IntVar(&splitByteCount, "b", 0, "The number of bytes at which to split the log files. A zero (0) will disable splitting by bytes.")
	flag.StringVar(&splitPrefix, "prefix", "", "A custom prefix to use for log files.")
	flag.BoolVar(&help, "h", false, "Show this message.")
	flag.Parse()

	bareLog = log.New(os.Stderr, "", 0)

	if help {
		bareLog.Println("\nOmnilogger is an HTTP server that coalesces log data (line by line) from multiple sources to a common destination. This defaults to consecutively named log files of ~5000 lines.\n")
		flag.PrintDefaults()
		bareLog.Println("\n")
		os.Exit(0)
	}

	if err := ws.TestFileIO(); err != nil {
		log.Fatal(err)
	}

	initShutdownWatcher()
}

func main() {
	var wr io.Writer
	ch := make(chan []byte, 0)

	if forceStdout {
		wr = newWriteCloser(ioStdout)
	} else {
		wr = newWriteCloser(ioFile)
	}

	go coalesce(ch, wr)

	// adapters are closures and therefore executed in reverse order
	http.Handle("/", Adapt(parseRequest(ch), parseCustomHeader, checkAuth(), ensurePost(), checkShutdown()))
	http.ListenAndServe(":"+port, nil)

}

// coalesce runs in it's own goroutine and ranges over a channel and writing the
// data to an io.Writer. All our goroutines send data on this channel and this
// func coalesces them in to one stream.
func coalesce(ch chan []byte, wr io.Writer) {
	for b := range ch {
		wr.Write(append(b, '\n'))
	}
}
