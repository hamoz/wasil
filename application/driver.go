package application

import (
	"github.com/hamoz/wasil/domain/entity"
	"github.com/hamoz/wasil/domain/service"
)

type DriverApp struct {
	service service.DriverService
}

type DriverAppInterface interface {
	SaveDriver(*entity.Driver) (*entity.Driver, map[string]string)
	GetDrivers() ([]entity.Driver, error)
	GetDriver(id string) (*entity.Driver, error)
}

func (app *DriverApp) SaveDriver(driver *entity.Driver) (*entity.Driver, error) {
	return app.service.Register(driver)
}

func (app *DriverApp) GetDriver(id string) (*entity.Driver, error) {
	return app.service.Find(id)
}

func (app *DriverApp) DeleteDriver(id string) error {
	return app.service.Delete(id)
}
