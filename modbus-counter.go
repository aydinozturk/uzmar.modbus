package main

import (
	"fmt"
	"github.com/goburrow/modbus"
	"github.com/tbrandon/mbserver"
	"log"
	"math/rand"
	"strconv"
	"time"
	modbusModel "uzmar.modbus.timecounter/model"
	converterParser "uzmar.modbus.timecounter/package"
)

func asyncListener(client modbus.Client) {
	result, _ := client.ReadDiscreteInputs(0, 16)
	decompiler(result)
	modbusModel.GetInstance()
}

func decompiler(result []byte) {
	if len(result) > 1 {
		binaryNumber := strconv.FormatUint(uint64(result[0]), 2)
		binaryNumber2 := strconv.FormatUint(uint64(result[1]), 2)
		sectionOneLen := len(binaryNumber)
		sectionTwoLen := len(binaryNumber2)
		for i := 0; i < sectionTwoLen; i++ {
			fmt.Println("DI", sectionTwoLen+7-i, "=", binaryNumber2[i:i+1])
		}
		for i := 0; i < sectionOneLen; i++ {
			fmt.Println("DI", sectionOneLen-1-i, "=", binaryNumber[i:i+1])
		}
	}
}

func main() {
	modbusModel.GetInstance()
	modbusServer := mbserver.NewServer()
	modbusErr := modbusServer.ListenTCP("0.0.0.0:8502")
	if modbusErr != nil {
		log.Printf("%v\n", modbusErr)
	}
	defer modbusServer.Close()

	modbusClientHandler := modbus.NewTCPClientHandler("0.0.0.0:8502")
	modbusClientHandler.SlaveId = 0x1
	modbusClientHandler.Timeout = 3 * time.Second
	_ = modbusClientHandler.Connect()
	defer modbusClientHandler.Close()

	modbusClient := modbus.NewClient(modbusClientHandler)
	modbusClient.WriteMultipleRegisters(10, 4, converterParser.ConvertInt64ToBytes(rand.Int63()))
	//modbusClientHandler.Logger = log.New(os.Stdout, "debug:", log.LstdFlags)

	handler := modbus.NewTCPClientHandler("192.168.127.253:502")
	handler.SlaveId = 0x1
	handler.Timeout = 3 * time.Second
	_ = handler.Connect()
	defer handler.Close()

	client := modbus.NewClient(handler)
	for {
		go asyncListener(client)
		time.Sleep(1 * time.Second)
	}
}
