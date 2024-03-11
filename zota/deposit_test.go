package zota

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test base structure
type testD struct {
	name                  string
	expectedDepositResult DepositResult
	mockSDK               SDK
	mockDepositOrder      DepositOrder
	expectedError         error
}

// ClientMockSuccess mock successful api response
// implement httpClient interface
type ClientMockSuccess struct{}

func (c *ClientMockSuccess) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
    "code": "200",
    "data": {
        "depositUrl": "https://api.com/api/v1/deposit/init/8b3a6b89697e8ac8f45d964bcc90c7ba41764acd/",
        "merchantOrderID": "QvE8dZshpKhaOmHY",
        "orderID": "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd"
    }
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

// ClientMockSuccess mock api error response
// implement httpClient interface
type ClientMockErr struct{}

func (c *ClientMockErr) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
   "code": "400",
   "message": "bad request"
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

// ClientMockErrDo mock client do error response
// implement httpClient interface
type ClientMockErrDo struct{}

func (c *ClientMockErrDo) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
   "code": "400",
   "message": "bad request"
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, fmt.Errorf("do error")
}

// ClientMockUnexpectedResp mock unexpected api response
// implement httpClient interface
type ClientMockUnexpectedResp struct{}

func (c *ClientMockUnexpectedResp) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`Unexpected data`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

func Test_Deposit(t *testing.T) {
	// init http mocks
	cMockSuccess := ClientMockSuccess{}
	cMockAPIErr := ClientMockErr{}
	cMockAPIErrDo := ClientMockErrDo{}
	cMockAPIUnexpectedResp := ClientMockUnexpectedResp{}

	// init all test cases
	tests := []testD{
		{
			// --------------------Test Successful Api response----------------------
			name: "Deposit Request Success",
			expectedDepositResult: DepositResult{
				Code: "200",
				Data: DepositResultData{
					DepositURL:      "https://api.com/api/v1/deposit/init/8b3a6b89697e8ac8f45d964bcc90c7ba41764acd/",
					MerchantOrderID: "QvE8dZshpKhaOmHY",
					OrderID:         "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd",
				},
			},
			expectedError: nil,
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockSuccess,
			},
			mockDepositOrder: DepositOrder{
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
				CheckoutURL:         "https://vt.com/",
				CustomParam:         "",
				Language:            "EN",
			},
		},
		{
			// --------------------Test Http Error Api response----------------------
			name: "Deposit Http Request API Error",
			expectedDepositResult: DepositResult{
				Code:    "400",
				Message: "bad request",
			},
			expectedError: nil,
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIErr,
			},
			mockDepositOrder: DepositOrder{
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
				CheckoutURL:         "https://vt.com/",
				CustomParam:         "",
				Language:            "EN",
			},
		},
		{
			// --------------------Test Http Do Error Api response----------------------
			name: "Deposit Request API Error",
			expectedDepositResult: DepositResult{
				Code:    "",
				Message: "",
			},
			expectedError: fmt.Errorf("do error"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIErrDo,
			},
			mockDepositOrder: DepositOrder{
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
				CheckoutURL:         "https://vt.com/",
				CustomParam:         "",
				Language:            "EN",
			},
		},
		{
			// --------------------Test Unmarshal Err----------------------
			name:                  "Deposit Request Unexpected Resp",
			expectedDepositResult: DepositResult{},
			expectedError:         fmt.Errorf("json Unmarshal err:invalid character 'U' looking for beginning of value"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIUnexpectedResp,
			},
			mockDepositOrder: DepositOrder{
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
				CheckoutURL:         "https://vt.com/",
				CustomParam:         "",
				Language:            "EN",
			},
		}, {
			// --------------------Test Validate Deposit----------------------
			name:                  "Deposit Struct Validate",
			expectedDepositResult: DepositResult{},
			expectedError:         fmt.Errorf("MerchantOrderID is required"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIUnexpectedResp,
			},
			mockDepositOrder: DepositOrder{},
		},
		{
			// --------------------Test Validations----------------------
			name:                  "Deposit Request Validation Error",
			expectedDepositResult: DepositResult{},
			expectedError:         fmt.Errorf("MerchantID is required."),
			//empty struct will trigger validation
			mockSDK:          SDK{},
			mockDepositOrder: DepositOrder{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			//execute the Deposit request
			res, err := test.mockSDK.Deposit(test.mockDepositOrder)

			assert.IsType(t, DepositResult{}, res)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedDepositResult, res)
		})
	}
}

// TestDepositResult_SetMockResponse test mock logic
func TestDepositResult_SetMockResponse(t *testing.T) {
	cMockSuccess := ClientMockSuccess{}

	//mock the response
	mock := &DepositResult{
		Code: "200",
		Data: DepositResultData{
			DepositURL:      "http://some.mock",
			MerchantOrderID: "123",
			OrderID:         "1234",
		},
		Message: "SomeMockMsg",
	}
	mock.SetMockResponse()

	//init deposit
	sdk := &SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        SANDBOX,
		HttpClient:        &cMockSuccess,
	}
	res, err := sdk.Deposit(DepositOrder{
		MerchantOrderID:     "134e4f44t651",
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
		Language:            "EN",
	})

	//assert that the response is same as the mocked struct
	assert.Equal(t, mock, &res)
	assert.Equal(t, nil, err)

	//do new new deposit
	res, err = sdk.Deposit(DepositOrder{
		MerchantOrderID:     "134e4f44t651",
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
		Language:            "EN",
	})

	//assert that the response is NOT same as the mocked struct
	assert.NotEqual(t, mock, &res)
	assert.Equal(t, nil, err)
}

func TestDeposit_Validate(t *testing.T) {

	tests := []test{
		{
			name: "success",
			mock: &DepositOrder{
				MerchantOrderID:     "134e4f44t651",
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
				Language:            "EN",
			},
			expectedError: nil,
		}, {
			name:          "missing field",
			mock:          &DepositOrder{},
			expectedError: fmt.Errorf("MerchantOrderID is required"),
		},
	}

	for _, test := range tests {
		err := test.mock.(*DepositOrder).validate()
		assert.Equal(t, test.expectedError, err)
	}
}
