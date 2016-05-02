package main

import (
	"github.com/henderjon/omnilogger/counter"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	sysSigChan     chan os.Signal
	shutdownLogger = log.New(os.Stderr, "", 0) // log to stderr without the timestamps
	statusInterval = 10 * time.Minute
)

func init() {
	sysSigChan = make(chan os.Signal, 1)
	signal.Notify(sysSigChan, os.Interrupt) // syscall.SIGINT
	signal.Notify(sysSigChan, syscall.SIGTERM)
}

// watchShutdown turn on our signal watching goroutine
func monitorStatus(shutdown chan struct{}) {

	var sig os.Signal
Loop:
	for {
		select {
		case sig = <-sysSigChan:
			close(shutdown) // idiom via: http://dave.cheney.net/2013/04/30/curious-channels
			break Loop
		case <-time.After(statusInterval):
			printStatus()
		}
	}

	shutdownLogger.Printf("\n.signal: %s; shutting down...\n", sig.String())
	wg.Wait()
	printStatus()
	shutdownLogger.Printf(".shutdown: program exit at %s\n", time.Now().Format(time.RFC3339))
	os.Exit(1)
}

// print a status line of total data collected over the life of our server
func printStatus() {
	total := byteCounter.Current(counter.Megabyte)
	shutdownLogger.Printf(".status: collected %dM from %d hits in %s\n", total, hitCounter.Current(uint64(0)), byteCounter.Since().String())
}
