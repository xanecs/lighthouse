package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/xanecs/lighthouse/config"
	"github.com/xanecs/lighthouse/driver"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/firmata"
)

// Hardware represents all connected controllers
type Hardware struct {
	connections []gobot.Connection
	devices     map[string]driver.Device
}

// NewHardware creates a new Hardware object from a board config
func NewHardware(cfg map[string]config.BoardConfig) (*Hardware, error) {
	connections := make([]gobot.Connection, 0)
	devices := make(map[string]driver.Device)

	for _, boardConfig := range cfg {
		adaptor := firmata.NewAdaptor(boardConfig.Serial)
		connections = append(connections, adaptor)

		for deviceName, deviceConfig := range boardConfig.Dev {
			if devices[deviceName] != nil {
				return nil, errors.New("Name " + deviceName + " has been given twice")
			}
			device, err := driver.NewDriver(deviceConfig, adaptor)
			if err != nil {
				return nil, err
			}
			devices[deviceName] = device
		}
	}

	return &Hardware{connections, devices}, nil
}

// Start connects to the boards
func (h *Hardware) Start() (err error) {
	for _, c := range h.connections {
		err = c.Connect()
		if err != nil {
			return
		}
	}
	return
}

// Stop disconnects from the boards
func (h *Hardware) Stop() {
	for _, c := range h.connections {
		c.Finalize()
	}
}

// Restore all Hardware from database
func (h *Hardware) Restore(b *broker) {
	for name, dev := range h.devices {
		status, err := b.fetchStatus(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		dev.Restore(status)
	}
}

// Listen processes Messages from a channel
func (h *Hardware) Listen(in chan Message, out chan Status, chanErr chan error) {
	for {
		msg := <-in
		dev := h.devices[msg.Device]
		if dev == nil {
			chanErr <- errors.New("Device '" + msg.Device + "' not found")
			continue
		}
		if err := dev.HandleMessage(msg.Action, msg.Params); err != nil {
			chanErr <- err
			continue
		}
		out <- Status{msg.Device, dev.Status()}
	}
}
