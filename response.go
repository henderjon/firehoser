package main

import "fmt"

const (
	Success             = "success"
	ErrBadRequest       = "err: missing header"
	ErrMethodNotAllowed = "err: method not allowed"
)

// defines the properties of a response
type response struct {
	Message  string
	Recieved int
}

// returns the response as a TSV string
func (r response) String() string {
	return fmt.Sprintf("message\t%s\tbytes_received\t%d\n", r.Message, r.Recieved)
}

// returns the []bytes version of the TSV string
func (r response) Bytes() []byte {
	return []byte(r.String())
}
