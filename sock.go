package main

import (
	"bufio"
	"log"
	"net"
)

// was ripped from the examples as a TCP listener
func sock(out *log.Logger) {
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go HandleSock(conn, out)
	}
}

// scans the incoming data from the connection and writes it to stdout
func HandleSock(conn net.Conn, out *log.Logger) {

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// message = bytes.TrimRight(message, "\n")
		out.Println(scanner.Text())

		if err := scanner.Err(); err != nil {
			break
		}

		conn.Write((&response{
			"success", len(scanner.Text()), len(scanner.Text()),
		}).Bytes())
	}

	// Shut down the connection.
	conn.Close()
}
