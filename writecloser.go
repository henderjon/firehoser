package main

import (
	ws "github.com/henderjon/omnilogger/writesplitter"
	"io"
	"os"
	"time"
)

const (
	ioFile = 1 << iota // constants to label the desired destination
	ioStdout
	ioStderr
)

// getDest is a factory for various log destinations.
func newWriteCloser(t int) io.WriteCloser {
	var writer io.WriteCloser
	switch {
	case t == ioFile:
		writer = &ws.WriteSplitter{
			LineLimit: splitLineCount,
			ByteLimit: splitByteCount,
			Prefix:    splitPrefix,
			Created:   time.Now(),
		}
	case t == ioStderr:
		writer = os.Stderr
	default:
		fallthrough
	case t == ioStdout:
		writer = os.Stdout
	}
	return writer
}
