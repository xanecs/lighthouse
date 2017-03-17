package main

import "encoding/json"

// Message represents a message received from redis
type Message struct {
	Device string                 `json:"device"`
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params"`
}

// Parser parses messages from and to a chan
func Parser(in chan string, out chan Message, errOut chan error) {
	for {
		rawMsg := <-in
		var msg Message
		err := json.Unmarshal([]byte(rawMsg), &msg)
		if err != nil {
			errOut <- err
			continue
		}
		out <- msg
	}
}
