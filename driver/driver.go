package driver

import (
	"errors"

	"github.com/xanecs/lighthouse/config"
	"gobot.io/x/gobot"
)

// Device represents a device that can handle messages
type Device interface {
	HandleMessage(string, map[string]interface{}) error
	Status() map[string]interface{}
	Driver() gobot.Device
}

// NewDriver creates a new Driver from a device config
func NewDriver(cfg config.DeviceConfig, connection gobot.Connection) (Device, error) {
	switch cfg.Mode {
	case "direct":
		return newDirectDriver(cfg, connection)

	case "pwm":
		return newPwmDriver(cfg, connection)

	default:
		return nil, errors.New("Invalid mode")
	}
}
