package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// a custom error to signal that no file was closed
var (
	ErrNotAFile = errors.New("fwrite: invalid memory address or nil pointer dereference")
	ErrNotADir  = errors.New("fwrite: specified dir is not a dir")
)

func fwrite(payload []byte) {
	if len(payload) == 0 {
		return
	}

	f, e := create("", "TL1-")
	if e != nil {
		log.Println(e)
	}

	f.Write(payload)
	f.Close()
}

// CheckDir ensure that the given dir exists and is a dir
func checkDir(dir string) error {
	dir = filepath.Clean(dir)
	stat, e := os.Stat(dir)
	if os.IsNotExist(e) || !stat.IsDir() || os.IsPermission(e) {
		return ErrNotADir
	}
	return nil
}

// createFile is the file creating function that wraps os.Create
func create(dir, prefix string) (io.WriteCloser, error) {
	var f *os.File
	var e error

	filename := filepath.Join(dir, prefix+time.Now().Format(time.RFC3339Nano))
	f, e = os.Create(filename)
	return f, e
}
