package service

import (
	"fmt"

	"github.com/hamoz/wasil/commons/types"
	"github.com/hamoz/wasil/domain/entity"
	"github.com/hamoz/wasil/domain/repository"

	geo "github.com/kellydunn/golang-geo"
)

const KILO_PRICE int = 200
const METER_START int = 600

type RequestService interface {
	Create(request *entity.Request) (*entity.Request, error)
	Find(id string) (*entity.Request, error)
	FindByChannelUser(channel types.ChannelType, channeUserId string) (*entity.Request, error)
	Update(request *entity.Request) (*entity.Request, error)
	Delete(id string) error
	Price(request *entity.Request) (int, error)
}

type requestService struct {
	repo repository.RequestRepository
}

func NewRequestService(repo repository.RequestRepository) *requestService {
	return &requestService{repo}
}

func (srv *requestService) Create(request *entity.Request) (*entity.Request, error) {
	return srv.repo.Create(request)
}
func (srv *requestService) FindByChannelUser(channelType types.ChannelType, channelUserID string) (*entity.Request, error) {
	return srv.repo.FindByChannelUser(channelType, channelUserID)
}

func (srv *requestService) Find(id string) (*entity.Request, error) {
	return srv.repo.Find(id)
}
func (srv *requestService) Update(request *entity.Request) (*entity.Request, error) {
	return srv.repo.Update(request)
}
func (srv *requestService) Delete(id string) error {
	return srv.repo.Delete(id)
}

func (srv *requestService) Price(request *entity.Request) (int, error) {
	if request.FromLoc == nil {
		return -1, fmt.Errorf("from Location is empty")
	}
	if request.ToLocation == nil {
		return -1, fmt.Errorf("to Location is empty")
	}
	from := geo.NewPoint(float64(request.FromLoc.Lat), float64(request.FromLoc.Lng))
	to := geo.NewPoint(float64(request.ToLocation.Lat), float64(request.ToLocation.Lng))
	distance := from.GreatCircleDistance(to)
	fmt.Printf("Request Distance : %f", distance)
	price := float64(METER_START) + distance*float64(KILO_PRICE)
	return int(price), nil

}
