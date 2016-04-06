package main

import (
	ws "github.com/henderjon/omnilogger/writesplitter"
	"io"
	"log"
	"os"
)

const (
	ioFile = 1 << iota // constants to label the desired destination
	ioStdout
	ioStderr
)

type destination struct {
	*log.Logger
}

// Write Satisfies io.WriteCloser but guaruntees atomicity via log
func (d *destination) Write(s []byte) (int, error) {
	d.Println(string(s))
	return len(s), nil
}

// Write Satisfies io.WriteCloser but guaruntees atomicity via log
func (d *destination) Close() error {
	return d.Close()
}

// getDest is a factory for various log destinations.
func getDest(t int) io.WriteCloser {
	var writer io.WriteCloser
	switch {
	case t == ioFile:
		writer = &ws.WriteSplitter{
			LineLimit: splitLineCount,
			ByteLimit: splitByteCount,
			Prefix:    splitPrefix,
		}
	case t == ioStderr:
		writer = os.Stderr
	default:
		fallthrough
	case t == ioStdout:
		writer = os.Stdout
	}
	return &destination{log.New(writer, "", 0)}
}
