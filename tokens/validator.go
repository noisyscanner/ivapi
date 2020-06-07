package tokens

import (
	"fmt"
	"time"

	"bradreed.co.uk/iverbs/api/options"
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

// TODO: Test
func NewRedisTokenValidator(options options.Options) (*RedisTokenValidator, error) {
	// TODO: in OptionProvider
	conn, err := redis.Dial("tcp", options.Redis)
	if err != nil {
		return nil, err
	}

	return &RedisTokenValidator{conn: conn}, err
}

func (v *RedisTokenValidator) Validate(token string) (bool, error) {
	expiry, err := v.conn.Do("HGET", v.key, token)
	if err != nil {
		return false, err
	}

	// TODO: Handle nil token

	var (
		expiryStr string
		ok        bool
	)
	if expiryStr, ok = expiry.(string); !ok {
		return false, fmt.Errorf("%v not a string", expiry)
	}

	expiryTime, err := time.Parse(FORMAT, expiryStr)
	if err != nil {
		return false, err
	}

	return expiryTime.After(time.Now()), nil
}

