package driver

import (
	"errors"

	"github.com/xanecs/lighthouse/config"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type pwmDriver struct {
	driver     *gpio.DirectPinDriver
	inverted   bool
	power      bool
	brightness byte
}

func (p *pwmDriver) write() error {
	if !p.power {
		return p.driver.Off()
	}
	brightness := p.brightness
	if p.inverted {
		brightness = 255 - brightness
	}
	return p.driver.PwmWrite(brightness)
}

func (p *pwmDriver) HandleMessage(action string, params map[string]interface{}) error {
	switch action {
	case "on":
		p.power = true

	case "off":
		p.power = false

	case "brightness":
		val := params["brightness"]
		if val == nil {
			return errors.New("Missing parameter 'brightness'")
		}
		brightness, ok := val.(float64)
		if !ok {
			return errors.New("Invalid parameter 'brightness'")
		}
		p.brightness = uint8(255 * brightness)

	case "write":
		val := params["power"]
		if val == nil {
			return errors.New("Missing parameter 'power'")
		}
		power, ok := val.(bool)
		if !ok {
			return errors.New("Invalid parameter 'power'")
		}
		val = params["brightness"]
		if val == nil {
			return errors.New("Missing parameter 'brightness'")
		}
		brightness, ok := val.(float64)
		if !ok {
			return errors.New("Invalid parameter 'brightness'")
		}
		p.power = power
		p.brightness = uint8(brightness * 255)
	}
	return p.write()
}

func (p *pwmDriver) Status() map[string]interface{} {
	return map[string]interface{}{"power": p.power, "brightness": p.brightness / 255.0}
}

func (p *pwmDriver) Driver() gobot.Device {
	return p.driver
}

func newPwmDriver(cfg config.DeviceConfig, connection gobot.Connection) (*pwmDriver, error) {
	if len(cfg.Pins) != 1 {
		return nil, errors.New("Invalid number of pins")
	}
	return &pwmDriver{
		driver:     gpio.NewDirectPinDriver(connection, cfg.Pins[0]),
		inverted:   cfg.Inverted,
		power:      false,
		brightness: 255,
	}, nil
}
