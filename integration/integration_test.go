package integration

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/noisyscanner/gofly/gofly"
	apihttp "github.com/noisyscanner/ivapi/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const KEY = "iverbs"

func insertLang(lang *gofly.Language) {
	_, err := db.Exec("INSERT INTO languages (id, code, lang, locale, hasHelpers, hasReflexives) VALUES (?, ?, ?, ?, ?, ?)", lang.Id, lang.Code, lang.Lang, lang.Locale, lang.HasHelpers, lang.HasReflexives)
	Expect(err).To(BeNil())
}

func insertTense(lang *gofly.Language, tense *gofly.Tense) {
	_, err := db.Exec("INSERT INTO tenses (id, lang_id, identifier, displayName, `order`) VALUES (?, ?, ?, ?, ?)", tense.Id, lang.Id, tense.Identifier, tense.DisplayName, tense.Order)
	Expect(err).To(BeNil())
}

func insertPronoun(lang *gofly.Language, pronoun *gofly.Pronoun) {
	_, err := db.Exec("INSERT INTO pronouns (id, lang_id, identifier, displayName, `order`) VALUES (?, ?, ?, ?, ?)", pronoun.Id, lang.Id, pronoun.Identifier, pronoun.DisplayName, pronoun.Order)
	Expect(err).To(BeNil())
}

func insertVerb(lang *gofly.Language, verb *gofly.Verb) {
	_, err := db.Exec("INSERT INTO verbs (id, lang_id, infinitive, normalisedInfinitive, english, helperID) VALUES (?, ?, ?, ?, ?, NULL)", verb.Id, lang.Id, verb.Infinitive, verb.NormalisedInfinitive, verb.English)
	Expect(err).To(BeNil())
}

var _ = Describe("Integration: languages", func() {
	BeforeEach(func() {
		srv, rr = SetupServer()
	})

	Context("has language", func() {
		BeforeEach(func() {
			insertLang(frenchLang)
			insertTense(frenchLang, frenchTense)
			insertPronoun(frenchLang, frenchPronoun)
			insertVerb(frenchLang, frenchVerb)

			insertLang(germanLang)
			insertTense(germanLang, germanTense)
			insertPronoun(germanLang, germanPronoun)
			insertVerb(germanLang, germanVerb)
		})

		AfterEach(func() {
			_, err := db.Exec("DELETE FROM languages")
			Expect(err).To(BeNil())
		})

		Describe("GET /languages", func() {
			It("should return the list of languages as expected", func() {
				req, _ := http.NewRequest("GET", "/languages", nil)
				srv.Router.ServeHTTP(rr, req)

				langRes := &apihttp.LanguagesResponse{}
				err := json.Unmarshal(rr.Body.Bytes(), langRes)
				Expect(err).To(BeNil())

				Expect(langRes.Error).To(BeEmpty())
				Expect(langRes.Data).To(HaveLen(2))

				version := langRes.Data[0].Version
				schemaVersion := langRes.Data[0].SchemaVersion
				Expect(int64(version)).To(BeNumerically("==", time.Now().Unix()))
				Expect(int64(schemaVersion)).To(BeNumerically("==", time.Now().Unix()))
				frenchLang.Version = version
				frenchLang.SchemaVersion = schemaVersion

				Expect(langRes.Data[0]).To(BeEquivalentTo(frenchLang))
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
				token = GetToken(srv)
			})

			It("should invalidate the token", func() {
				req, _ := http.NewRequest("GET", "/languages/fr", nil)
				req.Header.Add("Authorization", token)
				srv.Router.ServeHTTP(rr, req)

				redisClient := redisPool.Get()
				defer redisClient.Close()
				_, err := redis.String(redisClient.Do("HGET", KEY, token))
				Expect(err).To(Equal(redis.ErrNil))
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
