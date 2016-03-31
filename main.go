package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"
)

const (
	ErrUnknownProtocol = "err: 'protocol' must be either 'http' or 'tcp'"
)

var (
	out            io.Writer
	protocol, port string
	help           bool
	wg             sync.WaitGroup
)

func init() {

	flag.StringVar(&port, "port", "8080", "the port used for the server")
	flag.StringVar(&protocol, "protocol", "http", "the protocol used for the server")
	flag.BoolVar(&help, "h", false, "show this message")
	flag.BoolVar(&help, "h", false, "show this message")

	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	initShutdownWatcher()
}

func main() {

	out = getDest(ioFile)

	switch protocol {
	case "http":
		web(out, port)
	case "tcp":
		sock(out, port)
	default:
		log.Fatalln(ErrUnknownProtocol)
		os.Exit(1)
	}
}
