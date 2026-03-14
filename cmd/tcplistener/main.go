package main

import (
	"fmt"
	"log"
	"net"

	"http_protocol/internal/request"
)

const inputFilePath = "messages.txt"
const port = ":42069"

func main() {

	// Open TCP Listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Could not establish TCP listener: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error in listener.Accept(): %s\n", err.Error())
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)

		fmt.Printf("Request line:\n - Method: %s\n - Target: %s\n - Version: %s\n",
			req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion)

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")

	}

}
