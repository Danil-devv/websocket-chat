package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Opt    *redis.Options
	Key    string
	Logger logrus.FieldLogger
}
