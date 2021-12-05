package config

import (
	"gorestkafka/internal/errors"
	"strconv"

	"github.com/fasthttp-contrib/websocket"
)

type ConsumeConfig struct {
	Timeout       int
	Broker        string
	Topic         string
	CommitEnable  bool
	ConsumerGroup string
}

func NewConsumeConfig(c *websocket.Conn) (*ConsumeConfig, error) {
	consumeConfig := ConsumeConfig{}

	timeout := c.Header("timeout")
	t := 10
	if len(timeout) > 0 {
		t, _ = strconv.Atoi(timeout)
	}

	broker := c.Header("broker")
	if len(broker) < 1 {
		return nil, &errors.RequestError{StatusCode: 400, Message: "broker not specified"}
	}
	topic := c.Header("topic")
	if len(topic) < 1 {
		return nil, &errors.RequestError{StatusCode: 400, Message: "topic not specified"}
	}

	consumerGroup := c.Header("consumer-group")
	commitEnable := c.Header("commit-enable")
	nc := len(consumerGroup) > 0
	if len(commitEnable) > 0 {
		nc = true
	}
	consumeConfig.Timeout = t
	consumeConfig.Broker = broker
	consumeConfig.Topic = topic
	consumeConfig.CommitEnable = nc
	consumeConfig.ConsumerGroup = consumerGroup

	return &consumeConfig, nil
}
