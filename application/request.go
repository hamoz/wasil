package application

import (
	"github.com/hamoz/wasil/domain/entity"
	"github.com/hamoz/wasil/domain/service"
)

type RequestApp struct {
	requestService service.RequestService
	driverService  service.DriverService
}

type RequestApplication interface {
	PlaceRequest(*entity.Request) (*entity.Request, error)
	GetRequest(id string) (*entity.Request, error)
	GetRequests() ([]entity.Request, error)
	DeleteRequest(id string) error
	Start()
}

func (rApp *RequestApp) NewRequestApp() *RequestApp {
	return nil
}
