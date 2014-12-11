package campfire

import (
	"fmt"
	"strings"
)

type Room struct {
	*Connection `json:"-"`

	ID              int     `json:"id,omitempty"`
	Full            bool    `json:"full,omitempty"`
	MembershipLimit int     `json:"membership_limit,omitempty"`
	Name            string  `json:"name,omitempty"`
	OpenToGuests    bool    `json:"open_to_guests,omitempty"`
	Topic           string  `json:"topic,omitempty"`
	Users           []*User `json:"users,omitempty"`
	UpdatedAt       string  `json:"updated_at,omitempty"`
	CreatedAt       string  `json:"created_at,omitempty"`
}

type RoomResult struct {
	Room *Room `json:"room"`
}

type RoomsResult struct {
	Rooms []*Room `json:"rooms"`
}

// Stream returns a Stream for you to follow the contents of the Room
func (r *Room) Stream() *Stream {
	return NewStream(r)
}

// Join joins the Room
func (r *Room) Join() error {
	return r.Connection.Post(fmt.Sprintf("/room/%d/join", r.ID), nil)
}

// Leave leaves the Room
func (r *Room) Leave() error {
	return r.Connection.Post(fmt.Sprintf("/room/%d/leave", r.ID), nil)
}

// Lock locks the Room
func (r *Room) Lock() error {
	return r.Connection.Post(fmt.Sprintf("/room/%d/lock", r.ID), nil)
}

// Unlock unlocks the Room
func (r *Room) Unlock() error {
	return r.Connection.Post(fmt.Sprintf("/room/%d/unlock", r.ID), nil)
}

// SendText sends a TextMessage to the Room
func (r *Room) SendText(message string) error {
	message = strings.Replace(message, "\n", "&#xA;", -1)
	return r.message(&Message{Type: "TextMessage", Body: message})
}

// SendPaste sends a PasteMessage to the Room
func (r *Room) SendPaste(content string) error {
	return r.message(&Message{Type: "PasteMessage", Body: content})
}

// SendSound sends a SoundMessage to the Room
func (r *Room) SendSound(name string) error {
	return r.message(&Message{Type: "SoundMessage", Body: name})
}

// SendTweet sends a TweetMessage to the Room
func (r *Room) SendTweet(url string) error {
	return r.message(&Message{Type: "TweetMessage", Body: url})
}

func (r *Room) message(m *Message) error {
	result := MessageResult{m}
	return r.Connection.Post(fmt.Sprintf("/room/%d/speak", r.ID), result)
}
