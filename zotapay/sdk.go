package zotapay

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

const (
	// VERSION string of SDK
	VERSION string = "v1.1.2"
	// SANDBOX endpoint URL
	SANDBOX string = "https://api.zotapay-sandbox.com"
	// LIVE production endpoint URL
	LIVE string = "https://api.zotapay.com"
)

// SDK represents the base SDK structure
// all properties are required, except HttpClient
// HttpClient implement httpClient interface if is empty will be initialized
type SDK struct {
	MerchantID        string
	MerchantSecretKey string
	EndpointID        string
	ApiBaseURL        string
	HttpClient        httpClient
}

// httpClient is the interface that wraps the basic http.Client Do method.
type httpClient interface {
	Do(req *http.Request) (ret *http.Response, err error)
}

// validate the instance of SDK
// and returns an error
func (s *SDK) validate() error {
	if s.MerchantID == "" {
		return fmt.Errorf("MerchantID is required.")
	}
	if s.MerchantSecretKey == "" {
		return fmt.Errorf("MerchantSecretKey is required.")
	}
	if s.EndpointID == "" {
		return fmt.Errorf("EndpointID is required.")
	}
	if s.ApiBaseURL == "" {
		return fmt.Errorf("ApiBaseURL is required.")
	}
	if s.ApiBaseURL == "" {
		return fmt.Errorf("ApiBaseURL is required.")
	}
	return nil
}

// initHttpClient is checking if *SDK.HttpClient has been populated
// in case it has not, crete new one.
func (s *SDK) initHttpClient() {
	if s.HttpClient == nil {
		//enforce tls min version > 1.2
		mTLSConfig := &tls.Config{}
		mTLSConfig.PreferServerCipherSuites = true
		mTLSConfig.MinVersion = tls.VersionTLS12
		tr := &http.Transport{
			TLSClientConfig: mTLSConfig,
		}

		s.HttpClient = &http.Client{
			Timeout:   time.Second * 10,
			Transport: tr,
		}
	}
}

// postJson makes an http-post request to a Zotapay API endpoint.
// returns the response as []byte, or an error.
func (s *SDK) httpDo(method string, url string, data []byte) (code int, body []byte, err error) {

	s.initHttpClient()

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("Zotapay Go SDK %v(%v; %v; %v)", VERSION, runtime.GOOS, runtime.GOARCH, runtime.Version()))

	resp, err := s.HttpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

// sign generate Zotapay API signature for Deposit and Payout
func (s SDK) sign(args ...string) (signature string) {

	str := ""
	for _, v := range args {
		str += v
	}

	h := sha256.New()
	h.Write([]byte(str + s.MerchantSecretKey))

	signature = hex.EncodeToString(h.Sum(nil))
	return
}
