package tokens

import (
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const TOKEN = "sometoken"

var _ = Describe("Tokens", func() {
	var _ = Describe("RedisTokenPersister", func() {
		var _ = Describe("TokenPersister", func() {
			var (
				pool           *MockRedisPool
				conn           *MockRedisConn
				clockMock      *clock.Mock
				tokenGenerator *RedisTokenPersister
			)

			BeforeEach(func() {
				pool = new(MockRedisPool)
				conn = new(MockRedisConn)
				clockMock = clock.NewMock()

				pool.On("Get").Return(conn)
				conn.On("Close").Return(nil) // TODO: Test failure

				tokenGenerator = &RedisTokenPersister{
					pool:  pool,
					clock: clockMock,
				}
			})

			Describe("PersistToken", func() {
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

			Describe("InvalidateToken", func() {
				It("should delete the field from Redis", func() {
					conn.On("Do", "HDEL", KEY, TOKEN).Return(1, nil)
					err := tokenGenerator.InvalidateToken(TOKEN)
					Expect(err).To(BeNil())
				})

				It("should return an error if Redis returns an error", func() {
					redisError := fmt.Errorf("Redis died")
					conn.On("Do", "HDEL", KEY, TOKEN).Return(0, redisError)

					err := tokenGenerator.InvalidateToken(TOKEN)
					Expect(err).To(Equal(redisError))
				})
			})
		})
	})
})
