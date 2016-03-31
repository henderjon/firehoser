package main

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// shutdown marks that the application is in shutdown mode.
var shutdownSig int32
var brokenPipeSig int32

func isShutdownMode() bool {
	s := atomic.LoadInt32(&shutdownSig)
	return s != 0
}

func isBrokenPipe() bool {
	s := atomic.LoadInt32(&brokenPipeSig)
	return s != 0
}

func initShutdownWatcher() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGPIPE)

	go func() {
		sig := <-c
		log.Printf("\t!Caught signal '%d: %s'; shutting down;\n", sig, sig.String())

		if sig == syscall.SIGPIPE {
			// atomicly indicate we are in shutdown mode.
			atomic.AddInt32(&brokenPipeSig, 1)
		}

		// atomicly indicate we are in shutdown mode.
		atomic.AddInt32(&shutdownSig, 1)

		wg.Wait()

		log.Println("... finished at", time.Now().Format(time.RFC3339))
		os.Exit(1)
	}()
}
