package main

import (
	// "log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	sysSigChan chan os.Signal
)

func init() {
	sysSigChan = make(chan os.Signal, 1)
	signal.Notify(sysSigChan, os.Interrupt) // syscall.SIGINT
	signal.Notify(sysSigChan, syscall.SIGTERM)
}

// watchShutdown turn on our signal watching goroutine
func watchShutdown(shutdown chan struct{}) {
	sig := <-sysSigChan
	close(shutdown) // idiom via: http://dave.cheney.net/2013/04/30/curious-channels

	helpLogger.Printf("\n.signal: %s (%d); shutting down...\n", sig.String(), sig)
	wg.Wait()
	printStatus()
	helpLogger.Printf(".shutdown: program exit at %s\n", time.Now().Format(time.RFC3339))
	os.Exit(1)
}
