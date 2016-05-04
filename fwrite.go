package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

// a custom error to signal that no file was closed
var (
	ErrNotAFile = errors.New("fwrite: invalid memory address or nil pointer dereference")
	ErrNotADir  = errors.New("fwrite: specified dir is not a dir")
)

func init() {
	if e := checkDir(logDir); e != nil {
		log.Fatal(e)
	}
}

func fwrite(name string, payload []byte) {
	if len(payload) == 0 {
		return
	}

	f, e := create(name)
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
func create(fname string) (io.WriteCloser, error) {
	return os.Create(fname)
}
