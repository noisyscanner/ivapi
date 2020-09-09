package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/noisyscanner/gofly/gofly"
	"github.com/noisyscanner/gofly/migrate"
	helpers "github.com/noisyscanner/ivapi/helpers"
	apihttp "github.com/noisyscanner/ivapi/http"
	"github.com/noisyscanner/ivapi/options"
	"github.com/noisyscanner/ivapi/server"
	. "github.com/onsi/gomega"
)

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

var opts = &options.Options{
	Port:           6000,
	Redis:          helpers.GetEnvElse("REDIS", "localhost:6379"),
	CacheDirectory: "../testCache",
}
var goflyConfig = GetTestConfig()

func Setup() (*sql.DB, redis.Conn) {
	db, err := setupTestDb(goflyConfig)
	Expect(err).To(BeNil())

	redisClient, err := server.ConnectToRedis(opts)
	Expect(err).To(BeNil())

	return db, redisClient
}

func TearDown(db *sql.DB, redisClient redis.Conn) {
	err := migrate.Down(goflyConfig)
	Expect(err).To(BeNil())

	redisClient.Close()
}

func SetupServer() (srv *apihttp.Server, rr *httptest.ResponseRecorder) {
	srv = server.GetServer(opts, goflyConfig)
	srv.Setup()
	rr = httptest.NewRecorder()
	return
}

func GetToken(srv *apihttp.Server) string {
	req, _ := http.NewRequest("POST", "/tokens", nil)
	rrLocal := httptest.NewRecorder()
	srv.Router.ServeHTTP(rrLocal, req)

	tokenRes := &apihttp.TokenResponse{}
	json.Unmarshal(rrLocal.Body.Bytes(), tokenRes)

	Expect(tokenRes.Error).To(BeEmpty())

	return tokenRes.Token
}
