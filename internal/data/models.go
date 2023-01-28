package data

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Models struct {
	Message MessageModelInterface
}

type Dependencies struct {
	Queue *amqp.Connection
	Cache *redis.Client
}

func NewModels(deps Dependencies) Models {
	return Models{
		Message: MessageModel{Cache: deps.Cache, Queue: deps.Queue},
	}
}
