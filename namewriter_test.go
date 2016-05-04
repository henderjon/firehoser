package main

import (
	"testing"
	"io"
)

func Test_ioWriter(t *testing.T) {

	// do we fulfill io.Writer?
	var _ io.Writer = (*nameWriter)(nil)

}

func Test_name(t *testing.T) {

	var nw = &nameWriter{
		"dir", "prefix",
	}

	x := nw.name()

	expected := "dir/prefix"
	if x[:10] != expected {
		t.Error("Response error: \nexpected\n", expected, "\nactual\n", x[:10])
	}
}
