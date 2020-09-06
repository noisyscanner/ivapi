package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	helpers "bradreed.co.uk/iverbs/api/helpers"
	apihttp "bradreed.co.uk/iverbs/api/http"
	"bradreed.co.uk/iverbs/api/options"
	"bradreed.co.uk/iverbs/api/server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/noisyscanner/gofly/gofly"
	"github.com/noisyscanner/gofly/migrate"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO migrations path

type TestConfigService struct {
	Driver string
	Host   string
	User   string
	Pass   string
	Port   int
	Db     string
}

func GetTestConfig() *TestConfigService {
	port, _ := strconv.Atoi(helpers.GetEnvElse("DB_PORT", "3306"))
	return &TestConfigService{
		Driver: helpers.GetEnvElse("DB_DRIVER", "mysql"),
		Host:   helpers.GetEnvElse("DB_HOST", "localhost"),
		User:   helpers.GetEnvElse("DB_USER", "root"),
		Pass:   helpers.GetEnvElse("DB_PASS", "iverbs"),
		Port:   port,
		Db:     helpers.GetEnvElse("DB_NAME", "ivapi_1"),
	}
}

func (c *TestConfigService) GetConfig() *gofly.DBConfig {
	return &gofly.DBConfig{
		Driver: c.Driver,
		Host:   c.Host,
		User:   c.User,
		Pass:   c.Pass,
		Port:   c.Port,
		Db:     c.Db,
	}
}

func (c *TestConfigService) ConnString() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/", c.User, c.Pass, c.Host, c.Port)
}

func setupTestDb(goflyConfig *TestConfigService) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", goflyConfig.ConnString())
	if err != nil {
		return
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS " + goflyConfig.Db)
	if err != nil {
		return
	}

	_, err = db.Exec("CREATE DATABASE " + goflyConfig.Db)
	if err != nil {
		return
	}

	db.Close()

	db, err = sql.Open("mysql", goflyConfig.GetConfig().DBString())
	if err != nil {
		return
	}

	err = migrate.Up(goflyConfig)
	return
}

var _ = Describe("Integration", func() {
	opts := &options.Options{
		Port:  6000,
		Redis: helpers.GetEnvElse("REDIS", "localhost:6379"),
	}
	goflyConfig := GetTestConfig()

	expLang := &gofly.Language{
		Id:     1,
		Code:   "fr",
		Lang:   "French",
		Locale: "fr_FR",
	}

	BeforeSuite(func() {
		db, err := setupTestDb(goflyConfig)

		_, err = db.Exec("INSERT INTO languages (id, code, lang, locale) VALUES (?, ?, ?, ?)", expLang.Id, expLang.Code, expLang.Lang, expLang.Locale)
		if err != nil {
			panic(err)
		}
	})

	AfterSuite(func() {
		err := migrate.Down(goflyConfig)
		if err != nil {
			panic(err)
		}
	})

	It("should return the list of languages as expected", func() {
		srv := server.GetServer(opts, goflyConfig)
		srv.Setup()

		req, err := http.NewRequest("GET", "/languages", nil)
		if err != nil {
			panic(err)
		}

		rr := httptest.NewRecorder()
		srv.Router.ServeHTTP(rr, req)

		langRes := &apihttp.LanguagesResponse{}
		json.Unmarshal(rr.Body.Bytes(), langRes)

		Expect(langRes.Error).To(BeEmpty())
		Expect(langRes.Data).To(HaveLen(1))
		Expect(langRes.Data[0]).To(BeEquivalentTo(expLang))
	})
})
