package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"

	"github.com/troian/hid"
)

// This is the only headset I care about atm
const VENDOR_LOGITECH uint16 = 0x046d
const ID_LOGITECH_PRO_X_1 uint16 = 0x0aba

func getHeadset() (*hid.Device, error) {
	devices := hid.Enumerate(VENDOR_LOGITECH, ID_LOGITECH_PRO_X_1)
	return devices[0].Open()
}

func main() {
	handle, err := getHeadset()
	if err != nil {
		log.Fatalf("can't get handle: %v", err)
	}
	defer handle.Close()

	reqBatMsg := []byte{
		0x11, 0xff, 0x06, 0x0b, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
	}
	_, err = handle.Write(reqBatMsg)
	if err != nil {
		handle.Close()
		log.Fatalf("failed hid write: %v", err)
	}

	res := make([]byte, 7)
	_, err = handle.Read(res)
	if err != nil {
		handle.Close()
		log.Fatalf("failed hid read: %v", err)
	}

	voltage := float64(binary.BigEndian.Uint16(res[4:6]))

	var level uint8
	if voltage <= 3525 {
		level = uint8((0.03 * voltage) - 101)
	} else if voltage > 4030 {
		level = 100
	} else {
		level = uint8(0.0000000037268473047*math.Pow(voltage, 4) - 0.00005605626214573775*math.Pow(voltage, 3) + 0.3156051902814949*math.Pow(voltage, 2) - 788.0937250298629*voltage + 736315.3077118985)
	}

	statusCode := uint8(res[6])
	var status string
	switch statusCode {
	case 1:
		status = "discharging"
	case 3:
		status = "charging"
	default:
		status = "unknown"
	}

	fmt.Printf("status: %s\nlevel: %d\n", status, level)
}
