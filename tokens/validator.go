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
	pool RedisPool
	key  string
}

func NewRedisTokenValidator(pool *redis.Pool) *RedisTokenValidator {
	return &RedisTokenValidator{
		pool: pool,
		key:  KEY, // TODO: In config?
	}
}

func (v *RedisTokenValidator) Validate(token string) (bool, error) {
	conn := v.pool.Get()
	defer conn.Close()

	expiry, err := redis.String(conn.Do("HGET", v.key, token))
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
