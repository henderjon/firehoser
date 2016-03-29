package main

import (
	"flag"
	"log"
	"os"
)

var (
	out            *log.Logger
	protocol, port string
	help           bool
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
		web(out, port)
	case "tcp":
		sock(out, port)
	default:
		log.Fatal("flag err: `protocol` must be either `http` or `tcp`.")
		os.Exit(1)
	}
	// sock(out)
}
