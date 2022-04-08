package service

import (
	"github.com/tbrandon/mbserver"
	"log"
)

func BuildModbusServer() {
	modbusServer := mbserver.NewServer()
	modbusErr := modbusServer.ListenTCP("0.0.0.0:502")
	if modbusErr != nil {
		log.Printf("%v\n", modbusErr)
	}
	defer modbusServer.Close()
}
