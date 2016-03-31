package main

import "fmt"

const (
	success             = "success"
	errShutdown         = "err: server is shutting down"
	errBadRequest       = "err: missing header"
	errMethodNotAllowed = "err: method not allowed"
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
