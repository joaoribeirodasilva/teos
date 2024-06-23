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
	name     string
	client   *redis.Client
	addr     string
	port     int
	db       int
	username string
	password string
}

func New(name string, addr string, port int, database int, username string, password string) *RedisDB {

	r := &RedisDB{
		name:     name,
		addr:     addr,
		port:     port,
		db:       database,
		username: username,
		password: password,
	}

	return r
}

func (r *RedisDB) Connect() error {

	opts := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.addr, r.port),
		DB:       r.db,
		Username: r.username,
		Password: r.password,
	}

	slog.Info(fmt.Sprintf("service %s connecting to redis db on %s database %d", r.name, r.addr, r.db))
	r.client = redis.NewClient(&opts)

	_, err := r.client.Ping(context.Background()).Result()
	if err != nil {

		slog.Error(fmt.Sprintf("service %s failed to connect to redis server on '%s'. ERR: %s", r.name, r.addr, err.Error()))
		return err
	}

	slog.Info(fmt.Sprintf("redis service %s database %s connected", r.name, r.addr))

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

func (r *RedisDB) Get(key string) (*string, error) {

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

func (r *RedisDB) MGet(prefix string) (*map[string]*string, error) {

	var cursor uint64
	ctx := context.Background()
	items := make(map[string]*string)

	for {
		keys, cursor, err := r.client.Scan(ctx, cursor, prefix, 1000).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			val, err := r.Get(key)
			if err != nil {
				return nil, err
			}
			items[key] = val
		}

		if cursor == 0 {
			break
		}
	}

	return &items, nil
}
