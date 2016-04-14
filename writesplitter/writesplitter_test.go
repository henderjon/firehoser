package writesplitter

import (
	"io"
	"os"
	"testing"
)

func TestWriteSplitter(t *testing.T) {
	var expected string
	var O *WriteSplitter
	var e error
	var _ io.WriteCloser = (*WriteSplitter)(nil)

	ltd := 500
	O = LineSplitter(ltd, "/abs/path/", "derp-")
	if e != nil {
		t.Error(e)
	}

	if O.Limit != ltd {
		t.Error("Limit error: \nexpected\n", ltd, "\nactual\n", O.Limit)
	}

	expected = "/abs/path"
	if O.Dir != expected {
		t.Error("Path error (dir): \nexpected\n", expected, "\nactual\n", O.Dir)
	}

	expected = "derp-"
	if O.Prefix != expected {
		t.Error("Path error (prefix): \nexpected\n", expected, "\nactual\n", O.Prefix)
	}

	O = ByteSplitter(ltd, "rel/path/", "derp-")

	expected = "rel/path"
	if O.Dir != expected {
		t.Error("Path error (dir): \nexpected\n", expected, "\nactual\n", O.Dir)
	}

	if O.Bytes != true {
		t.Error("Type error (bytes): \nexpected\n", true, "\nactual\n", O.Bytes)
	}
}

func TestCheckDir(t *testing.T) {
	if e := CheckDir(os.TempDir()); e != nil {
		t.Error("CheckDir error: TempDir not writable")
	}

	if e := CheckDir("."); e != nil {
		t.Error("CheckDir error: PWD not writable")
	}
}
