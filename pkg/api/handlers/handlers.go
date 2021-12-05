package handlers

import (
	"context"
	"encoding/json"
	"gorestkafka/internal/config"
	"log"
	"time"

	"github.com/fasthttp-contrib/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/valyala/fasthttp"
)

type BodyMessage interface{}

type ProduceBody struct {
	Messages []BodyMessage `json:"messages"`
}

var upgrader = websocket.New(consume)

func ProduceHandler(ctx *fasthttp.RequestCtx) {
	ctxBackground := context.Background()

	produceConfig, re := config.NewProduceConfig(ctx)
	if re != nil {
		ctx.Response.SetStatusCode(re.StatusCode)
		ctx.Response.SetBody([]byte(re.Message))
		return
	}

	log.Printf("get produce request on %s broker, %s topic", produceConfig.Broker, produceConfig.Topic)

	produceBody := ProduceBody{}
	json.Unmarshal(ctx.PostBody(), &produceBody)
	msgs := []kafka.Message{}
	for _, value := range produceBody.Messages {
		messageString, _ := json.Marshal(value)
		msgs = append(msgs, kafka.Message{Value: messageString})
	}

	conn, err := kafka.DialLeader(ctxBackground, "tcp", produceConfig.Broker, produceConfig.Topic, 0)

	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		msgs...,
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func ConsumeHandler(ctx *fasthttp.RequestCtx) {
	upgrader.Upgrade(ctx)
}

func consume(c *websocket.Conn) {
	consumerConfig, err := config.NewConsumeConfig(c)
	if err != nil {
		c.WriteMessage(1, []byte(err.Error()))
		c.Close()
		return
	}

	c.SetWriteDeadline(time.Now().Add(time.Duration(consumerConfig.Timeout) * time.Second))
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{consumerConfig.Broker},
		GroupID:         consumerConfig.ConsumerGroup,
		Topic:           consumerConfig.Topic,
		MinBytes:        10e3,
		MaxBytes:        10e6,
		SessionTimeout:  time.Duration(consumerConfig.Timeout) * time.Second,
		ReadLagInterval: -1,
	})
	defer c.Close()
	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Fatal("Failed to fetch message", err)
			break
		}
		log.Printf("Fetch message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		err = c.WriteMessage(1, m.Value)
		if err != nil {
			log.Fatal("Failed to write message : ", err)
			break
		}
		if consumerConfig.CommitEnable {
			log.Printf("Commit message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			r.CommitMessages(context.Background(), m)
		}
	}
	return
}
