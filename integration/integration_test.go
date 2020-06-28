package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"bradreed.co.uk/iverbs/api/options"
	"bradreed.co.uk/iverbs/api/server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/noisyscanner/gofly/gofly"
	"github.com/noisyscanner/gofly/migrate"
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

type TestConfigService struct{}

func (_ *TestConfigService) GetConfig() *gofly.DBConfig {
	return &gofly.DBConfig{
		Driver: "mysql",
		Host:   "mysql",
		User:   "root",
		Pass:   "iverbs",
		Port:   3306,
		Db:     "iverbs",
	}
}

var _ = Describe("Integration", func() {
	opts := &options.Options{
		Port:  6000,
		Redis: "redis:6379",
	}
	goflyConfig := &TestConfigService{}

	BeforeSuite(func() {
		err := migrate.Up(goflyConfig)
		if err != nil {
			panic(err)
		}
		// Expect(err).To(BeNil())
	})

	It("should work", func() {
		srv := server.GetServer(opts, goflyConfig)
		srv.Setup()

		req, err := http.NewRequest("GET", "/languages", nil)
		if err != nil {
			panic(err)
		}

		rr := httptest.NewRecorder()
		srv.Router.ServeHTTP(rr, req)

		fmt.Print(rr.Body)
	})
})
