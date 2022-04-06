package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	modbusModel "uzmar.modbus.timecounter/model"
)

type modbusService struct {
	DryInputs        map[string]bool
	CalculationTimer func(name string, status bool)
}

var lock = &sync.Mutex{}
var modbusServiceInstance *modbusService

func GetInstance() *modbusService {
	if modbusServiceInstance == nil {
		lock.Lock()
		log.Println("Creating new modbusService Instance")
		modbusServiceInstance = &modbusService{DryInputs: createAndFillNewDryInputs(), CalculationTimer: calculationTimer}
	}
	return modbusServiceInstance
}

func createAndFillNewDryInputs() map[string]bool {
	DB := modbusModel.GetDB()
	dryInputs := make(map[string]bool)
	for i := 0; i < 16; i++ {
		var dryItem modbusModel.DryInterface
		dryItemResult := DB.Where(&modbusModel.DryInterface{Name: "DI" + strconv.FormatInt(int64(i), 10)}).First(&dryItem)
		if errors.Is(dryItemResult.Error, gorm.ErrRecordNotFound) {
			diName := "DI" + strconv.FormatInt(int64(i), 10)
			DB.Create(&modbusModel.DryInterface{Name: diName, State: false, RunningTime: 0})
			dryInputs[diName] = false
		} else {
			dryInputs[dryItem.Name] = dryItem.State
		}
	}
	return dryInputs
}

func calculationTimer(name string, status bool) {
	DB := modbusModel.GetDB()
	input := modbusServiceInstance.DryInputs[name]
	if input != status {
		fmt.Println(name, status)
		DB.Where("Name = '?'", name).Updates(&modbusModel.DryInterface{State: status})
		modbusServiceInstance.DryInputs[name] = status
	}
}
