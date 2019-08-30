package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/libgo/cron/lock"
)

func init() {
	lock.Register("Redis", &Redis{})
}

type Redis struct {
	client *redis.Client
}

// Conf redis
type Conf struct {
	DSN string
}

// Open a new mysql locker instance
// uri format -> redis://:password@url/dbNum[optional,default 0]
func (m *Redis) Open(uri string) lock.Locker {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	return &Redis{client: client}
}

// Lock with setNX(n)
func (m *Redis) Lock(n string) error {
	success, err := m.client.SetNX(n,"Lock", 10*time.Second).Result()
	if !success {
		return fmt.Errorf("lock job %s error", n)
	}
	return err
}

// Unlock with del key(n)
func (m *Redis) Unlock(n string) error {
	_, err :=m.client.Del(n).Result()
	return err
}
