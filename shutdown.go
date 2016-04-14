package main

import (
	bc "github.com/henderjon/omnilogger/bytecounter"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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
func monitorStatus(shutdown chan struct{}) {

	var sig os.Signal
Loop:
	for {
		select {
		case sig = <-sysSigChan:
			close(shutdown) // idiom via: http://dave.cheney.net/2013/04/30/curious-channels
			break Loop
		case <-time.After(10 * time.Minute):
			printStatus()
		}
	}

	shutdownLogger.Printf("\n.signal: %s; shutting down...\n", sig.String())
	wg.Wait()
	printStatus()
	shutdownLogger.Printf(".shutdown: program exit at %s\n", time.Now().Format(time.RFC3339))
	os.Exit(1)
}

func countBytes() chan int {
	byteCount := make(chan int, 0)
	bc.IncrBy(byteCount)
	return byteCount
}

// print a status line of total data collected over the life of our server
func printStatus() {
	total, _ := bc.Current(bc.Kilobyte)
	shutdownLogger.Printf(".status: collected %dK in %s\n", total, bc.Since().String())
}
