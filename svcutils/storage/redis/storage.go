package redis

import (
	"context"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisStorage struct {
	rdb  *redis.Client
	key  Keys
	exp  time.Duration
	data any
}

type Storage interface {
	Get(ctx context.Context) error
	Set(ctx context.Context) error
	Del(ctx context.Context) error
	Ttl(ctx context.Context) (time.Duration, error)
	WithHash(ctx context.Context, hash string) Hash
}

func (rst *redisStorage) Get(ctx context.Context) error {
	if reflect.ValueOf(rst.data).Kind() != reflect.Ptr {
		return ErrNonDataPtr
	}

	err := rst.rdb.Get(ctx, rst.key.ToString()).Scan(rst.data)
	if err != nil {
		if err == redis.Nil {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (rst *redisStorage) Set(ctx context.Context) error {
	_, err := rst.rdb.Set(ctx, rst.key.ToString(), rst.data, rst.exp).Result()
	if err != nil {
		return err
	}
	return nil
}

func (rst *redisStorage) Del(ctx context.Context) error {
	_, err := rst.rdb.Del(ctx, rst.key.ToString()).Result()
	if err != nil {
		return err
	}
	return nil
}

func (rst *redisStorage) WithHash(ctx context.Context, hash string) Hash {
	return &redisHash{rst: rst, hash: hash}
}

func (rst *redisStorage) Ttl(ctx context.Context) (time.Duration, error) {
	dur, err := rst.rdb.TTL(ctx, rst.key.ToString()).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return dur, err
}
