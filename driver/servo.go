package driver

import (
	"errors"
	"fmt"
	"strconv"

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

func (s *servoDriver) Status() map[string]string {
	return map[string]string{"power": fmt.Sprint(s.power), "angle": fmt.Sprint(s.angle)}
}

func (s *servoDriver) Restore(status map[string]string) error {
	powerStr, ok := status["power"]
	if !ok {
		return errors.New("Missing parameter 'power'")
	}
	s.power = powerStr == trueStr

	angleStr, ok := status["angle"]
	if !ok {
		return errors.New("Missing parameter 'angle'")
	}
	angle, err := strconv.ParseUint(angleStr, 10, 8)
	if err != nil {
		return err
	}
	s.angle = uint8(angle)
	return s.write()
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
