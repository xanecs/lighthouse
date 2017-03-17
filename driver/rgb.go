package driver

import (
	"errors"

	"github.com/xanecs/lighthouse/config"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type rgbDriver struct {
	driver   *gpio.RgbLedDriver
	inverted bool
	power    bool
	color    color
}

type color struct {
	red   byte
	green byte
	blue  byte
}

func (r *rgbDriver) write() error {
	if !r.power {
		if r.inverted {
			return r.driver.SetRGB(255, 255, 255)
		}
		return r.driver.Off()
	}
	if r.inverted {
		return r.driver.SetRGB(255-r.color.red, 255-r.color.green, 255-r.color.blue)
	}
	return r.driver.SetRGB(r.color.red, r.color.green, r.color.blue)
}

func (r *rgbDriver) HandleMessage(action string, params map[string]interface{}) error {
	switch action {
	case cmdOn:
		r.power = true

	case cmdOff:
		r.power = false

	case "color":
		newColor, err := parseColor(params)
		if err != nil {
			return err
		}
		r.color = newColor

	case cmdWrite:
		val := params["power"]
		if val == nil {
			return errors.New("Missing parameter 'power'")
		}
		power, ok := val.(bool)
		if !ok {
			return errors.New("Invalid parameter 'power'")
		}

		newColor, err := parseColor(params)
		if err != nil {
			return err
		}

		r.power = power
		r.color = newColor
	}
	return r.write()
}

func (r *rgbDriver) Status() map[string]interface{} {
	return map[string]interface{}{"power": r.power, "color": r.color}
}

func (r *rgbDriver) Driver() gobot.Device {
	return r.driver
}

func parseColor(params map[string]interface{}) (color, error) {
	valR := params["red"]
	valG := params["green"]
	valB := params["blue"]

	if valR == nil || valG == nil || valB == nil {
		return color{}, errors.New("Missing color parameters")
	}

	red, okR := valR.(float64)
	green, okG := valG.(float64)
	blue, okB := valB.(float64)

	if !(okR && okG && okB) {
		return color{}, errors.New("Invalid color parameters")
	}

	return color{uint8(red), uint8(green), uint8(blue)}, nil
}

func newRgbDriver(cfg config.DeviceConfig, connection gpio.DigitalWriter) (*rgbDriver, error) {
	if len(cfg.Pins) != 3 {
		return nil, errors.New("Invalid number of pins")
	}
	return &rgbDriver{
		driver:   gpio.NewRgbLedDriver(connection, cfg.Pins[0], cfg.Pins[1], cfg.Pins[2]),
		inverted: cfg.Inverted,
		power:    false,
		color:    color{255, 255, 255},
	}, nil
}
