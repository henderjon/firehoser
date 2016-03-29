package main

import "fmt"

type response struct {
	Message           string
	Recieved, Written int
}

func (r response) String() string {
	return fmt.Sprintf("message\t%s\treceived\t%d\twritten\t%d\n", r.Message, r.Recieved, r.Written)
}

func (r response) Bytes() []byte {
	return []byte(r.String())
}
