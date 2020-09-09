package integration

import (
	"database/sql"
	"net/http/httptest"
	"testing"

	"github.com/gomodule/redigo/redis"
	apihttp "github.com/noisyscanner/ivapi/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var (
	srv         *apihttp.Server
	rr          *httptest.ResponseRecorder
	redisClient redis.Conn
	db          *sql.DB
)

var _ = Describe("Setup", func() {
	BeforeSuite(func() {
		db, redisClient = Setup()
	})

	AfterSuite(func() {
		TearDown(db, redisClient)
	})
})
