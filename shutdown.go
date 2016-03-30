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
var swapFile *log.Logger

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
		log.Printf("\n!Caught signal %d; shutting down; flushing to disk\n", sig)

		if sig == syscall.SIGPIPE {
			// atomicly indicate we are in shutdown mode.
			atomic.AddInt32(&brokenPipeSig, 1)
		}

		// atomicly indicate we are in shutdown mode.
		atomic.AddInt32(&shutdownSig, 1)

		wg.Wait()

		log.Println("... fin")
		os.Exit(1)

	}()
}

func getSwapFile() *log.Logger {
	if swapFile != nil {
		return swapFile
	}

	fn := time.Now().Format("2006-01-02T15.04.05Z0700.omnilog.swp")
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	swapFile = log.New(f, "", 0)
	return swapFile
}
