package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var ReddisClient = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("RedisAddr"),
	Password: os.Getenv("RedisPass"), // no password set
	DB:       0,                      // use default DB
})
