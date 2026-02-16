package main

import (
		// "os"
		"log"
		"fmt"
		"io"
		"errors"
		"strings"
		"net"
)

const inputFilePath = "messages.txt"
const port = ":42069"

func getLinesChannel(conn io.ReadCloser) <-chan string {
	// Create channel
	lines := make(chan string)
	
	go func(){
		// Create bytes slice 
		b := make([]byte, 8, 8)

		// Close context managers and channels at return
		defer conn.Close()
		defer close(lines)

		currentLineContents := ""
		for {
	
			n, err := conn.Read(b)
			if err != nil {
	
				if currentLineContents != "" {
					lines <- currentLineContents
					currentLineContents = ""
				}
	
				if errors.Is(err, io.EOF) {
					return
				}
	
				fmt.Printf("error: %s\n", err.Error())
				return
			}
	
			str := string(b[:n])
			parts := strings.Split(str,"\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- currentLineContents + parts[i]
				currentLineContents = ""
			}
	
			currentLineContents += parts[len(parts)-1]
	
		}
	}()

	return lines
}

func main(){


	// Open TCP Listener
	listener, err := net.Listen("tcp",port)
	if err != nil{
		log.Fatalf("Could not establish TCP listener: %s\n",err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error in listener.Accept(): %s\n",err.Error())
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())
		
		lines := getLinesChannel(conn)

		for line := range lines{
			fmt.Println(line)
		}
		
		fmt.Println("Connection to ", conn.RemoteAddr(),"closed")

	}




}