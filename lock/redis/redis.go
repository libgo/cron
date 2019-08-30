package redis

import (
	"github.com/go-redis/redis"
	"github.com/libgo/cron/lock"
	"github.com/libgo/logx"
)

func init() {
	lock.Register("redis", &Redis{})
}

type Redis struct {
	client *redis.Client
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

// Lock with SetNX(n)
func (m *Redis) Lock(n string) error {
	success, err := m.client.SetNX("job_lock_"+n, "Lock", 0).Result()
	if err != nil {
		logx.Errorf("cron: lock job %s error at redis: %s", n, err.Error())
		return lock.ErrLock
	}

	if !success {
		return lock.ErrLock
	}

	return nil
}

// Unlock with Del key(n)
func (m *Redis) Unlock(n string) error {
	d, err := m.client.Del("job_lock_" + n).Result()
	if err != nil {
		logx.Errorf("cron: unlock job %s error at redis: %s", n, err.Error())
		return lock.ErrUnlock
	}

	if d == 0 {
		logx.Errorf("cron: unlock job get wrong redis DEL return number")
		return lock.ErrUnlock
	}

	return nil
}
