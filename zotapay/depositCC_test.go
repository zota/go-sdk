package zotapay

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

// test base structure
type testCC struct {
	name                    string
	expectedDepositCCResult DepositCCResult
	mockSDK                 SDK
	mockDepositCCOrder      DepositCCOrder
	expectedError           error
}

// ClientMockSuccessCC mock successful api response
// implement httpClient interface
type ClientMockSuccessCC struct{}

func (c *ClientMockSuccessCC) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
    "code": "200",
    "data": {
        "status": "PROCESSING",
        "merchantOrderID": "QvE8dZshpKhaOmHY",
        "orderID": "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd"
    }
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

// ClientMockSuccessCC mock api error response
// implement httpClient interface
type ClientMockErrCC struct{}

func (c *ClientMockErrCC) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
   "code": "400",
   "message": "bad request"
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

// ClientMockErrCCDo mock client do error response
// implement httpClient interface
type ClientMockErrCCDo struct{}

func (c *ClientMockErrCCDo) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
   "code": "400",
   "message": "bad request"
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, fmt.Errorf("do error")
}

// ClientMockUnexpectedRespCC mock unexpected api response
// implement httpClient interface
type ClientMockUnexpectedRespCC struct{}

func (c *ClientMockUnexpectedRespCC) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`Unexpected data`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

func Test_DepositCC(t *testing.T) {
	// init http mocks
	cMockSuccess := ClientMockSuccessCC{}
	cMockAPIErr := ClientMockErrCC{}
	cMockAPIErrDo := ClientMockErrCCDo{}
	cMockAPIUnexpectedResp := ClientMockUnexpectedRespCC{}

	// init all test cases
	tests := []testCC{
		{
			// --------------------Test Successful Api response----------------------
			name: "Deposit CC Request Success",
			expectedDepositCCResult: DepositCCResult{
				Code: "200",
				Data: DepositCCResultData{
					Status:          "PROCESSING",
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
			mockDepositCCOrder: DepositCCOrder{
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
				// credit card data
				CardHolderName:      "TEST TEST",
				CardNumber:          "4222222222347466",
				CardExpirationMonth: "01",
				CardExpirationYear:  "21",
				CardCvv:             "111",
			},
		},
		{
			// --------------------Test Http Error Api response----------------------
			name: "Deposit Http Request API Error",
			expectedDepositCCResult: DepositCCResult{
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
			mockDepositCCOrder: DepositCCOrder{
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
				// credit card data
				CardHolderName:      "TEST TEST",
				CardNumber:          "4222222222347466",
				CardExpirationMonth: "01",
				CardExpirationYear:  "21",
				CardCvv:             "111",
			},
		},
		{
			// --------------------Test Http Do Error Api response----------------------
			name: "Deposit CC Request API Error",
			expectedDepositCCResult: DepositCCResult{
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
			mockDepositCCOrder: DepositCCOrder{
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
				// credit card data
				CardHolderName:      "TEST TEST",
				CardNumber:          "4222222222347466",
				CardExpirationMonth: "01",
				CardExpirationYear:  "21",
				CardCvv:             "111",
			},
		},
		{
			// --------------------Test Unmarshal Err----------------------
			name:                    "Deposit CC Request Unexpected Resp",
			expectedDepositCCResult: DepositCCResult{},
			expectedError:           fmt.Errorf("json Unmarshal err:invalid character 'U' looking for beginning of value"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIUnexpectedResp,
			},
			mockDepositCCOrder: DepositCCOrder{
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
				// credit card data
				CardHolderName:      "TEST TEST",
				CardNumber:          "4222222222347466",
				CardExpirationMonth: "01",
				CardExpirationYear:  "21",
				CardCvv:             "111",
			},
		}, {
			// --------------------Test Validate Deposit----------------------
			name:                    "Deposit Struct Validate",
			expectedDepositCCResult: DepositCCResult{},
			expectedError:           fmt.Errorf("MerchantOrderID is required"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIUnexpectedResp,
			},
			mockDepositCCOrder: DepositCCOrder{},
		},
		{
			// --------------------Test Validations----------------------
			name:                    "Deposit CC Request Validation Error",
			expectedDepositCCResult: DepositCCResult{},
			expectedError:           fmt.Errorf("MerchantID is required."),
			//empty struct will trigger validation
			mockSDK:            SDK{},
			mockDepositCCOrder: DepositCCOrder{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			//execute the Deposit request
			res, err := test.mockSDK.DepositCC(test.mockDepositCCOrder)

			assert.IsType(t, DepositCCResult{}, res)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedDepositCCResult, res)
		})
	}
}

// TestDepositCCResult_SetMockResponse test mock logic
func TestDepositCCResult_SetMockResponseCC(t *testing.T) {
	cMockSuccess := ClientMockSuccessCC{}

	//mock the response
	mock := &DepositCCResult{
		Code: "200",
		Data: DepositCCResultData{
			Status:          "PROCESSING",
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
	res, err := sdk.DepositCC(DepositCCOrder{
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
		// credit card data
		CardHolderName:      "TEST TEST",
		CardNumber:          "4222222222347466",
		CardExpirationMonth: "01",
		CardExpirationYear:  "21",
		CardCvv:             "111",
	})

	//assert that the response is same as the mocked struct
	assert.Equal(t, mock, &res)
	assert.Equal(t, nil, err)

	//do new new deposit
	res, err = sdk.DepositCC(DepositCCOrder{
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
		// credit card data
		CardHolderName:      "TEST TEST",
		CardNumber:          "4222222222347466",
		CardExpirationMonth: "01",
		CardExpirationYear:  "21",
		CardCvv:             "111",
	})

	//assert that the response is NOT same as the mocked struct
	assert.NotEqual(t, mock, &res)
	assert.Equal(t, nil, err)
}

func TestDeposit_ValidateCC(t *testing.T) {

	tests := []test{
		{
			name: "success",
			mock: &DepositCCOrder{
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
				// credit card data
				CardHolderName:      "TEST TEST",
				CardNumber:          "4222222222347466",
				CardExpirationMonth: "01",
				CardExpirationYear:  "21",
				CardCvv:             "111",
			},
			expectedError: nil,
		}, {
			name:          "missing field",
			mock:          &DepositCCOrder{},
			expectedError: fmt.Errorf("MerchantOrderID is required"),
		},
	}

	for _, test := range tests {
		err := test.mock.(*DepositCCOrder).validate()
		assert.Equal(t, test.expectedError, err)
	}
}
