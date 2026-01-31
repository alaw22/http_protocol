package main

import (
		"os"
		"log"
		"fmt"
		"io"
		"errors"
)

const inputFilePath = "messages.txt"

func main(){

	// Create bytes slice 
	b := make([]byte, 8, 8)

	// Open messages file
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open '%s': %s\n", inputFilePath, err)
	}
	defer f.Close()

	for {

		n, err := f.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Printf("error: %s\n", err.Error())
			break
		}

		str := string(b[:n])
		fmt.Printf("read: %s\n",str)

	}
}