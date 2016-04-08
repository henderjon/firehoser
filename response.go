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
func (r response) Json() string {
	return fmt.Sprintf(`{"Status":%d,"Message":"%s","Received":%d}`, r.Status, r.Message, r.Received)
}

// satisfies the encoding/json.Marshaler interface
func (r response) MarshalJSON() ([]byte, error) {
	return []byte(r.Json()), nil
}
