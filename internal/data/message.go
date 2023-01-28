package data

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type MessageModel struct {
	Cache *redis.Client
	Queue *amqp.Connection
}

type Message struct {
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type MessageModelInterface interface {
	CreateMessage(message Message) error
	GetMessages(key, reverseKey string) ([]Message, error)
}

func (m MessageModel) CreateMessage(message Message) error {
	ch, err := m.Queue.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()
	q, err := ch.QueueDeclare(
		"messageQueue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}
	log.Printf(" [x] Sent %s\n", body)
	return nil
}

func (m MessageModel) GetMessages(key, reverseKey string) ([]Message, error) {
	ctx := context.TODO()
	var result []Message
	numKey, _ := m.Cache.Exists(ctx, key).Result()
	numRvsKey, _ := m.Cache.Exists(ctx, reverseKey).Result()
	num := numKey + numRvsKey

	if num > 0 {
		var keyToUse string
		if numKey == 1 {
			keyToUse = key
		}
		if numRvsKey == 1 {
			keyToUse = reverseKey
		}
		val, err := m.Cache.Get(ctx, keyToUse).Result()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(val), &result)
		if err != nil {
			return nil, err
		}
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result, nil
}
