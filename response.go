package main

import "fmt"

// defines the properties of a response
type response struct {
	Message           string
	Recieved, Written int
}

// returns the response as a TSV string
func (r response) String() string {
	return fmt.Sprintf("message\t%s\treceived\t%d\twritten\t%d\n", r.Message, r.Recieved, r.Written)
}

// returns the []bytes version of the TSV string
func (r response) Bytes() []byte {
	return []byte(r.String())
}
