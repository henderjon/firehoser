package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

type nameWriter struct {
	dir, prefix string
}

func (nw nameWriter) name() string {
	if len(nw.prefix) == 0 {
		nw.prefix = defaultPrefix
	}
	return filepath.Join(nw.dir, nw.prefix+time.Now().Format(time.RFC3339Nano))
}

func (nw nameWriter) Write(b []byte) (int, error) {
	var e error

	if len(b) > 1 {
		e = ioutil.WriteFile(nw.name(), b, defaultPerms)
		if e != nil {
			log.Println(e)
		}
	}

	return len(b), e
}
