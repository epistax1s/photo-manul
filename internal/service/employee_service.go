package service

import (
	"github.com/epistax1s/photo-manul/internal/model"
	"github.com/epistax1s/photo-manul/internal/repository"
)

type EmployeeService interface {
	CreateEmployee(employee *model.Employee) error
	UpdateEmployee(employee *model.Employee) error
	GetEmployeeByID(employeeID int64) (*model.Employee, error)
}

type employeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

func (service *employeeService) CreateEmployee(employee *model.Employee) error {
	return service.repo.Create(employee)
}

func (service *employeeService) UpdateEmployee(employee *model.Employee) error {
	return service.repo.Update(employee)
}

func (service *employeeService) GetEmployeeByID(employeeID int64) (*model.Employee, error) {
	return service.repo.FindByEmployeeID(employeeID)
}
