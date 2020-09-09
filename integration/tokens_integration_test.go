package integration

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	apihttp "github.com/noisyscanner/ivapi/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration: tokens", func() {
	BeforeEach(func() {
		srv, rr = SetupServer()
	})

	Describe("POST /tokens", func() {
		It("should return a token and persist to Redis", func() {
			req, _ := http.NewRequest("POST", "/tokens", nil)
			srv.Router.ServeHTTP(rr, req)

			tokenRes := &apihttp.TokenResponse{}
			json.Unmarshal(rr.Body.Bytes(), tokenRes)

			Expect(tokenRes.Error).To(BeEmpty())

			expiryStr, err := redis.String(redisClient.Do("HGET", KEY, tokenRes.Token))
			Expect(err).To(BeNil())

			expiryTime, err := time.Parse(time.RFC3339, expiryStr)
			Expect(err).To(BeNil())
			Expect(expiryTime.Unix()).To(BeNumerically(">", time.Now().Unix()))
		})
	})
})
