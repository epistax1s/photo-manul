package server

import (
	manul "github.com/epistax1s/photo-manul/internal/bot"
	"github.com/epistax1s/photo-manul/internal/config"
	"github.com/epistax1s/photo-manul/internal/database"
	"github.com/epistax1s/photo-manul/internal/log"
	"github.com/epistax1s/photo-manul/internal/repository"
	"github.com/epistax1s/photo-manul/internal/service"
)

type Server struct {
	UserService     service.UserService
	EmployeeService service.EmployeeService
	Config          *config.Config
	Manul           *manul.Manul
}

func InitServer() *Server {
	config, err := config.LoadConfig()
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}

	log.InitLogger(&config.Log)

	log.Info("Server initialization: telegram bot initialization")

	manul, err := manul.InitTelegramBot(&config.Bot)
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}

	log.Info("Server initialization: database initialization")

	db, err := database.InitDatabase(config)
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}

	log.Info("Server initialization: run migration")

	migrationManager := database.NewMigrationManager(db)
	migrationManager.RunMigration()

	log.Info("Server initialization: services initialization")

	userService := service.NewUserService(
		repository.NewUserRepository(db),
	)

	employeeService := service.NewEmployeeService(
		repository.NewEmployeeRepository(db),
	)

	return &Server{
		UserService:     userService,
		EmployeeService: employeeService,
		Config:          config,
		Manul:           manul,
	}
}
