package service

import (
	"github.com/hamoz/wasil/commons/types"
	"github.com/hamoz/wasil/domain/entity"
	"github.com/hamoz/wasil/domain/repository"
)

type PassengerService interface {
	Register(passenger *entity.Passenger) (*entity.Passenger, error)
	Find(id string) (*entity.Passenger, error)
	FindByChannelUser(channelType types.ChannelType, userId string) (*entity.Passenger, error)
	Update(passenger *entity.Passenger) error
	Delete(id string) error
}

type passengerService struct {
	repo repository.PassengerRepository
}

func NewPassengerService(passengerRepo repository.PassengerRepository) *passengerService {
	return &passengerService{passengerRepo}
}

func (srv *passengerService) Register(passenger *entity.Passenger) (*entity.Passenger, error) {
	return srv.repo.Create(passenger)
}

func (srv *passengerService) Find(id string) (*entity.Passenger, error) {
	return srv.repo.Find(id)
}

func (srv *passengerService) FindByChannelUser(channelType types.ChannelType, userId string) (*entity.Passenger, error) {
	return srv.repo.FindByChannelUser(channelType, userId)
}

func (srv *passengerService) Update(passenger *entity.Passenger) error {
	return srv.repo.Update(passenger)
}

func (srv *passengerService) Delete(id string) error {
	return srv.repo.Delete(id)
}
