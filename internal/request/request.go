package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/alaw22/http_protocol/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateDone
)

const crlf = "\r\n"
const bufferSize = 8

func (r *Request) parseSingle(data []byte) (int, error) {

	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)

		if err != nil {
			return 0, err
		}

		// if we get passed the top part then err == nil
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders

		return n, nil

	case requestStateParsingHeaders:
		// parse headers
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			r.state = requestStateDone
		}

		return n, nil

	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")

	default:
		return 0, fmt.Errorf("error: unknown request state")
	}

}

func (r *Request) parse(data []byte) (int, error) {

	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}

		totalBytesParsed += n

		if n == 0 {
			break
		}
	}

	return totalBytesParsed, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	rl := &RequestLine{}

	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:idx])

	requestParts := strings.Split(requestLineText, " ")

	if len(requestParts) != 3 {
		return nil, 0, fmt.Errorf("Request line should have 3 items delimited by a space: %v\n", requestParts)
	}

	// Check that the request method is all capital letters
	method_bytes := []byte(requestParts[0])
	for _, char := range method_bytes {
		if char < 65 || char > 90 {
			return nil, 0, fmt.Errorf("HTTP Method contains either a non-capital letter or non-letter '%s'\n", rl.Method)
		}
	}

	httpversionParts := strings.Split(requestParts[2], "/")
	// Check that httpversion split on / has 2 parts
	if len(httpversionParts) != 2 {
		return nil, 0, fmt.Errorf("malformed start-line: '%s'\n", requestLineText)
	}

	// n, Check that the protocol is HTTP
	if httpversionParts[0] != "HTTP" {
		return nil, 0, fmt.Errorf("Unrecognized HTTP version: '%s'\n", requestParts[2])
	}

	// Check that the version of the HTTP protocol is correct
	if httpversionParts[1] != "1.1" {
		return nil, 0, fmt.Errorf("Unrecognized HTTP version: '%s'\n", requestParts[2])
	}

	rl.Method = requestParts[0]
	rl.RequestTarget = requestParts[1]
	rl.HttpVersion = httpversionParts[1]

	return rl, idx + 2, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	bytes_buffer := make([]byte, bufferSize)

	readToIndex := 0 // keeps track of how much data we've read from the `io.Reader` into the buffer

	for r.state != requestStateDone {
		if readToIndex >= len(bytes_buffer) {
			temp := make([]byte, len(bytes_buffer)*2)
			copy(temp, bytes_buffer)
			bytes_buffer = temp
		}

		// read from buffer starting at readToIndex
		numBytesRead, err := reader.Read(bytes_buffer[readToIndex:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				if r.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", r.state, numBytesRead)

				}
				break
			}
			return nil, err
		}

		readToIndex += numBytesRead

		numBytesParsed, err := r.parse(bytes_buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(bytes_buffer, bytes_buffer[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return r, nil
}
