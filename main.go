package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/troian/hid"
)

// This is the only headset I care about atm
const (
	VENDOR_LOGITECH     uint16 = 0x046d
	ID_LOGITECH_PRO_X_1 uint16 = 0x0aba
)

func main() {
	ret, err := run()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(ret)
}

func run() (string, error) {
	handle, err := getHeadset()
	if err != nil {
		return "", fmt.Errorf("can't get handle: %w", err)
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
		return "", fmt.Errorf("failed hid write: %w", err)
	}

	res := make([]byte, 7)
	_, err = handle.Read(res)
	if err != nil {
		return "", fmt.Errorf("failed hid read: %w", err)
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

	return fmt.Sprintf("status: %s\nlevel: %d", status, level), nil
}

func getHeadset() (*hid.Device, error) {
	devices := hid.Enumerate(VENDOR_LOGITECH, ID_LOGITECH_PRO_X_1)
	return devices[0].Open()
}
