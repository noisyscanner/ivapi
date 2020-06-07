package tokens

import (
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/mock"
)

type MockRedisConn struct {
	mock.Mock
	redis.Conn
}

func (m *MockRedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	argArr := []interface{}{commandName}
	for _, arg := range args {
		argArr = append(argArr, arg)
	}

	returns := m.Called(argArr...)
	return returns.Get(0), returns.Error(1)
}
