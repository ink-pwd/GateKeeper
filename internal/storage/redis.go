package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ClientRedis struct {
	client *redis.Client
	ctx    context.Context
}

func NewClient(redClient *redis.Client) *ClientRedis {
	return &ClientRedis{
		client: redClient,
		ctx:    context.Background(),
	}
}

func (r *ClientRedis) Ping() bool {
	var (
		err error
	)

	_, err = r.client.Ping(r.ctx).Result()

	if err != nil {
		return false
	}
	return true
}

func (r *ClientRedis) Set(uuid string, userJson []byte, minutes int) error {
	/*записываем данные пользователя, а в качестве ключа используем uuid
	Время жизни указываем в минутах в env файле*/
	var (
		duration time.Duration
	)

	duration = time.Duration(minutes) * time.Minute
	return r.client.Set(r.ctx, uuid, userJson, duration).Err()
}

func (r *ClientRedis) Get(uuid string) (string, error) {

	return r.client.Get(r.ctx, uuid).Result()

}

func (r *ClientRedis) Exist(uuid string) (int64, error) {

	return r.client.Exists(r.ctx, uuid).Result()
}

func (r *ClientRedis) Del(uuid string) error {
	return r.client.Del(r.ctx, uuid).Err()
}
