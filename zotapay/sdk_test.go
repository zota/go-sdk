package zotapay

import (
	"crypto/tls"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

// test base structure
type test struct {
	name          string
	mock          interface{}
	mockOrder     interface{}
	expected      interface{}
	expectedError error
}

// Test_Validate test the SDK validation
func Test_Validate(t *testing.T) {

	tests := []test{
		{
			name: "Success",
			mock: &SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
			},
			expectedError: nil,
		}, {
			name: "Unexpected ApiBaseURL",
			mock: &SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        "http://some.wrong",
			},
			expectedError: fmt.Errorf("unexpected ApiBaseURL."),
		}, {
			name: "Missing EndpointID",
			mock: &SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				ApiBaseURL:        SANDBOX,
			},
			expectedError: fmt.Errorf("EndpointID is required."),
		}, {
			name: "Missing EndpointID",
			mock: &SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
			},
			expectedError: fmt.Errorf("ApiBaseURL is required."),
		}, {
			name: "Missing MerchantSecretKey",
			mock: &SDK{
				MerchantID: "API_MERCHANT_ID",
				EndpointID: "503368",
				ApiBaseURL: SANDBOX,
			},
			expectedError: fmt.Errorf("MerchantSecretKey is required."),
		}, {
			name: "Missing MerchantID",
			mock: &SDK{
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
			},
			expectedError: fmt.Errorf("MerchantID is required."),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			//validate the
			err := test.mock.(*SDK).validate()
			assert.Equal(t, test.expectedError, err)
		})
	}
}

// Test_Validate test the SDK validation
func Test_InitHttpClient(t *testing.T) {

	tests := []test{
		{
			name: "InitIfEmpty",
			mock: &SDK{
				HttpClient: nil,
			},
			expected: &http.Client{
				Timeout: time.Second * 10,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{PreferServerCipherSuites: true, MinVersion: tls.VersionTLS12},
				},
			},
		}, {
			name: "NotInitIfExists",
			mock: &SDK{
				HttpClient: &http.Client{
					Timeout: time.Second * 20,
				},
			},
			expected: &http.Client{
				Timeout: time.Second * 20,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			//validate the
			test.mock.(*SDK).initHttpClient()
			assert.Equal(t, test.expected, test.mock.(*SDK).HttpClient)
		})
	}
}

func Test_Sign(t *testing.T) {

	tests := []test{
		{
			name: "Success",
			mock: &SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
			},
			mockOrder: &DepositOrder{
				MerchantOrderID:     "134",
				MerchantOrderDesc:   "Test order description",
				OrderAmount:         "500",
				OrderCurrency:       "MYR",
				CustomerEmail:       "customer@email-address.com",
				CustomerFirstName:   "John",
				CustomerLastName:    "Doe",
				CustomerAddress:     "The Swan, Jungle St. 108",
				CustomerCountryCode: "US",
				CustomerCity:        "Los Angeles",
				CustomerState:       "CA",
				CustomerZipCode:     "84280",
				CustomerPhone:       "+1 420-100-1000",
				CustomerIP:          "127.0.0.1",
				CustomerBankCode:    "BBL",
				RedirectURL:         "https://some.endpoint/redirect",
				CallbackURL:         "https://some.endpoint/callback",
				CheckoutURL:         "https://some.endpoint/checkout",
				CustomParam:         "",
				Language:            "EN",
			},
			expected:      "ab46525abe045d8a8bf6d6f6f3139d39814ff480e8e55280ae852307ab83e66a",
			expectedError: nil,
		},
	}

	for _, test := range tests {
		sign := test.mock.(*SDK).sign(test.mockOrder.(*DepositOrder).MerchantOrderID, test.mockOrder.(*DepositOrder).OrderAmount, test.mockOrder.(*DepositOrder).CustomerEmail)
		assert.Equal(t, test.expected, sign)
	}
}
