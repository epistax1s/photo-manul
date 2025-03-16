package database

import (
	"github.com/epistax1s/photo-manul/internal/log"
	"github.com/epistax1s/photo-manul/internal/model"
	"gorm.io/gorm"
)

type MigrationManager struct {
	db *gorm.DB
}

func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db: db,
	}
}

func (manager *MigrationManager) RunMigration() {
	hasTable := manager.db.Migrator().HasTable(&model.User{})
	if !hasTable {
		err := manager.db.AutoMigrate(&model.User{}, &model.Employee{})
		if err != nil {
			log.Error("Migration failure", "err", err)
			return
		}
		log.Info("Migration success!")
	} else {
		log.Info("Migration was performed earlier")
	}
}
