package main

import (
	"log"
	"syscall"
)

func increaseFileLimit(n uint64) {

	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Fatal("Error Getting Rlimit", err)
		return
	}

	if rLimit.Cur < n {
		log.Printf("Increasing maximum number of open files to %d (it was originally set to %d)", n, rLimit.Cur)

		rLimit.Cur = n
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			log.Fatal("Error Setting Rlimit ", err)
		}
	}
}
