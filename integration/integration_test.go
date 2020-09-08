package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gomodule/redigo/redis"
	"github.com/noisyscanner/gofly/gofly"
	"github.com/noisyscanner/gofly/migrate"
	helpers "github.com/noisyscanner/ivapi/helpers"
	apihttp "github.com/noisyscanner/ivapi/http"
	"github.com/noisyscanner/ivapi/options"
	"github.com/noisyscanner/ivapi/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const KEY = "iverbs"

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
		Port:           6000,
		Redis:          helpers.GetEnvElse("REDIS", "localhost:6379"),
		CacheDirectory: "../testCache",
	}
	goflyConfig := GetTestConfig()
	var (
		srv *apihttp.Server
		rr  *httptest.ResponseRecorder
		db  *sql.DB
		err error
	)

	BeforeSuite(func() {
		db, err = setupTestDb(goflyConfig)
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

	BeforeEach(func() {
		srv = server.GetServer(opts, goflyConfig)
		srv.Setup()
		rr = httptest.NewRecorder()
	})

	getToken := func() (err error, token string) {
		req, _ := http.NewRequest("POST", "/tokens", nil)
		rrLocal := httptest.NewRecorder()
		srv.Router.ServeHTTP(rrLocal, req)

		tokenRes := &apihttp.TokenResponse{}
		json.Unmarshal(rrLocal.Body.Bytes(), tokenRes)

		if tokenRes.Error != "" {
			err = fmt.Errorf(tokenRes.Error)
		}

		token = tokenRes.Token
		return
	}

	Describe("POST /tokens", func() {
		It("should return a token and persist to Redis", func() {
			req, _ := http.NewRequest("POST", "/tokens", nil)
			srv.Router.ServeHTTP(rr, req)

			tokenRes := &apihttp.TokenResponse{}
			json.Unmarshal(rr.Body.Bytes(), tokenRes)

			Expect(tokenRes.Error).To(BeEmpty())

			redisClient, err := server.ConnectToRedis(opts)
			Expect(err).To(BeNil())
			defer redisClient.Close()

			expiryStr, err := redis.String(redisClient.Do("HGET", KEY, tokenRes.Token))
			Expect(err).To(BeNil())

			expiryTime, err := time.Parse(time.RFC3339, expiryStr)
			Expect(err).To(BeNil())
			Expect(expiryTime.Unix()).To(BeNumerically(">", time.Now().Unix()))
		})
	})

	Context("has language", func() {
		lang := &gofly.Language{
			Id:            1,
			Code:          "fr",
			Lang:          "French",
			Locale:        "fr_FR",
			HasHelpers:    true,
			HasReflexives: true,
		}

		tense := &gofly.Tense{
			Id:          1,
			Identifier:  "je",
			DisplayName: "Je",
			Order:       0,
		}

		pronoun := &gofly.Pronoun{
			Id:          1,
			Identifier:  "present",
			DisplayName: "Present",
			Order:       0,
		}

		verb := &gofly.Verb{
			Id:                   1,
			Infinitive:           "jour",
			NormalisedInfinitive: "jour",
			English:              "to play",
		}

		BeforeEach(func() {
			_, err := db.Exec("INSERT INTO languages (id, code, lang, locale, hasHelpers, hasReflexives) VALUES (?, ?, ?, ?, ?, ?)", lang.Id, lang.Code, lang.Lang, lang.Locale, lang.HasHelpers, lang.HasReflexives)
			if err != nil {
				panic(err)
			}

			_, err = db.Exec("INSERT INTO tenses (id, lang_id, identifier, displayName, `order`) VALUES (?, ?, ?, ?, ?)", tense.Id, lang.Id, tense.Identifier, tense.DisplayName, tense.Order)
			if err != nil {
				panic(err)
			}

			_, err = db.Exec("INSERT INTO pronouns (id, lang_id, identifier, displayName, `order`) VALUES (?, ?, ?, ?, ?)", pronoun.Id, lang.Id, pronoun.Identifier, pronoun.DisplayName, pronoun.Order)
			if err != nil {
				panic(err)
			}

			_, err = db.Exec("INSERT INTO verbs (id, lang_id, infinitive, normalisedInfinitive, english, helperID) VALUES (?, ?, ?, ?, ?, NULL)", verb.Id, lang.Id, verb.Infinitive, verb.NormalisedInfinitive, verb.English)
			if err != nil {
				panic(err)
			}
		})

		AfterEach(func() {
			_, err := db.Exec("DELETE FROM languages")
			if err != nil {
				panic(err)
			}
		})

		Describe("GET /languages", func() {
			It("should return the list of languages as expected", func() {
				req, _ := http.NewRequest("GET", "/languages", nil)
				srv.Router.ServeHTTP(rr, req)

				langRes := &apihttp.LanguagesResponse{}
				err := json.Unmarshal(rr.Body.Bytes(), langRes)
				Expect(err).To(BeNil())

				Expect(langRes.Error).To(BeEmpty())
				Expect(langRes.Data).To(HaveLen(1))

				version := langRes.Data[0].Version
				schemaVersion := langRes.Data[0].SchemaVersion
				Expect(int64(version)).To(BeNumerically("==", time.Now().Unix()))
				Expect(int64(schemaVersion)).To(BeNumerically("==", time.Now().Unix()))
				lang.Version = version
				lang.SchemaVersion = schemaVersion

				Expect(langRes.Data[0]).To(BeEquivalentTo(lang))
			})
		})

		Describe("GET /languages/{code}", func() {
			contents := `{"test": "json"}`
			contentsBytes := []byte(contents)

			var (
				err             error
				token           string
				gzContentsBytes []byte
			)

			BeforeEach(func() {
				err, token = getToken()
				Expect(err).To(BeNil())
			})

			Context("language exists", func() {
				BeforeEach(func() {
					err = os.Mkdir(opts.CacheDirectory, 0755)
					Expect(err).To(BeNil())

					cacheFile := opts.CacheDirectory + "/fr.json.full"
					err := ioutil.WriteFile(cacheFile, contentsBytes, 0644)
					Expect(err).To(BeNil())

					gzCacheFile := opts.CacheDirectory + "/fr.json.full.gz"
					gzContentsBytes, err = gofly.ZipBytes(contentsBytes)
					Expect(err).To(BeNil())
					err = ioutil.WriteFile(gzCacheFile, gzContentsBytes, 0644)
					Expect(err).To(BeNil())
				})

				AfterEach(func() {
					os.RemoveAll(opts.CacheDirectory)
				})

				It("should return an existing language", func() {
					req, _ := http.NewRequest("GET", "/languages/fr", nil)
					req.Header.Add("Authorization", token)
					srv.Router.ServeHTTP(rr, req)

					Expect(rr.Body.String()).To(Equal(contents))
				})

				It("should return an existing language as GZIP", func() {
					req, _ := http.NewRequest("GET", "/languages/fr", nil)
					req.Header.Add("Authorization", token)
					req.Header.Add("Accept-Encoding", "gzip")
					srv.Router.ServeHTTP(rr, req)

					Expect(rr.Body.Bytes()).To(Equal(gzContentsBytes))
				})
			})

			Context("language does not exist", func() {
				It("should return an error if the cache file does not exist", func() {
					req, _ := http.NewRequest("GET", "/languages/fr", nil)
					req.Header.Add("Authorization", token)
					srv.Router.ServeHTTP(rr, req)

					errRes := &apihttp.LanguagesResponse{}
					json.Unmarshal(rr.Body.Bytes(), errRes)

					Expect(errRes.Error).To(Equal("Language not found"))
				})
			})
		})
	})
})
