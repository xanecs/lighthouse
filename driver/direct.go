package driver

import (
	"errors"

	"github.com/xanecs/lighthouse/config"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type directDriver struct {
	driver   *gpio.DirectPinDriver
	inverted bool
	state    bool
}

func boolToByte(v bool) byte {
	if v {
		return 1
	}
	return 0
}

func (d *directDriver) write() error {
	// != means xor here
	err := d.driver.DigitalWrite(boolToByte(d.state != d.inverted))
	return err
}

func (d *directDriver) HandleMessage(action string, params map[string]interface{}) error {
	switch action {
	case "on":
		d.state = true

	case "off":
		d.state = false

	case "write":
		val := params["power"]
		if val == nil {
			return errors.New("Missing parameter 'power'")
		}

		power, ok := val.(bool)
		if !ok {
			return errors.New("Invalid parameter 'power'")
		}

		d.state = power

	default:
		return errors.New("Invalid command")
	}
	return d.write()
}

func (d *directDriver) Status() map[string]interface{} {
	return map[string]interface{}{"power": d.state}
}

func (d *directDriver) Driver() gobot.Device {
	return d.driver
}

func newDirectDriver(cfg config.DeviceConfig, connection gobot.Connection) (*directDriver, error) {
	if len(cfg.Pins) != 1 {
		return nil, errors.New("Invalid number of pins")
	}
	return &directDriver{
		driver:   gpio.NewDirectPinDriver(connection, cfg.Pins[0]),
		inverted: cfg.Inverted,
		state:    false,
	}, nil
}
