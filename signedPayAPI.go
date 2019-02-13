package gosignedpayapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Client hold baseURI and virify flag
type Client struct {
	baseURI string
	verify  bool
}

// API store client info, mechant and private key
type API struct {
	client     Client
	merchantID []byte
	privateKey []byte
}

// NewAPI return new instance of Api
func NewAPI(merchantID []byte,
	privateKey []byte,
	baseURI string) *API {
	api := new(API)
	api.merchantID = merchantID
	api.privateKey = privateKey
	api.client.baseURI = baseURI
	api.client.verify = true
	return api
}

// Charge step of payment
func (a API) Charge(attributes []byte) ([]byte, error) {
	return a.sendRequest("charge", attributes)
}

// Status return status of paiment
func (a API) Status(attributes []byte) ([]byte, error) {
	return a.sendRequest("status", attributes)
}

// Refund return payment
func (a API) Refund(attributes []byte) ([]byte, error) {
	return a.sendRequest("refund", attributes)
}

// InitPayment start of payment initialisation
func (a API) InitPayment(attributes []byte) ([]byte, error) {
	return a.sendRequest("init-payment", attributes)
}

func (a API) sendRequest(method string, attributes []byte) ([]byte, error) {
	req := a.makeRequest(method, attributes)
	fmt.Println(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic("Error while request")
	}
	defer resp.Body.Close()
	fmt.Println(resp.Header)
	fmt.Println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return []byte{}, err
	}
	return body, nil
}

func (a *API) makeRequest(method string, attributes []byte) *http.Request {
	a.client.baseURI += method
	r, err := http.NewRequest("POST", a.client.baseURI, bytes.NewBuffer(attributes))
	if err != nil {
		log.Printf("Error creating request: %v", err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Merchant", string(a.merchantID))
	r.Header.Set("Signature", getSignature(attributes, a.merchantID, a.privateKey))
	return r
}

func getSignature(body []byte, m []byte, pk []byte) string {
	data := append(m, body...)
	data = append(data, m...)
	encData := hmac.New(sha512.New, pk)
	encData.Write(data)
	return base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(encData.Sum(nil))))
}
