package driver

import (
	"errors"

	"github.com/xanecs/lighthouse/config"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

const (
	cmdOn    = "on"
	cmdOff   = "off"
	cmdWrite = "write"
	cmdPower = "power"
)

const (
	errInvalidCmd = "Invalid command"
)

// Device represents a device that can handle messages
type Device interface {
	HandleMessage(string, params) error
	Status() map[string]string
	Restore(map[string]string) error
}

// NewDriver creates a new Driver from a device config
func NewDriver(cfg config.DeviceConfig, connection gobot.Connection) (Device, error) {
	switch cfg.Mode {
	case "direct":
		return newDirectDriver(cfg, connection)

	case "pwm":
		return newPwmDriver(cfg, connection)

	case "rgb":
		dd, ok := connection.(gpio.DigitalWriter)
		if !ok {
			return nil, errors.New("rgb mode is not supported on this board")
		}
		return newRgbDriver(cfg, dd)

	case "servo":
		return newServoDriver(cfg, connection)

	default:
		return nil, errors.New("Invalid mode")
	}
}
