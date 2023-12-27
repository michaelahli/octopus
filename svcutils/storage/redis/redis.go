package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type redisClient struct {
	rdb *redis.Client
}

type Redis interface {
	Client() *redis.Client
	Prepare(key Keys, data any, exp time.Duration) (Storage, error)
	Close() error
}

func NewConnection(cfg *viper.Viper, log *logrus.Entry, dbnum dbloc) (Redis, error) {
	connString := fmt.Sprintf("%s:%s", cfg.GetString("redis.host"), cfg.GetString("redis.port"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     connString,
		Password: cfg.GetString("redis.password"),
		DB:       dbnum.ToInt(),
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.WithError(err).Error("redis connection")
		return nil, err
	}

	log.Infoln("Successfully established connection to", connString)
	log.Infof("Ping redis client : %v \n", pong)

	return &redisClient{rdb: rdb}, nil
}

func (rc *redisClient) Prepare(key Keys, data any, exp time.Duration) (Storage, error) {
	if data == nil {
		return nil, ErrNilDataStorage
	}
	return &redisStorage{rdb: rc.rdb, key: key, data: data, exp: exp}, nil
}

func (rc *redisClient) Close() error {
	return rc.rdb.Close()
}

func (rc *redisClient) Client() *redis.Client {
	return rc.rdb
}
