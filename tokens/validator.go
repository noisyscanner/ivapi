package tokens

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

const FORMAT = time.RFC3339

type TokenValidator interface {
	Validate(token string) (bool, error)
}

type RedisTokenValidator struct {
	conn redis.Conn
	key  string
}

func NewRedisTokenValidator(conn redis.Conn) *RedisTokenValidator {
	return &RedisTokenValidator{
		conn: conn,
		key:  KEY, // TODO: In config?
	}
}

func (v *RedisTokenValidator) Validate(token string) (bool, error) {
	expiry, err := redis.String(v.conn.Do("HGET", v.key, token))
	if err != nil {
		// TODO: Test
		if err == redis.ErrNil {
			err = nil
		}
		return false, err
	}

	expiryTime, err := time.Parse(FORMAT, expiry)
	if err != nil {
		return false, err
	}

	return expiryTime.After(time.Now()), nil
}

