package writesplitter

import (
	"log"
	"os"
	"time"
)

const (
	Kilobyte  = 1024
	Megabyte  = 1024 * 1024
	formatStr = "2006-01-02T15.04.05.999999999Z0700.log"
)

// WriteSplitter represents a disk bound io.WriteCloser that splits the input
// across consecutively named files based on either the number of bytes or the
// number of lines. Splitting does not guarantee true byte/line split
// precision as it does not parse the incoming data. The decision to split is
// before the underlying write operation based on the previous invocation. In
// other words, if a []byte sent to `Write()` contains enough bytes or new
// lines ('\n') to exceed the given limit, a new file won't be generated until
// the *next* invocation of `Write()`. If both LineLimit and ByteLimit is set,
// preference is given to LineLimit. By default, no splitting occurs because
// both LineLimit and ByteLimit are zero (0).
type WriteSplitter struct {
	LineLimit int      // how many write ops (typically one per line) before splitting the file
	ByteLimit int      // how many bytes before splitting the file
	Prefix    string   // files are named "Prefix + nano-precision-timestamp.log"
	numBytes  int      // internal byte count
	numLines  int      // internal line count
	handle    *os.File // embedded file
}

func init() {
	testFileIO()
}

// Close is a passthru and satisfies io.Closer. Subsequent writes will return an
// error.
func (ws *WriteSplitter) Close() error {
	return ws.handle.Close()
}

// Write satisfies io.Writer and internally manages file io. Write also limits
// each WriteSplitter to only one open file at a time.
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

// newFile creates a new file with the given prefix
func newFile(prefix string) (*os.File, error) {
	fn := prefix + time.Now().Format(formatStr)
	return os.Create(fn)
}

// testFileIO creates and removes a file at init() to ensure that the current dir is writable
func testFileIO() {
	fn := "test.tmp"
	_, err := os.Create(fn)
	if err != nil {
		log.Fatal("WriteSplitter cannot write files to this location")
	}
	os.Remove(fn)
}
