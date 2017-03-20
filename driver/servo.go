package driver

import (
	"errors"

	"github.com/xanecs/lighthouse/config"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type servoDriver struct {
	drivers []*gpio.DirectPinDriver
	power   bool
	angle   byte
}

func (s *servoDriver) write() error {
	for _, driver := range s.drivers {
		var err error
		if !s.power {
			err = driver.DigitalWrite(0)
		} else {
			err = driver.ServoWrite(s.angle)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *servoDriver) HandleMessage(action string, p params) error {
	switch action {
	case cmdOn:
		s.power = true

	case cmdOff:
		s.power = false

	case cmdPower:
		power, err := p.getBool("power")
		if err != nil {
			return err
		}
		s.power = power

	case "angle":
		angle, err := p.getFloat("angle")
		if err != nil {
			return err
		}
		s.angle = uint8(angle)

	case cmdWrite:
		power, err := p.getBool("power")
		if err != nil {
			return err
		}

		angle, err := p.getFloat("angle")
		if err != nil {
			return err
		}

		s.power = power
		s.angle = uint8(angle)

	default:
		return errors.New(errInvalidCmd + " " + action)

	}
	return s.write()
}

func (s *servoDriver) Status() map[string]interface{} {
	return map[string]interface{}{"power": s.power, "angle": s.angle}
}

func newServoDriver(cfg config.DeviceConfig, connection gobot.Connection) (*servoDriver, error) {
	if len(cfg.Pins) < 1 {
		return nil, errors.New("Invalid number of pins")
	}
	drivers := make([]*gpio.DirectPinDriver, len(cfg.Pins))
	for i, pin := range cfg.Pins {
		drivers[i] = gpio.NewDirectPinDriver(connection, pin)
	}
	return &servoDriver{drivers, false, 0}, nil
}
