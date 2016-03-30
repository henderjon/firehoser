package main

import (
	"flag"
	"log"
	"os"
	"sync"
)

const (
	ErrUnknownProtocol = "err: 'protocol' must be either 'http' or 'tcp'."
)

var (
	out            *log.Logger
	protocol, port string
	help           bool
	wg             sync.WaitGroup
)

func init() {

	flag.StringVar(&port, "port", "8080", "the port used for the server")
	flag.StringVar(&protocol, "protocol", "http", "the protocol used for the server")
	flag.BoolVar(&help, "h", false, "show this message")

	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// use log because fmt isn't goroutine safe
	out = log.New(os.Stdout, "", 0)
}

func main() {
	switch protocol {
	case "http":
		initShutdownWatcher()
		web(out, port)
	case "tcp":
		initShutdownWatcher()
		sock(out, port)
	default:
		log.Fatalln(ErrUnknownProtocol)
		os.Exit(1)
	}
	// sock(out)
}
