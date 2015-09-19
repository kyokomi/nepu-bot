package main

import (
	"gopkg.in/redis.v2"

	"github.com/kyokomi/goroku"
	"github.com/kyokomi/slackbot"
	"golang.org/x/net/context"
)

type RedisRepository struct {
	redisDB *redis.Client
}

func NewRedisRepository() slackbot.Repository {
	s := &RedisRepository{}
	s.redisDB = goroku.MustRedis(goroku.OpenRedis(context.Background()))
	return s
}

func (s RedisRepository) Close() error {
	if s.redisDB != nil {
		return s.redisDB.Close()
	}
	return nil
}

func (r *RedisRepository) Save(key string, value string) error {
	return r.redisDB.Set(key, value).Err()
}

func (r *RedisRepository) Load(key string) (string, error) {
	return r.redisDB.Get(key).Result()
}

var _ slackbot.Repository = (*RedisRepository)(nil)
