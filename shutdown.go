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
	shutdownSig   int32 // atomically signal shutdow
	brokenPipeSig int32 // atomically signal a broken pipe
)

// isShutdownMode checks to see if we're shutting down
func isShutdownMode() bool {
	s := atomic.LoadInt32(&shutdownSig)
	return s != 0
}

// isBrokenPipe checks to see if we've suffered a broken pipe
func isBrokenPipe() bool {
	s := atomic.LoadInt32(&brokenPipeSig)
	return s != 0
}

// initShutdownWatcher turn on our signal watching goroutine
func initShutdownWatcher() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGPIPE)

	go func() {
		sig := <-c
		bareLog.Printf("\n!Caught signal '%d: %s'; shutting down...\n", sig, sig.String())

		if sig == syscall.SIGPIPE {
			// atomicly indicate we are in shutdown mode.
			atomic.AddInt32(&brokenPipeSig, 1)
		}

		// atomicly indicate we are in shutdown mode.
		atomic.AddInt32(&shutdownSig, 1)

		wg.Wait()

		bareLog.Printf("\nshutdown finished at %s\n", time.Now().Format(time.RFC3339))
		os.Exit(1)
	}()
}
