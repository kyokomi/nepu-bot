package httpie

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// NewStream returns a Stream
func NewStream(endpoint Endpoint, auth Authorizer, consumer Consumer) *Stream {
	return &Stream{
		endpoint:   endpoint,
		authorizer: auth,
		consumer:   consumer,
		data:       make(chan []byte, 50),
		errors:     make(chan error, 50),
		stop:       make(chan bool),
	}
}

type Stream struct {
	data       chan []byte
	errors     chan error
	stop       chan bool
	endpoint   Endpoint
	authorizer Authorizer
	consumer   Consumer
}

// Connect starts the stream
func (s *Stream) Connect() {
	resp, err := s.connect()
	if err != nil {
		s.errors <- err
		return
	}

	s.consume(resp)
}

// Data returns a channel that chunks of the
// feed will be communicated upon
func (s *Stream) Data() chan []byte {
	return s.data
}

// Errors returns a channel that stream errors
// will be sent over
func (s *Stream) Errors() chan error {
	return s.errors
}

// Disconnect forcefully disconnects from the stream
func (s *Stream) Disconnect() {
	s.stop <- true
}

func (s *Stream) connect() (*http.Response, error) {
	client := &http.Client{}
	req := &http.Request{Header: http.Header{}}

	s.endpoint.ApplyTo(req)
	if s.authorizer != nil {
		s.authorizer.Authorize(req)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status code received: %s", resp.StatusCode))
	}

	return resp, nil
}

func (s *Stream) consume(resp *http.Response) {
	reader := bufio.NewReader(resp.Body)

	var (
		b   []byte
		err error
	)

	for {
		select {
		case <-s.stop:
			resp.Body.Close()
			close(s.stop)
			close(s.errors)
			close(s.data)
			return
		default:
			b, err = s.consumer.Consume(reader)

			if err != nil {
				resp.Body.Close()
				time.Sleep(10 * time.Second)

				if resp, err = s.connect(); err != nil {
					s.errors <- err
					continue
				}

				reader = bufio.NewReader(resp.Body)
			}

			s.data <- b
		}
	}
}
