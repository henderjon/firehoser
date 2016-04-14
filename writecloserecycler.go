package main

import (
	ws "github.com/henderjon/omnilogger/writesplitter"
	"io"
	"log"
	"os"
)

type writeCloserRecycler func(io.WriteCloser) io.WriteCloser

// newWriteCloser is a factory for various io.WriteClosers.
func writeCloser(dir, prefix string) writeCloserRecycler {
	e := ws.CheckDir(dir)
	if e != nil {
		log.Fatal(e, dir)
	}
	return func(wc io.WriteCloser) io.WriteCloser {
		if forceStdout {
			return os.Stdout // don't ever worry about recycling stdout
		}

		// recycle the old io.WriteCloser with a new WriteSplitter, an explicit Close is better than waiting for GC
		if wc != nil {
			if e := wc.Close(); e != ws.ErrNotAFile {
				log.Fatal(e)
			}
		}

		if byBytes {
			wc = ws.ByteSplitter(limit, dir, prefix)
		} else {
			wc = ws.LineSplitter(limit, dir, prefix)
		}
		return wc
	}
}
