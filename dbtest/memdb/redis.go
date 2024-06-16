package memdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	client   *redis.Client
	addr     string
	db       int
	username string
	password string
}

func (r *RedisDB) Connect() error {

	opts := redis.Options{
		Addr:     r.addr,
		DB:       r.db,
		Username: r.username,
		Password: r.password,
	}

	slog.Info(fmt.Sprintf("connecting to redis db on %s database %d", r.addr, r.db))
	r.client = redis.NewClient(&opts)

	_, err := r.client.Ping(context.Background()).Result()
	if err != nil {

		slog.Error(fmt.Sprintf("failed to connect to redis server on '%s'. ERR: %s", r.addr, err.Error()))
		return err
	}

	slog.Info(fmt.Sprintf("redis database %s connected", r.addr))

	return nil
}

func (r *RedisDB) Set(key string, val interface{}, ttl int) error {

	json, err := json.Marshal(val)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to marshal data to JSON. ERR: %s", err.Error()))
		return err
	}

	duration := int64(time.Second) * int64(ttl)
	if err := r.client.Set(context.Background(), key, string(json), time.Duration(duration)); err.Err() != nil {
		slog.Error(fmt.Sprintf("failed to set to  Redis server on '%s'. ERR: %s", r.addr, err.Err().Error()))
		return err.Err()
	}

	return nil
}

func (r *RedisDB) Get(key string, val interface{}, ttl int) (*string, error) {

	result, err := r.client.Get(context.Background(), "rec").Result()
	if err != nil {
		slog.Error(fmt.Sprintf("[COMMON::REDISDB::Get] failed to set to redis server on '%s'. ERR: %s", r.addr, err.Error()))
		return nil, err
	}

	return &result, nil
}

func (r *RedisDB) Del(keys ...string) error {

	_, err := r.client.Del(context.Background(), keys...).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("[COMMON::REDISDB::Get] failed to delete from redis server on '%s'. ERR: %s", r.addr, err.Error()))
		return err
	}

	return nil
}
