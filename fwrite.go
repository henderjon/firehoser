package main

import(
	"path/filepath"
	"os"
	"time"
	"log"
)

func fwrite(payload []byte) {
	if len(payload) == 0 {
		return
	}
	filename := filepath.Join("TL1-" + time.Now().Format(time.RFC3339Nano))
	f, e := os.Create(filename)
	if e != nil {
		log.Fatal(e)
	}
	defer f.Close()
	f.Write(payload)
}
