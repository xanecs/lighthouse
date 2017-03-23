package main

import (
	"github.com/xanecs/lighthouse/config"
	redis "gopkg.in/redis.v5"
)

type broker struct {
	client *redis.Client
	pubsub *redis.PubSub
	topic  string
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
		topic:  cfg.Topic,
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

func (b *broker) updateStatus(status chan Status, errors chan error) {
	for {
		s := <-status
		err := b.client.HMSet(b.topic+":"+s.Device, s.Params).Err()
		if err != nil {
			errors <- err
		}

	}
}

func (b *broker) fetchStatus(device string) (map[string]string, error) {
	return b.client.HGetAll(b.topic + ":" + device).Result()
}

func (b *broker) Close() {
	b.pubsub.Close()
	b.client.Close()
}
