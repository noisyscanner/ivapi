package tokens

import (
	"github.com/benbjohnson/clock"
	"github.com/gomodule/redigo/redis"
	"time"
)

const KEY = "iverbs"
const DURATION = time.Second * 30

type TokenPersister interface {
	PersistToken(token string) error
	InvalidateToken(token string) error
}

type RedisTokenPersister struct {
	conn  redis.Conn
	clock clock.Clock
}

func (p *RedisTokenPersister) PersistToken(token string) error {
	now := p.clock.Now()
	expiryTime := now.Add(DURATION)

	_, err := p.conn.Do("HSET", KEY, token, expiryTime.Format(FORMAT))
	return err
}

func (p *RedisTokenPersister) InvalidateToken(token string) error {
	_, err := p.conn.Do("HDEL", KEY, token)
	return err
}

func NewRedisTokenPersister(conn redis.Conn) *RedisTokenPersister {
	return &RedisTokenPersister{
		conn:  conn,
		clock: clock.New(),
	}
}

