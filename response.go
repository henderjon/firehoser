package main

import (
	"fmt"
	"net/http"
)

// defines the properties of a response
type response struct {
	Status   int    // the HTTP status code
	Message  string // a descriptive message for the user
	Received int    // the number of bytes received (excludes new lines)
}

func newResponse(s, r int) response {
	return response{s, http.StatusText(s), r}
}

// returns a JSON string
func (r response) JSON() string {
	return fmt.Sprint(`{"Status":`, r.Status, `,"Message":"`, r.Message, `","Received":`, r.Received, `}`)
}

// satisfies the encoding/json.Marshaler interface
func (r response) MarshalJSON() ([]byte, error) {
	return []byte(r.JSON()), nil
}
