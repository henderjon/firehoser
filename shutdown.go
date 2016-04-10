package main

import (
	// "log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	shutdownSig int32 // atomically signal shutdown
	shutdownCh  chan os.Signal
)

// isShutdownMode checks to see if we're shutting down
func isShutdownMode() bool {
	s := atomic.LoadInt32(&shutdownSig)
	return s != 0
}

func signalShutdown() {
	atomic.AddInt32(&shutdownSig, 1)
}

func init() {
	shutdownCh = make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	signal.Notify(shutdownCh, syscall.SIGTERM)
}

// watchShutdown turn on our signal watching goroutine
func watchShutdown() {
	sig := <-shutdownCh
	helpLogger.Printf("\n.signal: %s (%d); shutting down...\n", sig.String(), sig)

	// atomically indicate we are in shutdown mode.
	signalShutdown()

	wg.Wait()

	printStatus()

	helpLogger.Printf(".shutdown: program exit at %s\n", time.Now().Format(time.RFC3339))
	os.Exit(1)
}
