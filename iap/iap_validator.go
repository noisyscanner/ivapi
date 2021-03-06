package iap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const SANDBOX_URL = "https://sandbox.itunes.apple.com"
const PROD_URL = "https://buy.itunes.apple.com"
const RESPONSE_CODE_SANDBOX = 21007
const RESPONSE_CODE_VALID = 0

type AppStoreResponse struct {
	Status int
}

type HttpClient interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type IapValidator struct {
	Client HttpClient
}

func NewIapValidator() *IapValidator {
	return &IapValidator{
		Client: http.DefaultClient,
	}
}

func getItunesUrl(isSandbox bool) string {
	var host string
	if isSandbox {
		host = SANDBOX_URL
	} else {
		host = PROD_URL
	}
	return fmt.Sprintf("%s/verifyReceipt", host)
}

func (v *IapValidator) validateIapToken(receipt []byte, isSandbox bool) (isValid bool, err error) {
	payload := map[string]string{
		"receipt-data": string(receipt),
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Print("Could not encode App Store request body", err)
		return
	}

	url := getItunesUrl(isSandbox)
	res, err := v.Client.Post(url, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		log.Print("Request to App Store failed", err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("Could not read App Store response", err)
		return
	}

	resJson := &AppStoreResponse{}
	err = json.Unmarshal(body, resJson)
	if err != nil {
		log.Print("Could not unmarshal App Store response", err)
		return
	}

	if resJson.Status == RESPONSE_CODE_SANDBOX {
		return v.validateIapToken(receipt, true)
	}

	isValid = resJson.Status == RESPONSE_CODE_VALID
	return
}

func (v *IapValidator) ValidateIapToken(receipt []byte) (bool, error) {
	return v.validateIapToken(receipt, false)
}
