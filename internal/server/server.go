package server

import (
	manul "github.com/epistax1s/photo-manul/internal/bot"
	"github.com/epistax1s/photo-manul/internal/config"
	"github.com/epistax1s/photo-manul/internal/log"
)

type Server struct {
	Config *config.Config
	Manul  *manul.Manul
}

func InitServer() *Server {
	config, err := config.LoadConfig()
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}

	log.InitLogger(&config.Log)
	log.Info("Server initialization: loading the configuration")

	manul, err := manul.InitTelegramBot(&config.Bot)
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}
	log.Info("Server initialization: telegram bot initialization")

	return &Server{
		Config: config,
		Manul:  manul,
	}
}
