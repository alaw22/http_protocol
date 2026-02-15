package main

import (
		"os"
		"log"
		"fmt"
		"io"
		"errors"
		"strings"
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

	currentLineContents := ""
	for {

		n, err := f.Read(b)
		if err != nil {

			if currentLineContents != "" {
				fmt.Printf("read: %s\n", currentLineContents)
				currentLineContents = ""
			}

			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Printf("error: %s\n", err.Error())
			break
		}

		str := string(b[:n])
		parts := strings.Split(str,"\n")
		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s%s\n", currentLineContents, parts[i])
			currentLineContents = ""
		}

		currentLineContents += parts[len(parts)-1]

	}


}