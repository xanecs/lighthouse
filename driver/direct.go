package driver

import (
	"errors"

	"github.com/xanecs/lighthouse/config"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type directDriver struct {
	drivers  []*gpio.DirectPinDriver
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
	for _, driver := range d.drivers {
		err := driver.DigitalWrite(boolToByte(d.state != d.inverted))
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *directDriver) HandleMessage(action string, params map[string]interface{}) error {
	switch action {
	case cmdOn:
		d.state = true

	case cmdOff:
		d.state = false

	case cmdWrite:
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
		return errors.New(errInvalidCmd + " " + action)
	}
	return d.write()
}

func (d *directDriver) Status() map[string]interface{} {
	return map[string]interface{}{"power": d.state}
}

func newDirectDriver(cfg config.DeviceConfig, connection gobot.Connection) (*directDriver, error) {
	if len(cfg.Pins) < 1 {
		return nil, errors.New("Invalid number of pins")
	}
	drivers := make([]*gpio.DirectPinDriver, len(cfg.Pins))
	for i, pin := range cfg.Pins {
		drivers[i] = gpio.NewDirectPinDriver(connection, pin)
	}
	return &directDriver{
		drivers:  drivers,
		inverted: cfg.Inverted,
		state:    false,
	}, nil
}
