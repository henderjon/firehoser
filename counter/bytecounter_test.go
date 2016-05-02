package counter

import (
	"testing"
)

func TestCounter(t *testing.T) {
	inbound := NewCounter()

	for i := 0; i <= int(Kilobyte); i += 1 {
		inbound.IncrBy(uint64(1))
	}

	expected := uint64(1)
	if inbound.Current(Kilobyte) != expected {
		t.Error("Counting error: \nexpected\n", expected, "\nactual\n", inbound.Current(Kilobyte))
	}
}
