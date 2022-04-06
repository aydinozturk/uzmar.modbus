package main

import (
	"fmt"
	"github.com/goburrow/modbus"
	"github.com/tbrandon/mbserver"
	"log"
	"strconv"
	"strings"
	"time"
	"uzmar.modbus.timecounter/service"
	converterParser "uzmar.modbus.timecounter/src"
)

func asyncReader(client modbus.Client) {
	result, _ := client.ReadDiscreteInputs(0, 16)
	unscramble(result)
}

func unscramble(result []byte) {
	if len(result) > 1 {
		binary := strings.Trim(fmt.Sprintf("%08b", result[1]), " ")
		binary += strings.Trim(fmt.Sprintf("%08b", result[0]), " ")
		diIncrement := 0
		for i := len(binary) - 1; i >= 0; i-- {
			dryName := "DI" + strconv.FormatInt(int64(diIncrement), 10)
			status, _ := strconv.ParseBool(binary[i : i+1])
			service.GetInstance().CalculationTimer(dryName, status)
			diIncrement++
		}
	}
}

func main() {
	modbusServer := mbserver.NewServer()
	modbusErr := modbusServer.ListenTCP("0.0.0.0:8502")
	if modbusErr != nil {
		log.Printf("%v\n", modbusErr)
	}
	defer modbusServer.Close()
	service.GetInstance()
	modbusClientHandler := modbus.NewTCPClientHandler("0.0.0.0:8502")
	modbusClientHandler.SlaveId = 0x1
	modbusClientHandler.Timeout = 3 * time.Second
	_ = modbusClientHandler.Connect()
	defer modbusClientHandler.Close()

	modbusClient := modbus.NewClient(modbusClientHandler)
	modbusClient.WriteMultipleRegisters(10, 2, converterParser.ConvertInt32ToBytes(315360000))
	//modbusClientHandler.Logger = log.New(os.Stdout, "debug:", log.LstdFlags)

	handler := modbus.NewTCPClientHandler("192.168.1.20:502")
	handler.SlaveId = 0x1
	handler.Timeout = 3 * time.Second
	_ = handler.Connect()
	defer handler.Close()

	client := modbus.NewClient(handler)
	for {
		go asyncReader(client)
		time.Sleep(1 * time.Second)
	}
}
