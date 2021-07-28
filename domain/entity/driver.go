package entity

import "encoding/json"

type Driver struct {
	ID            string
	Name          string
	Phone         string
	Channel       string
	ChannelId     string
	Registered    bool
	IdentityType  string
	IdentiyNumber string
	CarMark       string
	CarModel      int
}

func (t *Driver) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Driver) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	return nil
}
