package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"sync"
)

var (
	sysSigChan     chan os.Signal
	shutdownLogger = log.New(os.Stderr, "", 0) // log to stderr without the timestamps
)

func init() {
	sysSigChan = make(chan os.Signal, 1)
	signal.Notify(sysSigChan, os.Interrupt) // syscall.SIGINT
	signal.Notify(sysSigChan, syscall.SIGTERM)
}

// watchShutdown turn on our signal watching goroutine
func monitorStatus(shutdown chan struct{}, wg sync.WaitGroup) {

	var sig os.Signal

	select {
	case sig = <-sysSigChan:
		close(shutdown) // idiom via: http://dave.cheney.net/2013/04/30/curious-channels
	}

	shutdownLogger.Printf("\n.signal: %s; shutting down...\n", sig.String())
	wg.Wait()
	shutdownLogger.Printf(".shutdown: program exit at %s\n", time.Now().Format(time.RFC3339))
	os.Exit(1)
}
