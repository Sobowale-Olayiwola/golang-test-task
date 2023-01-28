package main

import (
	"context"
	"flag"
	"log"
	"twitch_chat_analysis/internal/data"

	"github.com/redis/go-redis/v9"
)

type config struct {
	port     int
	env      string
	cacheURL string
}

type application struct {
	models data.Models
	config config
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (testing|development|staging|production)")
	flag.StringVar(&cfg.cacheURL, "cacheURL", "localhost:6379", "Redis connection url")

	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.cacheURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.TODO()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to Redis:", err.Error())
	}
	deps := data.Dependencies{
		Cache: rdb,
	}

	app := &application{
		models: data.NewModels(deps),
		config: cfg,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
