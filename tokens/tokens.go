package tokens

import (
	"github.com/dimonomid/clock"
	"github.com/gomodule/redigo/redis"
	"time"
)

const DURATION = time.Second * 30

type TokenPersister interface {
	PersistToken(token string) error
}

type RedisTokenPersister struct {
	conn  redis.Conn
	clock clock.Clock
}

func (g *RedisTokenPersister) PersistToken(token string) error {
	now := g.clock.Now()
	expiryTime := now.Add(DURATION)

	_, err := g.conn.Do("HSET", KEY, token, expiryTime.Format(FORMAT))
	return err
}
