package redisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Name     string
	Client   *redis.Client
	Addr     string
	Db       int
	Username string
	Password string
}

func New(name string, addr string, db int, username string, password string) *RedisDB {
	r := new(RedisDB)
	r.Name = name
	r.Client = nil
	r.Addr = addr
	r.Db = db
	r.Username = username
	r.Password = password
	return r
}

func (r *RedisDB) Connect() error {

	opts := redis.Options{
		Addr:     r.Addr,
		DB:       r.Db,
		Username: r.Username,
		Password: r.Password,
	}

	slog.Info(fmt.Sprintf("[COMMON::REDISDB::Connect] connecting to redis '%s' database %d at %s...\n", r.Name, opts.DB, opts.Addr))
	r.Client = redis.NewClient(&opts)

	_, err := r.Client.Ping(context.Background()).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("[COMMON::REDISDB::Connect] failed to connect to '%s' with Redis server(s). ERR: %s", r.Name, err.Error()))
		return err
	}
	slog.Info(fmt.Sprintf("[COMMON::REDISDB::Connect] redis database %s connected", opts.Addr))
	return nil
}

func (r *RedisDB) Set(key string, val interface{}, ttl int) error {

	json, err := json.Marshal(val)
	if err != nil {
		slog.Error(fmt.Sprintf("[COMMON::REDISDB::Set] failed to marshal data to JSON. ERR: %s", err.Error()))
		return err
	}

	duration := int64(time.Second) * int64(ttl)
	if err := r.Client.Set(context.Background(), key, string(json), time.Duration(duration)); err.Err() != nil {
		slog.Error(fmt.Sprintf("[COMMON::REDISDB::Set] failed to set to '%s' Redis server. ERR: %s", r.Name, err.Err().Error()))
		return err.Err()
	}

	return nil
}

func (r *RedisDB) Get(key string, val interface{}, ttl int) (*string, error) {

	result, err := r.Client.Get(context.Background(), "rec").Result()
	if err != nil {
		slog.Error(fmt.Sprintf("[COMMON::REDISDB::Get] failed to set to '%s' Redis server. ERR: %s", r.Name, err.Error()))
		return nil, err
	}

	return &result, nil
}

func (r *RedisDB) Del(keys ...string) error {

	_, err := r.Client.Del(context.Background(), keys...).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("[COMMON::REDISDB::Get] failed to delete from '%s' Redis server. ERR: %s", r.Name, err.Error()))
		return err
	}
	return nil
}
