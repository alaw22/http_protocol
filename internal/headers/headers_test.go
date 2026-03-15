package headers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("HOST: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 30, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["context"] = "your mom"
	data = []byte("Host: localhost:42069\r\nTiMe: 15:10:23\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "your mom", headers["context"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[n:])
	assert.Equal(t, "15:10:23", headers["time"])
	assert.Equal(t, 16, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 30, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[n:])
	assert.Equal(t, 0, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	fmt.Println("Error:", err.Error())
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
