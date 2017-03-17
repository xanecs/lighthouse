package main

import (
	"github.com/xanecs/lighthouse/config"
	redis "gopkg.in/redis.v5"
)

type broker struct {
	client *redis.Client
	pubsub *redis.PubSub
}

func newBroker(cfg config.RedisConfig) (*broker, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: "",
		DB:       0,
	})
	pubsub, err := client.Subscribe(cfg.Topic)
	if err != nil {
		return nil, err
	}
	return &broker{
		client: client,
		pubsub: pubsub,
	}, nil
}

func (b *broker) listen(messages chan string, errors chan error) {
	for {
		msg, err := b.pubsub.ReceiveMessage()
		if err != nil {
			errors <- err
			continue
		}
		messages <- msg.Payload
	}
}

func (b *broker) Close() {
	b.pubsub.Close()
	b.client.Close()
}
