package main

import (
	"io"
	"testing"
)

func TestNewWriteCloser(t *testing.T) {
	var ok bool

	if _, ok = newWriteCloser(999).(io.WriteCloser); !ok {
		t.Error("Interface error: fallthrough")
	}

	if _, ok = newWriteCloser(ioStdout).(io.WriteCloser); !ok {
		t.Error("Interface error: ioStdout")
	}

	if _, ok = newWriteCloser(ioStderr).(io.WriteCloser); !ok {
		t.Error("Interface error: ioStderr")
	}

	if _, ok = newWriteCloser(ioFile).(io.WriteCloser); !ok {
		t.Error("Interface error: ioFile (lines)")
	}

	splitByteCount = 50
	if _, ok = newWriteCloser(ioFile).(io.WriteCloser); !ok {
		t.Error("Interface error: ioFile (bytes)")
	}
}
