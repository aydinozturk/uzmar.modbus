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
	DryInputs        map[string]modbusModel.DryInput
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

func createAndFillNewDryInputs() map[string]modbusModel.DryInput {
	DB := modbusModel.GetDB()
	dryInputs := make(map[string]modbusModel.DryInput)
	for i := 0; i < 16; i++ {
		var dryItem modbusModel.DryInput
		dryItemResult := DB.Where(&modbusModel.DryInput{Name: "DI" + strconv.FormatInt(int64(i), 10)}).First(&dryItem)
		diName := "DI" + strconv.FormatInt(int64(i), 10)
		if errors.Is(dryItemResult.Error, gorm.ErrRecordNotFound) {
			DB.Create(&modbusModel.DryInput{Name: diName, Status: false, RunningTime: 0}).First(&dryItem)
		}
		dryInputs[diName] = dryItem
	}
	return dryInputs
}

func calculationTimer(name string, status bool) {
	DB := modbusModel.GetDB()
	dryInput := modbusServiceInstance.DryInputs[name]
	if dryInput.Status != status {
		fmt.Println(dryInput.Name, dryInput.Status)
		if status {
			DB.Create(&modbusModel.DryInputLog{DryInputID: dryInput.ID, Status: status})
		} else {
			var dryItemLog modbusModel.DryInputLog
			DB.Model(&modbusModel.DryInputLog{}).Where("dry_input_id = ?", dryInput.ID).Update("Status", status).First(&dryItemLog)

			dryInput.RunningTime += uint64(dryItemLog.UpdatedAt.UnixNano()) - uint64(dryItemLog.CreatedAt.UnixNano())
		}

		DB.Model(&modbusModel.DryInput{}).Where("Name = ?", dryInput.Name).Updates(dryInput).First(&dryInput)
		modbusServiceInstance.DryInputs[dryInput.Name] = dryInput
	}
}
