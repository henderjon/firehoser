package bytecounter

import (
	"testing"
)

func TestCounter(t *testing.T) {
	inbound := make(chan int)
	IncrBy(inbound)

	for i := 0; i <= int(Kilobyte); i += 1 {
		inbound <- 1
	}
	close(inbound)

	expected := uint64(1)
	if Current(Kilobyte) != expected {
		t.Error("Counting error: \nexpected\n", expected, "\nactual\n", Current(Kilobyte))
	}
}
