package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	// If there isn't a CRLF then we haven't gotten enough data
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	// if the CRLF is the first thing encountered in the bytes slice
	// then we have reached the end of the headers section
	if idx == 0 {
		return 0, true, nil
	}

	fieldLineText := string(data[:idx])

	fieldName, fieldValue, found := strings.Cut(fieldLineText, ":")

	if !found {
		return 0, false, fmt.Errorf("incorrect field line format, missing ':' -> '%s'\n", fieldLineText)
	}

	if strings.ContainsAny(fieldName, " \t\r\n") {
		return 0, false, fmt.Errorf("incorrect field line format, field-name has whitespace: '%s'\n", fieldName)
	}

	trimmedFieldValue := strings.TrimSpace(fieldValue)

	h[fieldName] = trimmedFieldValue

	return idx + 2, false, nil
}

func NewHeaders() Headers {
	return Headers{}
}
