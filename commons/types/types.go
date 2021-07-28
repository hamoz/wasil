package types

import (
	"context"
	"encoding/json"
)

type Message struct {
	ID          string
	Type        MessageType
	ChannelType ChannelType
	ChannelID   string
	Text        string
	Loc         *Location
	Contact     *Contact
	Sender      *User
	Destination *User
	params      map[string]string
}

type MessageType string

// StateType represents an extensible state type in the state machine.
type StateType string

const (
	Default               StateType = ""
	Started               StateType = "Started"
	OriginLocationUpdated StateType = "OriginLocationUpdated"
	DestLocationUpdated   StateType = "DestLocationUpdated"
	Confirmed             StateType = "Confirmed"
	Cancelled             StateType = "Cancelled"
	ContactUpdated        StateType = "ContactUpdated"
	Placed                StateType = "Placed"
)

const (
	TextMessage       = MessageType("Text")
	LocationMessage   = MessageType("Location")
	ContactMessage    = MessageType("Contact")
	StartMessage      = MessageType("Start")
	GetContactMessage = MessageType("GetContact")
	Unknown           = MessageType("Uknown")
	Error             = MessageType("Error")
)

type ChannelType string

const (
	Telegram = ChannelType("Telegram")
	WhatsApp = ChannelType("WhatsApp")
)

type Location struct {
	Name string
	Lat  float32
	Lng  float32
	Dist float32
}

type User struct {
	ID          string
	ChannelType ChannelType
	Phone       string
	Name        string
}

type Contact struct {
	Phone string
	Name  string
}

func (t *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Message) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	return nil
}

type Channel interface {
	Receive(ctx context.Context) []Message
	Handle(ctx context.Context, message Message)
	Send(ctx context.Context, message Message)
	Start()
	Stop()
	Init()
}

type Settings struct {
	Name             string
	ChannelType      ChannelType
	ChannelID        string
	OutQueue         string
	InQueue          string
	AdditionalParams map[string]string
}

func (t *Message) GetPayload() string {
	switch t.Type {
	case TextMessage:
		return t.Text
	case LocationMessage:
		b, err := json.Marshal(t.Loc)
		if err != nil {
			return err.Error()
		}
		return string(b)
	case ContactMessage:
		b, err := json.Marshal(t.Loc)
		if err != nil {
			return err.Error()
		}
		return string(b)
	default:
		return "Unsupported type"
	}
}

func (address *User) GetCompositID() string {
	return string(address.ChannelType) + ":" + address.ID
}
