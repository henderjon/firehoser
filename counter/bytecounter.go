package counter

import (
	"time"
)

// the uint64 representation of a Kilobyte, Megabyte, Gigabyte
const (
	Kilobyte uint64 = 1024
	Megabyte uint64 = Kilobyte * Kilobyte
	Gigabyte uint64 = Kilobyte * Megabyte
)

type Counter struct {
	count uint64 // 9223372036854775806
	created time.Time
}

func NewCounter() *Counter {
	return &Counter{
		count: uint64(0),
		created: time.Now(),
	}
}

// IncrBy collects a number to add to the total. Presumably the total number of count
func (c *Counter) IncrBy(i uint64) {
	c.count += uint64(i)
}

// Current returns the total number of count collected according to the given meter
func (c *Counter) Current(meter uint64) uint64 {
	switch {
	case meter == Kilobyte:
		fallthrough
	case meter == Megabyte:
		fallthrough
	case meter == Gigabyte:
		return (c.count / meter)
	default:
		return c.count
	}
}

// Since returns the age of the byte collector
func (c *Counter) Since() time.Duration {
	return time.Since(c.created)
}
