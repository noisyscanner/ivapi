package tokens

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tokens", func() {
	var _ = Describe("RedisTokenValidator", func() {
		var pool *MockRedisPool
		var conn *MockRedisConn
		var redisTokenValidator *RedisTokenValidator

		BeforeEach(func() {
			pool = new(MockRedisPool)
			conn = new(MockRedisConn)
			pool.On("Get").Return(conn)
			conn.On("Close").Return(nil) // TODO: Test failure

			redisTokenValidator = &RedisTokenValidator{
				pool: pool,
				key:  KEY,
			}
		})

		var _ = Describe("Validate", func() {
			token := "sometoken"

			It("should return an error if Redis returns an error", func() {
				errToReturn := fmt.Errorf("Redis error")
				conn.On("Do", "HGET", KEY, token).Return(nil, errToReturn)

				isValid, err := redisTokenValidator.Validate(token)
				Expect(err).To(Equal(errToReturn))
				Expect(isValid).To(BeFalse())
			})

			It("should return an error if the expiry time is not a string", func() {
				malformedExpiry := 12345
				conn.On("Do", "HGET", KEY, token).Return(malformedExpiry, nil)

				isValid, err := redisTokenValidator.Validate(token)
				Expect(err).NotTo(BeNil())
				Expect(isValid).To(BeFalse())
			})

			It("should return an error if the expiry time is the wrong format", func() {
				malformedExpiry := "wrongformat"
				conn.On("Do", "HGET", KEY, token).Return(malformedExpiry, nil)

				isValid, err := redisTokenValidator.Validate(token)
				Expect(err).NotTo(BeNil())
				Expect(isValid).To(BeFalse())
			})

			It("should return false if the token has expired", func() {
				tokenDuration, _ := time.ParseDuration("30s")
				expiry := time.Now().Add(-tokenDuration).Format(time.RFC3339)

				conn.On("Do", "HGET", KEY, token).Return(expiry, nil)

				isValid, err := redisTokenValidator.Validate(token)
				Expect(err).To(BeNil())
				Expect(isValid).To(BeFalse())
			})

			It("should return true if the token expires in the future", func() {
				tokenDuration, _ := time.ParseDuration("30s")
				expiry := time.Now().Add(tokenDuration).Format(time.RFC3339)

				conn.On("Do", "HGET", KEY, token).Return(expiry, nil)

				isValid, err := redisTokenValidator.Validate(token)
				Expect(err).To(BeNil())
				Expect(isValid).To(BeTrue())
			})
		})
	})
})
