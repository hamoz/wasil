package entity

import "encoding/json"

type Passenger struct {
	ID            string
	Name          string
	Phone         string
	ChannelType   string
	ChannelUserID string
	Registered    bool
}

func (t *Passenger) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Passenger) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	return nil
}
