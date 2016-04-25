package main

import (
	ws "github.com/henderjon/omnilogger/writesplitter"
	"io"
	"log"
	"os"
)

// writeCloserRecyclerFactory returns a factory that returns writeCloserRecycler that are bound to the given dir and prefix
type writeCloserRecyclerFactory func(dir, prefix string) writeCloserRecycler

// writeCloserRecycler io.WriteCloser closes and opens a writesplitter bound to the same destination
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
			if e := wc.Close(); e != nil && e != ws.ErrNotAFile {
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
