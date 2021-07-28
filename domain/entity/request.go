package entity

import (
	"encoding/json"
	"time"

	"github.com/hamoz/wasil/commons/types"
)

type Request struct {
	ID            string
	PassengerID   string
	DriverID      string
	ChannelType   types.ChannelType
	ChannelUserID string
	Time          time.Time
	FromLoc       *types.Location
	ToLocation    *types.Location
	Status        types.StateType
}

type RequestStatus int

const Initial = RequestStatus(0)
const WaitFromLocation = RequestStatus(1)
const WaitToLocation = RequestStatus(2)
const Pricing = RequestStatus(3)
const WaitContact = RequestStatus(4)
const Ready = RequestStatus(9)

func (t *Request) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Request) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	return nil
}

func (t *Request) GetCompsitID() {
}
