package main

import (
	"flag"
	"log"
	"twitch_chat_analysis/internal/data"

	amqp "github.com/rabbitmq/amqp091-go"
)

type config struct {
	port     int
	env      string
	queueURL string
}
type application struct {
	models data.Models
	config config
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4001, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.queueURL, "queueURL", "amqp://user:password@localhost:7001/", "RabbitMQ connection url")

	flag.Parse()

	conn, err := amqp.Dial(cfg.queueURL)
	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err.Error())
	}
	defer conn.Close()

	dependencies := data.Dependencies{
		Queue: conn,
	}

	app := &application{
		models: data.NewModels(dependencies),
		config: cfg,
	}
	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
