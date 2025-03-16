package model

type Employee struct {
	ID           int64  `gorm:"column:id;type:bigserial;primaryKey"`
	EmployeeID   int64  `gorm:"column:employee_id;type:bigint;unique;not null"`
	EmployeeName string `gorm:"column:employee_name;type:varchar(255);not null"`
	ImagePath    string `gorm:"column:image_path;type:varchar(255)"`
}

const (
	EmployeeTable           = "employees"
	EmployeeIDColumn        = "employee_id"
	EmployeeNameColumn      = "employee_name"
	EmployeeImagePathColumn = "image_path"
)

func (Employee) TableName() string {
	return EmployeeTable
}
