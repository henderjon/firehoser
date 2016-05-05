package main

import (
	"log"
	"os"
	"text/template"
)

func welcome() {

	measure := "kb"
	if scale {
		measure = "mb"
	}

	t, _ := template.New("welcome").Parse(`
Omnilogger is a simple web server to coalesce data via HTTP POST.

Current settings:
  - Server Port: {{.Port}}
  - Bearer Token: {{.Pswd}}
  - Custom Header: {{.Header}}
  - Request Buffer Size: {{.Buf}}
  - Number of Workers: {{.Workers}}
  - Worker Size: {{.Size}}{{.Scale}}
  - Log Dir: {{.Dir}}

`)

	if e := t.Execute(os.Stderr, struct {
		Port    string
		Pswd    string
		Header  string
		Buf     int
		Workers int
		Size    int
		Scale   string
		Dir     string
	}{
		Port:    port,
		Pswd:    pswd,
		Header:  customHeader,
		Buf:     requestBuffer,
		Workers: numWorkers,
		Size:    size,
		Scale:   measure,
		Dir:     logDir,
	}); e != nil {
		log.Fatal(e)
	}
}
