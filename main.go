package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"
)

const (
	errUnknownProtocol = "err: 'protocol' must be either 'http' or 'tcp'"
)

var (
	out            io.Writer      // where to write the output
	wg             sync.WaitGroup // ensure that our goroutines finish before shut down
	protocol       string         // http or tcp
	port           string         // the port on which to listen
	forceStdout    bool           // skip disk io and allow output redirection
	splitLineCount int            // how many lines per log file
	splitByteCount int            // how many bytes per log file
	splitPrefix    string         // the prefix for the log file(s) name
	help           bool           // I forgot my options
)

func init() {

	flag.StringVar(&port, "port", "8080", "The port used for the server.")
	flag.StringVar(&protocol, "protocol", "http", "The protocol used for the server.")
	flag.BoolVar(&forceStdout, "c", false, "Send output to stdout and not disk also disregards -l, -b, and -prefix.")
	flag.IntVar(&splitLineCount, "l", 5000, "The number of lines at which to split the log files. A zero (0) will disable splitting by lines.")
	flag.IntVar(&splitByteCount, "b", 0, "The number of bytes at which to split the log files. A zero (0) will disable splitting by bytes.")
	flag.StringVar(&splitPrefix, "prefix", "", "The prefix to use for log files.")
	flag.BoolVar(&help, "h", false, "Show this message.")
	flag.Parse()

	if help {
		log.Println("Omnilogger is an HTTP or TCP server that coalesces log data (line by line) from multiple sources to a common destination (defaults to consecutively named log files of ~5000 lines).")
		flag.PrintDefaults()
		os.Exit(0)
	}

	initShutdownWatcher()
}

func main() {

	if forceStdout {
		out = getDest(ioStdout)
	} else {
		out = getDest(ioFile)
	}

	switch protocol {
	case "http":
		web(out, port)
	case "tcp":
		sock(out, port)
	default:
		log.Fatalln(errUnknownProtocol)
		os.Exit(1)
	}
}
