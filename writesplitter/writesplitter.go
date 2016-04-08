package writesplitter

import (
	"io"
	"os"
	"time"
)

const (
	Kilobyte  = 1024        // const for specifying ByteLimit
	Megabyte  = 1024 * 1024 // const for specifying ByteLimit
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
	LineLimit int            // how many write ops (typically one per line) before splitting the file
	ByteLimit int            // how many bytes before splitting the file
	Prefix    string         // files are named "Prefix + nano-precision-timestamp.log"
	numBytes  int            // internal byte count
	numLines  int            // internal line count
	handle    io.WriteCloser // embedded file
}

// LineSplitter returns a WriteSplitter set to split at the given number of lines
func LineSplitter(limit int, prefix string) io.WriteCloser {
	return &WriteSplitter{LineLimit: limit, Prefix: prefix}
}

// ByteSplitter returns a WriteSplitter set to split at the given number of bytes
func ByteSplitter(limit int, prefix string) io.WriteCloser {
	return &WriteSplitter{ByteLimit: limit, Prefix: prefix}
}

// Close is a passthru and satisfies io.Closer. Subsequent writes will return an
// error.
func (ws *WriteSplitter) Close() error {
	return ws.handle.Close()
}

// Write satisfies io.Writer and internally manages file io. Write also limits
// each WriteSplitter to only one open file at a time.
func (ws *WriteSplitter) Write(p []byte) (int, error) {

	var n int
	var e error

	if ws.handle == nil {
		ws.handle, e = newFile(ws.Prefix)
	}

	switch {
	case ws.LineLimit > 0 && ws.numLines >= ws.LineLimit:
		fallthrough
	case ws.ByteLimit > 0 && ws.numBytes >= ws.ByteLimit:
		ws.Close()
		ws.handle, e = newFile(ws.Prefix)
		ws.numLines, ws.numBytes = 0, 0
	}

	if e != nil {
		return 0, e
	}

	n, e = ws.handle.Write(p)
	ws.numLines += 1
	ws.numBytes += n
	return n, e
}

// newFile creates a new file with the given prefix
func newFile(prefix string) (io.WriteCloser, error) {
	fn := prefix + time.Now().Format(formatStr)
	// fs is an abstraction layer for os allowing us to mock the filesystem for testing
	return fs.Create(fn)
}

// TestFileIO creates and removes a file in the local dir to ensure that it is
// writable. It doesn't use the fs layer because it should be used to test the
// writability of the actual filesystem. This test is unnecessary for mock filesystems
func TestFileIO() error {
	fn := "test.tmp"
	if _, err := os.Create(fn); err != nil {
		return err
	}
	os.Remove(fn)
	return nil
}

/// This is for mocking a filesystem. Used exclusively for testing
///-----------------------------------------------------------------------------

// creator is the interface used to represent the func(s) used for creating a
// file. This abstraction allows us to swap the underlying filesystem with something
// NOT the actual filesystem. useed for testing
type creator interface {
	Create(name string) (file, error)
}

// file is the interface used for our underlying File. os.File has a larger interface
// but within this scope only Write and Close matter
type file interface {
	io.WriteCloser
}

// mocks the os package allowing access to a Create func
var fs creator = ofs{}

// type ofs is a passthru to os.Create()
type ofs struct{}

// satisfies Creator, allowing a passthru to os
func (ofs) Create(name string) (file, error) {
	return os.Create(name)
}
