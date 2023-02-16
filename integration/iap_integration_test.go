package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	apihttp "github.com/noisyscanner/ivapi/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

type FakeClient struct {
	mock.Mock
}

func (m *FakeClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	args := m.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func response(respBody string) *http.Response {
	return &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(respBody)),
		ContentLength: int64(len(respBody)),
		Header:        make(http.Header, 0),
	}
}

func successResponse() (resp *http.Response) {
	respBody := `{"status": 0}`
	return response(respBody)
}

func sandBoxResponse() (resp *http.Response) {
	respBody := `{"status": 21007}`
	return response(respBody)
}

var _ = Describe("Integration: tokens", func() {
	var fakeClient *FakeClient
	BeforeEach(func() {
		srv, rr = SetupServer()
		fakeClient = new(FakeClient)
		srv.ReceiptValidator.Client = fakeClient
	})

	Describe("POST /iapvalidate", func() {
		receipt := []byte("sdfsf")
		receiptBody := &apihttp.ReceiptBody{
			Receipt: string(receipt),
		}
		receiptBodyStr, _ := receiptBody.MarshalJSON()
		expectedAppleBody := []byte(fmt.Sprintf(`{"receipt-data":"%s"}`, receipt))
		var token string

		BeforeEach(func() {
			token = GetToken(srv)
		})

		It("should validate the iap receipt (prod)", func() {
			fakeClient.
				On(
					"Post",
					"https://buy.itunes.apple.com/verifyReceipt",
					"application/json",
					mock.MatchedBy(func(reader io.Reader) bool {
						body, err := ioutil.ReadAll(reader)
						return err == nil && bytes.Equal(body, expectedAppleBody)
					}),
				).
				Return(successResponse(), nil).
				Once()

			req, _ := http.NewRequest(http.MethodPost, "/iapvalidate", bytes.NewReader(receiptBodyStr))
			req.Header.Set("Authorization", token)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Length", strconv.Itoa(len(receiptBodyStr)))
			srv.Router.ServeHTTP(rr, req)

			iapRes := &apihttp.IapResponse{}
			body := rr.Body.Bytes()
			json.Unmarshal(body, iapRes)

			Expect(iapRes.Error).To(BeEmpty())
			Expect(iapRes.Success).To(BeTrue())
		})

		It("should validate the iap receipt (sandbox)", func() {
			fakeClient.
				On(
					"Post",
					"https://buy.itunes.apple.com/verifyReceipt",
					"application/json",
					mock.MatchedBy(func(reader io.Reader) bool {
						body, err := ioutil.ReadAll(reader)
						return err == nil && bytes.Equal(body, expectedAppleBody)
					}),
				).
				Return(sandBoxResponse(), nil)

			fakeClient.
				On(
					"Post",
					"https://sandbox.itunes.apple.com/verifyReceipt",
					"application/json",
					mock.Anything, // anything as we can't read from a reader twice
				).
				Return(successResponse(), nil)

			req, _ := http.NewRequest("POST", "/iapvalidate", bytes.NewReader(receiptBodyStr))
			req.Header.Set("Authorization", token)
			srv.Router.ServeHTTP(rr, req)

			iapRes := &apihttp.IapResponse{}
			body := rr.Body.Bytes()
			json.Unmarshal(body, iapRes)

			Expect(iapRes.Error).To(BeEmpty())
			Expect(iapRes.Success).To(BeTrue())
		})
	})
})
