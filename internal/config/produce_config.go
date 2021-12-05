package config

import (
	"gorestkafka/internal/errors"
	"strconv"

	"github.com/valyala/fasthttp"
)

type ProduceConfig struct {
	Broker    string
	Topic     string
	Partition int
}

func NewProduceConfig(ctx *fasthttp.RequestCtx) (*ProduceConfig, *errors.RequestError) {
	broker := string(ctx.Request.Header.Peek("broker"))
	topic := ctx.UserValue("topic").(string)
	partition := string(ctx.Request.Header.Peek("partition"))
	p := 0
	if len(partition) > 0 {
		p, _ = strconv.Atoi(partition)
	}
	if len(broker) < 1 {
		return nil, &errors.RequestError{StatusCode: 400, Message: "broker not specified"}
	}

	return &ProduceConfig{Broker: broker, Topic: topic, Partition: p}, nil
}
