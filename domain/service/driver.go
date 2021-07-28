package service

import (
	"github.com/hamoz/wasil/domain/entity"
	"github.com/hamoz/wasil/domain/repository"
)

type DriverService interface {
	Register(driver *entity.Driver) (*entity.Driver, error)
	Find(id string) (*entity.Driver, error)
	Update(driver *entity.Driver) error
	Delete(id string) error
}

type driverService struct {
	repo repository.DriverRepository
}

func NewDriverService(driverRepo repository.DriverRepository) *driverService {
	return &driverService{driverRepo}
}

func (srv *driverService) Register(driver *entity.Driver) (*entity.Driver, error) {
	return srv.repo.Create(driver)
}

func (srv *driverService) Find(id string) (*entity.Driver, error) {
	return srv.repo.Find(id)
}

func (srv *driverService) Update(driver *entity.Driver) error {
	return srv.repo.Update(driver)
}

func (srv *driverService) Delete(id string) error {
	return srv.repo.Delete(id)
}
