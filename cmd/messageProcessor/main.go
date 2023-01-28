package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"twitch_chat_analysis/internal/data"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type config struct {
	queueURL string
	cacheURL string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	var cfg config
	flag.StringVar(&cfg.queueURL, "queueURL", "amqp://user:password@localhost:7001/", "RabbitMQ connection url")
	flag.StringVar(&cfg.cacheURL, "cacheURL", "localhost:6379", "Redis connection url")

	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.cacheURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.TODO()
	_, err := rdb.Ping(ctx).Result()
	failOnError(err, "Failed to connect to Redis:")

	conn, err := amqp.Dial(cfg.queueURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"messageQueue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var newMsg data.Message
			log.Printf("Received a message: %s", d.Body)
			json.Unmarshal([]byte(d.Body), &newMsg)
			var existingData []data.Message
			key := fmt.Sprintf("%s-%s", newMsg.Sender, newMsg.Receiver)
			reverseKey := fmt.Sprintf("%s-%s", newMsg.Receiver, newMsg.Sender)
			numKey, _ := rdb.Exists(ctx, key).Result()
			numRvsKey, _ := rdb.Exists(ctx, reverseKey).Result()
			num := numKey + numRvsKey
			if num > 0 {
				var keyToUse string
				if numKey == 1 {
					keyToUse = key
				}
				if numRvsKey == 1 {
					keyToUse = reverseKey
				}
				val, _ := rdb.Get(ctx, keyToUse).Result()
				json.Unmarshal([]byte(val), &existingData)
				existingData = append(existingData, newMsg)
				arrMarshalled, _ := json.Marshal(existingData)
				rdb.Set(ctx, keyToUse, string(arrMarshalled), 0)
			} else {
				arr := []data.Message{newMsg}
				arrMarshalled, _ := json.Marshal(arr)
				rdb.Set(ctx, key, string(arrMarshalled), 0)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
