package repository

import (
	"fmt"

	"github.com/epistax1s/photo-manul/internal/model"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(employee *model.Employee) error
	Update(employee *model.Employee) error
	FindByEmployeeID(employeeID int64) (*model.Employee, error)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (repo *employeeRepository) Create(employee *model.Employee) error {
	return repo.db.Create(employee).Error
}

func (repo *employeeRepository) Update(employee *model.Employee) error {
	return repo.db.Save(employee).Error
}

func (repo *employeeRepository) FindByEmployeeID(employeeID int64) (*model.Employee, error) {
	var employee model.Employee

	result := repo.db.
		Where(fmt.Sprintf("%s = ?", model.EmployeeIDColumn), employeeID).
		First(&employee)

	if result.Error != nil {
		return nil, result.Error
	}

	return &employee, nil
}
