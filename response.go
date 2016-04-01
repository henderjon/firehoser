package main

import "fmt"

const (
	success             = "success"
	errShutdown         = "err: server is shutting down"
	errBadRequest       = "err: missing header"
	errMethodNotAllowed = "err: method not allowed"
	errForbidden        = "err: forbidden"
)

// defines the properties of a response
type response struct {
	Status   int    // the HTTP status code
	Message  string // a descriptive message for the user
	Received int    // the number of bytes received (excludes new lines)
}

// returns the response as a TSV string
func (r response) String() string {
	return fmt.Sprintf("Status\t%s\tMessage\t%s\tReceived\t%d\n", r.Message, r.Received)
}

// returns the []bytes version of the TSV string
func (r response) Bytes() []byte {
	return []byte(r.String())
}
