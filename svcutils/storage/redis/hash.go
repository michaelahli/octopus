package redis

import (
	"context"
	"reflect"

	"github.com/go-redis/redis/v8"
)

type Hash interface {
	HMGet(ctx context.Context) error
	HMSet(ctx context.Context) error
	HDel(ctx context.Context) error
}

type redisHash struct {
	rst  *redisStorage
	hash string
}

// HMGet() : Get using hashed key
// hash can contains multiple keys
func (rh *redisHash) HMGet(ctx context.Context) error {
	if reflect.ValueOf(rh.rst.data).Kind() != reflect.Ptr {
		return ErrNonDataPtr
	}

	err := rh.rst.rdb.HGet(ctx, rh.hash, rh.rst.key.ToString()).Scan(rh.rst.data)
	if err != nil {
		if err == redis.Nil {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// HMSet() : Set using hashed key
// hash can contains multiple keys
func (rh *redisHash) HMSet(ctx context.Context) error {
	_, err := rh.rst.rdb.HSet(ctx, rh.hash, rh.rst.key.ToString()).Result()
	if err != nil {
		return err
	}
	return nil
}

// HMSet() : Delete using hashed key
// hash can contains multiple keys
func (rh *redisHash) HDel(ctx context.Context) error {
	_, err := rh.rst.rdb.HDel(ctx, rh.hash, rh.rst.key.ToString()).Result()
	if err != nil {
		return err
	}
	return nil
}
