package main

import (
	ws "github.com/capdig/omnilogger/writesplitter"
	"io"
	"log"
	"os"
)

const (
	ioFile = 1 << iota
	ioStdout
	ioStderr
)

type destination struct {
	*log.Logger
}

func (d *destination) Write(s []byte) (int, error) {
	d.Println(string(s))
	return len(s), nil
}

func getDest(t int) io.Writer {
	var writer io.Writer
	switch {
	case t == ioFile:
		writer = &ws.WriteSplitter{LineLimit: 5000, Prefix: "test-log-"}
	case t == ioStderr:
		writer = os.Stderr
	default:
		fallthrough
	case t == ioStdout:
		writer = os.Stdout
	}
	return &destination{log.New(writer, "", 0)}
}
