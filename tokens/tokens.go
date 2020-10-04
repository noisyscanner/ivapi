package tokens

import (
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gomodule/redigo/redis"
)

const KEY = "iverbs"
const DURATION = time.Second * 30

type TokenPersister interface {
	PersistToken(token string) error
	InvalidateToken(token string) error
}

type RedisPool interface {
	Get() redis.Conn
}

type RedisTokenPersister struct {
	pool  RedisPool
	clock clock.Clock
}

func (p *RedisTokenPersister) PersistToken(token string) error {
	now := p.clock.Now()
	expiryTime := now.Add(DURATION)

	conn := p.pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", KEY, token, expiryTime.Format(FORMAT))
	return err
}

func (p *RedisTokenPersister) InvalidateToken(token string) error {
	conn := p.pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", KEY, token)
	return err
}

func NewRedisTokenPersister(pool *redis.Pool) *RedisTokenPersister {
	return &RedisTokenPersister{
		pool:  pool,
		clock: clock.New(),
	}
}
