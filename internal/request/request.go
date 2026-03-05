package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parseRequestLine(data []byte) error {
	str := string(data)
	requestline := strings.Split(str, "\r\n")[0]
	requestParts := strings.Split(requestline, " ")

	if len(requestParts) != 3 {
		return fmt.Errorf("Request line should have 3 items delimited by a space: %v\n", requestParts)
	}

	// Check that the request method is all capital letters
	method_bytes := []byte(requestParts[0])
	for _, char := range method_bytes {
		if char < 65 || char > 90 {
			return fmt.Errorf("HTTP Method contains either a non-capital letter or non-letter '%s'\n", r.RequestLine.Method)
		}
	}

	httpversionParts := strings.Split(requestParts[2], "/")
	// Check that httpversion split on / has 2 parts
	if len(httpversionParts) != 2 {
		return fmt.Errorf("malformed start-line: '%s'\n", str)
	}

	// Check that the protocol is HTTP
	if httpversionParts[0] != "HTTP" {
		return fmt.Errorf("Unrecognized HTTP version: '%s'\n", requestParts[2])
	}

	// Check that the version of the HTTP protocol is correct
	if httpversionParts[1] != "1.1" {
		return fmt.Errorf("Unrecognized HTTP version: '%s'\n", requestParts[2])
	}

	r.RequestLine.Method = requestParts[0]
	r.RequestLine.RequestTarget = requestParts[1]
	r.RequestLine.HttpVersion = httpversionParts[1]

	return nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{}

	// read from reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = r.parseRequestLine(data)
	if err != nil {
		return nil, err
	}

	return r, nil
}
