package main

import (
	"time"
)

const (
	KiloByte = 1024
	MegaByte = KiloByte * KiloByte
	GigaByte = MegaByte * KiloByte
)

var firstPing time.Time = time.Now()

func watchStatus() {
	for {
		printStatus()
		time.Sleep(5 * time.Second)
	}
}

func printStatus() {
	helpLogger.Printf(".status: collected %dK in %s.\n", totalBytes/KiloByte, time.Since(firstPing).String())
}
