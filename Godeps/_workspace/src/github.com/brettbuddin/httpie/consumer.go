package httpie

import (
	"bufio"
	"bytes"
)

var (
	NewLine        = Delimeter{'\n'}
	CarriageReturn = Delimeter{'\r'}
	Space          = Delimeter{' '}
	Comma          = Delimeter{','}
)

// Consumer is implemented to provide a method
// of breaking up and consuming a stream
type Consumer interface {
	Consume(*bufio.Reader) ([]byte, error)
}

type Delimeter struct {
	Delim byte
}

// Consume reads up to the next delimeter
func (d Delimeter) Consume(reader *bufio.Reader) ([]byte, error) {
	b, err := reader.ReadBytes(d.Delim)

	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(b), nil
}
