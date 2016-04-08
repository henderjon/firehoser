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
)

// isShutdownMode checks to see if we're shutting down
func isShutdownMode() bool {
	s := atomic.LoadInt32(&shutdownSig)
	return s != 0
}

func signalShutdown() {
	atomic.AddInt32(&shutdownSig, 1)
}

// initShutdownWatcher turn on our signal watching goroutine
func initShutdownWatcher() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		sig := <-c
		bareLog.Printf("\n!Caught signal '%d: %s'; shutting down...\n", sig, sig.String())

		// atomically indicate we are in shutdown mode.
		signalShutdown()

		wg.Wait()

		bareLog.Printf("\nshutdown finished at %s\n", time.Now().Format(time.RFC3339))
		os.Exit(1)
	}()
}
