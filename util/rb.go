package util

import (
	"strings"
)

// Buffer implements a very simple structure which contains at most size elements.  Every time a new element is added,
// the head of the buffer is evicted.
type Buffer struct {
	Values [] string
}

func (Buffer) New(size int) Buffer {
	return Buffer{Values: make([]string, size)}
}

func (buffer * Buffer) Append(value string) * Buffer {
	a := buffer.Values[1:]
	buffer.Values = append(a, value)
	return buffer
}

func (buffer * Buffer) Join() string {
	var str strings.Builder

	for i := 0; i < len(buffer.Values); i++ {
		str.WriteString(buffer.Values[i])
		str.WriteString("\n")
	}

	return str.String()
}