package service

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"time"
	modbusModel "uzmar.modbus.timecounter/model"
)

type counterService struct {
	DryInputs        map[string]modbusModel.DryInput
	CalculationTimer func(name string, status bool)
}

var lock = &sync.Mutex{}
var modbusServiceInstance *counterService

func GetCounterServiceInstance(devices ...string) *counterService {
	if modbusServiceInstance == nil {
		lock.Lock()
		log.Println("Creating new modbusService Instance")
		modbusServiceInstance = &counterService{DryInputs: createAndFillNewDryInputs(devices...), CalculationTimer: calculationTimer}
	}
	return modbusServiceInstance
}

func createAndFillNewDryInputs(devices ...string) map[string]modbusModel.DryInput {
	DB := modbusModel.GetDB()
	dryInputs := make(map[string]modbusModel.DryInput)
	for _, device := range devices {
		for i := 0; i < 16; i++ {
			var dryItem modbusModel.DryInput
			diName := "DI_" + device + "_" + strconv.FormatInt(int64(i), 10)
			dryItemResult := DB.Where(&modbusModel.DryInput{Name: diName}).First(&dryItem)
			if errors.Is(dryItemResult.Error, gorm.ErrRecordNotFound) {
				DB.Create(&modbusModel.DryInput{Name: diName, Status: false, RunningTime: 0}).First(&dryItem)
			}
			dryInputs[diName] = dryItem
		}
	}
	return dryInputs
}

func calculationTimer(name string, status bool) {
	DB := modbusModel.GetDB()
	dryInput := modbusServiceInstance.DryInputs[name]
	if dryInput.Status != status {
		if !status {
			dryInput.RunningTime += uint64(time.Now().UTC().Sub(dryInput.UpdatedAt.UTC()).Seconds())
		}
		dryInput.Status = status
		DB.Save(&dryInput)
		modbusServiceInstance.DryInputs[name] = dryInput
	}
}
