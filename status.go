package main

import (
	"time"
)

const (
	Kilobyte       uint64 = 1024                // the uint64 representation of a Kilobyte
	Megabyte       uint64 = Kilobyte * Kilobyte // the uint64 representation of a Megabyte
	Gigabyte       uint64 = Kilobyte * Megabyte // the uint64 representation of a Gigabyte
	statusInterval        = 10 * time.Minute    // the interval for printing a status line
)

var firstPing time.Time = time.Now()

// every 10 minutes, print a summary of data collected (in KB) and the age of our server
func watchStatus() {
	for {
		printStatus()
		time.Sleep(statusInterval)
	}
}

// print a status line of total data collected over the life of our server
func printStatus() {
	helpLogger.Printf(".status: collected %dK in %s\n", totalBytes/Kilobyte, time.Since(firstPing).String())
}
