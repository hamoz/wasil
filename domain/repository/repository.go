package repository

import (
	"github.com/hamoz/wasil/commons/types"
	"github.com/hamoz/wasil/domain/entity"
)

type DriverRepository interface {
	Create(driver *entity.Driver) (*entity.Driver, error)
	Find(id string) (*entity.Driver, error)
	Update(driver *entity.Driver) error
	Delete(key string) error
}

type LocationRepository interface {
	Add(location *types.Location) error
	Find(id string) (*types.Location, error)
	Update(driver *types.Location) error
	Delete(id string) error
	FindByDistance(center *types.Location, count int, distance float64) []types.Location
}

type PassengerRepository interface {
	Create(passenger *entity.Passenger) (*entity.Passenger, error)
	Find(id string) (*entity.Passenger, error)
	FindByChannelUser(channelType types.ChannelType, userId string) (*entity.Passenger, error)
	Update(passenger *entity.Passenger) error
	Delete(id string) error
}

type RequestRepository interface {
	Create(request *entity.Request) (*entity.Request, error)
	Find(id string) (*entity.Request, error)
	FindByChannelUser(channel types.ChannelType, userId string) (*entity.Request, error)
	Update(request *entity.Request) (*entity.Request, error)
	Delete(id string) error
}
