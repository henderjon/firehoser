package main

import (
	"bufio"
	"net/http"
	"sync"
)

const (
	customHeader = "X-Omnilogger-Stream" // a custom header to validate intent
	methodPost   = "POST"                // because http doesn't have this ...
)

var (
	wg sync.WaitGroup // ensure that our goroutines finish before shut down
)

// Adapter is a decorator that takes a handler and returns a handler.  The
// returned handler does something before calling the handler that was passed in.
// In this way, small pieces of logic can broken into functions that are
// essentially chained together. Casting the returned closures to
// http.HandlerFunc() allows each to satisfy http.Handler. Below, most of the adapters
// are closures themselves. This isn't *necessary* in this instance, but would be
// if they needed to receive/wrap arguments.
type Adapter func(http.Handler) http.Handler

// Adapt takes a group of adapters and runs them in order. The passed handler is
// decorated in each step allowing small pieces of logic to be chained together
// and wrapped around the base handler. Because they are closures the end up
// executing in the reverse order of how they were passed in.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// checkShutdown takes a handler and returns a handler that ensures we're not shutting
// down before calling the passed handler
func checkShutdown(shutdown chan struct{}) Adapter {
	return func(fn http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// graceful shutdown, reject new requests
			select {
			case <-shutdown:
				s := http.StatusServiceUnavailable
				http.Error(rw, (newResponse(s, 0)).JSON(), s)
				return
			default:
			}
			fn.ServeHTTP(rw, req)
		})
	}
}

// ensurePost takes a handler and returns a handler that ensures the current request
// is using the HTTP method POST before calling the passed handler
func ensurePost() Adapter {
	return func(fn http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			s := http.StatusMethodNotAllowed
			// ensure a POST
			if req.Method != methodPost {
				http.Error(rw, (newResponse(s, 0)).JSON(), s)
				return
			}
			fn.ServeHTTP(rw, req)
		})
	}
}

// checkAuth takes a handler and returns a handler that checks the request headers
// for the Authorization header (e.g. 'Authorization: Bearer this-is-a-string')
// and makes sure it matches the given password (if applicable) before calling the
// passed handler
func checkAuth() Adapter {
	return func(fn http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// if no password was given to the server, leave the doors wide open
			if pswd == "" {
				fn.ServeHTTP(rw, req)
				return
			}

			s := http.StatusForbidden

			// make sure the header exists
			a, ok := req.Header["Authorization"]
			if !ok {
				http.Error(rw, (newResponse(s, 0)).JSON(), s)
				return
			}

			// go returns a map, so loop over it
			for _, v := range a {
				// Bearer
				if len(v) < 7 {
					continue
				}
				if v[7:] == pswd {
					fn.ServeHTTP(rw, req)
					return
				}
			}

			http.Error(rw, (newResponse(s, 0)).JSON(), s)
			return
		})
	}
}

// parseCustomHeader takes a handler and returns a handler that checks the request
// headers for the 'X-Omnilogger-Stream' header ... potentially making it's value useful
// before calling the passed handler
func parseCustomHeader(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// must have custom header (@TODO future stream separation?)
		if _, ok := req.Header[customHeader]; !ok {
			s := http.StatusBadRequest
			http.Error(rw, (newResponse(s, 0)).JSON(), s)
			return
		}
		fn.ServeHTTP(rw, req)
	})
}

// parseRequest returns a handler that reads the body of the request and sends
// it on the channel to be coalesced
func parseRequest(data chan *payload) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// if we get here, don't let the program goroutine die before the goroutine finishes
		wg.Add(1)
		defer wg.Done() // cover the short-circuit returns
		defer req.Body.Close()
		s := http.StatusOK
		// read the request body
		rn := 0
		scanner := bufio.NewScanner(req.Body)
		for scanner.Scan() {

			data <- &payload{stream: req.Header[customHeader][0], data: scanner.Bytes()}
			rn += len(scanner.Bytes())

			if err := scanner.Err(); err != nil {
				// using the default scanner, this is most likely an error with the underlying
				// Reader which would most likely indicate an error in the request
				s = http.StatusBadRequest
				http.Error(rw, (newResponse(s, 0)).JSON(), s)
				return
			}
		}
		http.Error(rw, (newResponse(s, rn)).JSON(), s)
	})
}
