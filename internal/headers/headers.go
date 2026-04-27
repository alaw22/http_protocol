package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"
const validCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-.^_`|~"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	// If there isn't a CRLF then we haven't gotten enough data
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	// if the CRLF is the first thing encountered in the bytes slice
	// then we have reached the end of the headers section
	if idx == 0 {
		return 2, true, nil
	}

	fieldLineText := string(data[:idx])

	fieldName, fieldValue, found := strings.Cut(fieldLineText, ":")

	if !found {
		return 0, false, fmt.Errorf("incorrect field line format, missing ':' -> '%s'\n", fieldLineText)
	}

	if strings.ContainsAny(fieldName, " \t\r\n") {
		return 0, false, fmt.Errorf("incorrect field line format, field-name has whitespace: '%s'\n", fieldName)
	}

	// Check that all characters in fieldname are valid
	for _, char := range fieldName {
		if !bytes.ContainsRune([]byte(validCharacters), char) {
			return 0, false, fmt.Errorf("incorrect field line format, field-name has an invalid character: '%s'-->'%c'\n", fieldName, char)
		}
	}

	trimmedFieldValue := strings.TrimSpace(fieldValue)
	lowerFieldName := strings.ToLower(fieldName)

	if val, ok := h[lowerFieldName]; ok {
		h[lowerFieldName] = fmt.Sprintf("%s, %s", val, trimmedFieldValue)
	} else {
		h[lowerFieldName] = trimmedFieldValue
	}

	return idx + 2, false, nil
}

func NewHeaders() Headers {
	return Headers{}
}
