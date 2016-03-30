package main

import (
	"bufio"
	"log"
	"net"
	"time"
)

var timeout = 3 * time.Second

// was ripped from the examples as a TCP listener
func sock(out *log.Logger, port string) {
	l, err := net.Listen("tcp", ":"+port)
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

		// graceful shutdown does not accept new connections
		if isShutdownMode() {
			conn.Write((&response{
				ErrShutdown, 0,
			}).Bytes())
			log.Println(ErrShutdown)
			return
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleSock(conn, out)
	}
}

// scans the incoming data from the connection and writes it to stdout
func handleSock(conn net.Conn, out *log.Logger) {

	// if we get here, don't let the program goroutine die before the goroutine finishes
	wg.Add(1)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {

		if isBrokenPipe() {
			out = getSwapFile()
		}

		out.Println(scanner.Text())

		if err := scanner.Err(); err != nil {
			log.Println(err.Error())
			break
		}

		conn.Write((&response{
			Success, len(scanner.Text()),
		}).Bytes())

		// let current connections finish writing
		// if isShutdownMode() { conn.Close() }

		// if nothing is happening, close the connection
		conn.SetDeadline(time.Now().Add(timeout))

	}

	wg.Done()

	// Shut down the connection.
	conn.Close()
}
