package bytecounter

import (
	"time"
	"errors"
)

// the uint64 representation of a Kilobyte, Megabyte, Gigabyte
const (
	Kilobyte uint64 = 1024
	Megabyte uint64 = Kilobyte * Kilobyte
	Gigabyte uint64 = Kilobyte * Megabyte
)

var (
	bytes = uint64(0)
	created = time.Now()
)

// IncrBy collects a number to add to the total. Presumably the total number of bytes
func IncrBy(count chan int) {
	go func(){
		for i := range count {
			bytes += uint64(i)
		}
	}()
}

// Current returns the total number of bytes collected according to the given meter
func Current(meter uint64) (uint64, error) {
	switch {
	case meter == Kilobyte:
		fallthrough
	case meter == Megabyte:
		fallthrough
	case meter == Gigabyte:
		return (bytes / meter), nil
	}
	return 0, errors.New("not a valid meter")
}

// Since returns the age of the byte collector
func Since() time.Duration {
	return time.Since(created)
}
