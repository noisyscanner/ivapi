package tokens

import (
	"fmt"
	"time"

	"github.com/dimonomid/clock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const TOKEN = "sometoken"

var _ = Describe("Tokens", func() {
	var _ = Describe("RedisTokenPersister", func() {
		var _ = Describe("TokenPersister", func() {
			var (
				conn           *MockRedisConn
				clockMock      *clock.Mock
				tokenGenerator *RedisTokenPersister
			)

			BeforeEach(func() {
				conn = new(MockRedisConn)
				clockMock = clock.NewMock()

				tokenGenerator = &RedisTokenPersister{
					conn:  conn,
					clock: clockMock,
				}
			})

			It("should save the token in Redis with an expiry time 30s in the future", func() {
				now := time.Now()
				clockMock.Set(now)

				expectedTime := now.Add(time.Second * 30).Format(FORMAT)
				conn.On("Do", "HSET", KEY, TOKEN, expectedTime).Return(nil, nil)

				err := tokenGenerator.PersistToken(TOKEN)
				Expect(err).To(BeNil())
			})

			It("should return an error if Redis returns an error", func() {
				now := time.Now()
				clockMock.Set(now)

				expectedTime := now.Add(time.Second * 30).Format(FORMAT)
				redisError := fmt.Errorf("Redis died")
				conn.On("Do", "HSET", KEY, TOKEN, expectedTime).Return(nil, redisError)

				err := tokenGenerator.PersistToken(TOKEN)
				Expect(err).To(Equal(redisError))
			})
		})
	})
})
