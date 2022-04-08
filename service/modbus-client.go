package service

import (
	"fmt"
	"github.com/goburrow/modbus"
	"strconv"
	"strings"
	"time"
)

type ModbusClientBuilder struct {
	Address  string
	SlaveId  byte
	Timeout  time.Duration
	DeviceId string
}
type modbusClient struct {
	AsyncReader func()
}

var modbusMasterClient modbus.Client
var modbusDeviceClient modbus.Client
var deviceId string

func asyncReader() {
	result, _ := modbusDeviceClient.ReadDiscreteInputs(0, 16)
	unscramble(result)
}

func unscramble(result []byte) {
	if len(result) > 1 {
		binary := strings.Trim(fmt.Sprintf("%08b", result[1]), " ")
		binary += strings.Trim(fmt.Sprintf("%08b", result[0]), " ")
		diIncrement := 0
		for i := len(binary) - 1; i >= 0; i-- {
			dryName := "DI_" + deviceId + "_" + strconv.FormatInt(int64(diIncrement), 10)
			status, _ := strconv.ParseBool(binary[i : i+1])
			GetCounterServiceInstance().CalculationTimer(dryName, status)
			diIncrement++
		}
	}
}
func createModbusClient(builder *ModbusClientBuilder) *modbus.TCPClientHandler {
	handler := modbus.NewTCPClientHandler(builder.Address)
	handler.SlaveId = builder.SlaveId
	handler.Timeout = builder.Timeout
	_ = handler.Connect()
	defer handler.Close()
	return handler
}

func BuildModbusClient(builder *ModbusClientBuilder) modbusClient {
	modbusDeviceClient = modbus.NewClient(createModbusClient(builder))
	modbusMasterClient = modbus.NewClient(createModbusClient(&ModbusClientBuilder{Address: "0.0.0.0:502", SlaveId: 0x1, Timeout: 3 * time.Second, DeviceId: "0"}))
	deviceId = builder.DeviceId
	return modbusClient{AsyncReader: asyncReader}
}
