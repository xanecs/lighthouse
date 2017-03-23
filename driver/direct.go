package driver

import (
	"errors"
	"fmt"

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

func (d *directDriver) HandleMessage(action string, p params) error {
	switch action {
	case cmdOn:
		d.state = true

	case cmdOff:
		d.state = false

	case cmdWrite:
		power, err := p.getBool("power")
		if err != nil {
			return err
		}
		d.state = power

	default:
		return errors.New(errInvalidCmd + " " + action)
	}
	return d.write()
}

func (d *directDriver) Status() map[string]string {
	return map[string]string{"power": fmt.Sprint(d.state)}
}

func (d *directDriver) Restore(status map[string]string) error {
	v, ok := status["power"]
	if !ok {
		return errors.New("Missing parameter power")
	}
	d.state = v == trueStr
	return d.write()
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
