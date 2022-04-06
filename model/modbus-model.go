package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"sync"
)

type DryInterface struct {
	gorm.Model
	Name            string
	State           bool
	RunningTime     uint64
	DryInterfaceLog []DryInterface `gorm:"foreignKey:ID"`
}

type DryInterfaceLog struct {
	gorm.Model
	State          bool
	DryInterfaceID uint
}

var lock = &sync.Mutex{}
var modbusModelInstance *gorm.DB

func GetDB() *gorm.DB {
	if modbusModelInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		log.Printf("Creating new db instance")
		modbusModelInstance = connectAndAutoMigrate()
	}
	return modbusModelInstance
}

func connectAndAutoMigrate() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect database")
	}
	db.AutoMigrate(&DryInterfaceLog{}, &DryInterface{})
	return db
}
