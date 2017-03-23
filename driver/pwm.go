package driver

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/xanecs/lighthouse/config"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type pwmDriver struct {
	drivers    []*gpio.DirectPinDriver
	inverted   bool
	power      bool
	brightness byte
}

func (p *pwmDriver) write() error {
	for _, driver := range p.drivers {
		var err error
		if !p.power {
			if p.inverted {
				err = driver.PwmWrite(255)
			} else {
				err = driver.Off()
			}
		} else {
			brightness := p.brightness
			if p.inverted {
				brightness = 255 - brightness
			}
			err = driver.PwmWrite(brightness)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *pwmDriver) HandleMessage(action string, par params) error {
	switch action {
	case cmdOn:
		p.power = true

	case cmdOff:
		p.power = false

	case "brightness":
		brightness, err := par.getFloat("brightness")
		if err != nil {
			return err
		}
		p.brightness = uint8(255 * brightness)

	case cmdPower:
		power, err := par.getBool("power")
		if err != nil {
			return err
		}
		p.power = power

	case cmdWrite:
		power, err := par.getBool("power")
		if err != nil {
			return err
		}
		brightness, err := par.getFloat("brightness")
		if err != nil {
			return err
		}
		p.power = power
		p.brightness = uint8(brightness * 255)

	default:
		return errors.New(errInvalidCmd + " " + action)
	}
	return p.write()
}

func (p *pwmDriver) Status() map[string]string {
	return map[string]string{"power": fmt.Sprint(p.power), "brightness": fmt.Sprint(float64(p.brightness) / 255.0)}
}

func (p *pwmDriver) Restore(status map[string]string) error {
	pwr, ok := status["power"]
	if !ok {
		return errors.New("Missing parameter 'power'")
	}
	p.power = pwr == trueStr
	brightStr, ok := status["brightness"]
	if !ok {
		return errors.New("Missing parameter 'brightness'")
	}
	bright, err := strconv.ParseFloat(brightStr, 64)
	if err != nil {
		return err
	}
	p.brightness = byte(bright * 255)
	return p.write()
}

func newPwmDriver(cfg config.DeviceConfig, connection gobot.Connection) (*pwmDriver, error) {
	if len(cfg.Pins) < 1 {
		return nil, errors.New("Invalid number of pins")
	}
	drivers := make([]*gpio.DirectPinDriver, len(cfg.Pins))
	for i, pin := range cfg.Pins {
		drivers[i] = gpio.NewDirectPinDriver(connection, pin)
	}
	return &pwmDriver{
		drivers:    drivers,
		inverted:   cfg.Inverted,
		power:      false,
		brightness: 255,
	}, nil
}
