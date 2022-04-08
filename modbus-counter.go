package main

import (
	"time"
	"uzmar.modbus.timecounter/service"
)

func main() {
	service.BuildModbusServer()
	service.GetCounterServiceInstance([]string{"0", "1"}...)
	a := service.BuildModbusClient(&service.ModbusClientBuilder{Address: "192.168.1.34:502", SlaveId: 0x1, Timeout: 3 * time.Second, DeviceId: "0"})
	b := service.BuildModbusClient(&service.ModbusClientBuilder{Address: "192.168.1.20:502", SlaveId: 0x1, Timeout: 3 * time.Second, DeviceId: "1"})
	//modbusClient := modbus.NewClient(modbusClientHandler)
	//modbusClient.WriteMultipleRegisters(10, 2, converterParser.ConvertInt32ToBytes(315360000))
	//modbusClientHandler.Logger = log.New(os.Stdout, "debug:", log.LstdFlags)
	for {
		go a.AsyncReader()
		go b.AsyncReader()
		time.Sleep(1 * time.Second)
	}
}
