package main

import (
	"log"
	"os"
)

var (
	out *log.Logger
)

func init() {
	out = log.New(os.Stdout, "", 0)
}

func main() {
	web(out)
}
