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

func getLinesChannel(f io.ReadCloser) <-chan string {
	// Create channel
	lines := make(chan string)
	
	go func(){
		// Create bytes slice 
		b := make([]byte, 8, 8)

		// Close context managers and channels at return
		defer f.Close()
		defer close(lines)

		currentLineContents := ""
		for {
	
			n, err := f.Read(b)
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



	// Open messages file
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open '%s': %s\n", inputFilePath, err)
	}
	// defer f.Close()

	lines := getLinesChannel(f)

	for line := range lines{
		fmt.Printf("read: %s\n",line)
	}


}