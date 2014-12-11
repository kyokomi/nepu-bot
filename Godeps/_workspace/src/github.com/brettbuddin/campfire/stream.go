package campfire

import (
	"encoding/json"
	"fmt"
	"github.com/brettbuddin/httpie"
	"net/url"
)

type Stream struct {
	base     *httpie.Stream
	room     *Room
	outgoing chan *Message
	stop     chan bool
}

// NewStream returns a Stream
func NewStream(room *Room) *Stream {
	return &Stream{
		room:     room,
		outgoing: make(chan *Message),
		stop:     make(chan bool),
	}
}

// Connect starts the stream
func (s *Stream) Connect() {
	url := &url.URL{
		Scheme: "https",
		Host:   "streaming.campfirenow.com",
		Path:   fmt.Sprintf("/room/%d/live.json", s.room.ID),
	}

	s.base = httpie.NewStream(
		httpie.Get{url},
		httpie.BasicAuth{s.room.Connection.Token, "X"},
		httpie.CarriageReturn,
	)

	go s.base.Connect()

	for {
		select {
		case <-s.stop:
			close(s.outgoing)
			return
		case data := <-s.base.Data():
			var m Message
			err := json.Unmarshal(data, &m)

			if err != nil {
				continue
			}

			m.Connection = s.room.Connection
			s.outgoing <- &m
		}
	}
}

func (s *Stream) Messages() chan *Message {
	return s.outgoing
}

// Disconnect stops the stream
func (s *Stream) Disconnect() {
	s.stop <- true
	s.base.Disconnect()
}
