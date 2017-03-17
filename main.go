package main

import (
	"fmt"
	"os"

	"github.com/xanecs/lighthouse/config"
)

func main() {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		panic(err)
	}

	hardware, err := NewHardware(cfg.Boards)
	if err != nil {
		panic(err)
	}
	err = hardware.Start()
	if err != nil {
		panic(err)
	}
	defer hardware.Stop()

	broker, err := newBroker(cfg.Redis)
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	chanErr := make(chan error)
	brokerMsg := make(chan string)
	parsedMsg := make(chan Message)

	go broker.listen(brokerMsg, chanErr)
	go Parser(brokerMsg, parsedMsg, chanErr)
	go hardware.Listen(parsedMsg, chanErr)

	for {
		err := <-chanErr
		fmt.Fprintln(os.Stderr, err)
	}
}
