package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"
)

var (
	out            io.Writer      // where to write the output
	wg             sync.WaitGroup // ensure that our goroutines finish before shut down
	tcp            bool           // use tcp as opposed to http
	port           string         // the port on which to listen
	forceStdout    bool           // skip disk io and allow output redirection
	splitLineCount int            // how many lines per log file
	splitByteCount int            // how many bytes per log file
	splitPrefix    string         // the prefix for the log file(s) name
	help           bool           // I forgot my options
)

func init() {

	flag.StringVar(&port, "port", "8080", "The port used for the server.")
	flag.BoolVar(&tcp, "tcp", false, "Use TCP as opposed to (the default) HTTP.")
	flag.BoolVar(&forceStdout, "c", false, "Send output to stdout and not disk also disregards -l, -b, and -prefix.")
	flag.IntVar(&splitLineCount, "l", 5000, "The number of lines at which to split the log files. A zero (0) will disable splitting by lines.")
	flag.IntVar(&splitByteCount, "b", 0, "The number of bytes at which to split the log files. A zero (0) will disable splitting by bytes.")
	flag.StringVar(&splitPrefix, "prefix", "", "The prefix to use for log files.")
	flag.BoolVar(&help, "h", false, "Show this message.")
	flag.Parse()

	if help {
		qLog := log.New(os.Stderr, "", 0)
		qLog.Println("\nOmnilogger is an HTTP (or TCP) server that coalesces log data (line by line) from multiple sources to a common destination (defaults to consecutively named log files of ~5000 lines).\n")
		flag.PrintDefaults()
		qLog.Println("\n")
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

	if tcp {
		sock(out, port)
	} else {
		web(out, port)
	}
}
