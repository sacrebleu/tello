package util

import (
	"strings"
)

type Buffer struct {
	values [] string
}

func (Buffer) New(size int) Buffer {
	return Buffer{values: make([]string, size)}
}

func (buffer * Buffer) Append(value string) * Buffer {
	a := buffer.values[1:]
	buffer.values = append(a, value)
	return buffer
}

func (buffer * Buffer) Join() string {
	var str strings.Builder

	for i := 0; i < len(buffer.values); i++ {
		str.WriteString(buffer.values[i])
		str.WriteString("\n")
	}

	return str.String()
}