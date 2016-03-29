package main

import (
	"log"
	"os"
)

var (
	out *log.Logger
)

func init() {
	// use log because fmt isn't goroutine safe
	out = log.New(os.Stdout, "", 0)
}

func main() {
	web(out)
}
