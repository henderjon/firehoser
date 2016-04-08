package main

import (
	ws "github.com/henderjon/omnilogger/writesplitter"
	"io"
	"os"
)

const (
	ioFile = 1 << iota // constants to label the desired destination
	ioStdout
	ioStderr
)

// newWriteCloser is a factory for various io.WriteClosers.
func newWriteCloser(t int) io.WriteCloser {
	var writer io.WriteCloser
	switch {
	case t == ioFile:
		if splitByteCount > 0 {
			writer = ws.ByteSplitter(splitByteCount, splitPrefix)
		}else{
			writer = ws.LineSplitter(splitLineCount, splitPrefix)
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
