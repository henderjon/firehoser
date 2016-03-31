package writesplitter

import (
	"os"
	"time"
	"log"
)

const (
	Kilobyte  = 1024
	Megabyte  = 1024 * 1024
	formatStr = "2006-01-02T15.04.05.999999999Z0700.log"
)

type WriteSplitter struct {
	LineLimit, ByteLimit int
	Prefix               string
	numBytes, numLines   int
	handle               *os.File
}

func init() {
	testFileIO()
}

func (ws *WriteSplitter) Close() error {
	err := ws.handle.Close()
	// @TODO, should I really allow more writes?
	ws.handle = nil
	return err
}

func (ws *WriteSplitter) Write(p []byte) (int, error) {

	var err error

	if ws.handle == nil {
		ws.handle, err = newFile(ws.Prefix)
	}

	switch {
	case ws.LineLimit > 0 && ws.numLines >= ws.LineLimit:
		fallthrough
	case ws.ByteLimit > 0 && ws.numBytes >= ws.ByteLimit:
		ws.Close()
		ws.handle, err = newFile(ws.Prefix)
		ws.numLines, ws.numBytes = 0, 0
	}

	if err != nil {
		return 0, nil
	}

	n, err := ws.handle.Write(p)
	ws.numLines += 1
	ws.numBytes += n
	return n, err
}

func newFile(prefix string) (*os.File, error) {
	fn := prefix + time.Now().Format(formatStr)
	return os.Create(fn)
}

func testFileIO() {
	fn := "test.tmp"
	_, err := os.Create(fn)
	if err != nil {
		log.Fatal("WriteSplitter cannot write files to this location")
	}
	os.Remove(fn)
}
