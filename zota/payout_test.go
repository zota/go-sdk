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
type testP struct {
	name                 string
	expectedPayoutResult PayoutResult
	mockSDK              SDK
	mockPayoutOrder      PayoutOrder
	expectedError        error
}

// PClientMockSuccess mock successful api response
// implement httpClient interface
type PClientMockSuccess struct{}

func (c *PClientMockSuccess) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
    "code": "200",
    "data": {
        "payoutUrl": "https://api.com/api/v1/payout/init/8b3a6b89697e8ac8f45d964bcc90c7ba41764acd/",
        "merchantOrderID": "QvE8dZshpKhaOmHY",
        "orderID": "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd"
    }
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

// PClientMockSuccess mock api error response
// implement httpClient interface
type PClientMockErr struct{}

func (c *PClientMockErr) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
   "code": "400",
   "message": "bad request"
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

// PClientMockErrDo mock client do error response
// implement httpClient interface
type PClientMockErrDo struct{}

func (c *PClientMockErrDo) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
   "code": "400",
   "message": "bad request"
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, fmt.Errorf("do error")
}

// PClientMockUnexpectedResp mock unexpected api response
// implement httpClient interface
type PClientMockUnexpectedResp struct{}

func (c *PClientMockUnexpectedResp) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`Unexpected data`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

func Test_Payout(t *testing.T) {
	// init http mocks
	cMockSuccess := PClientMockSuccess{}
	cMockAPIErr := PClientMockErr{}
	cMockAPIErrDo := PClientMockErrDo{}
	cMockAPIUnexpectedResp := PClientMockUnexpectedResp{}

	// init all test cases
	tests := []testP{
		{
			// --------------------Test Successful Api response----------------------
			name: "Payout Request Success",
			expectedPayoutResult: PayoutResult{
				Code: "200",
				Data: PayoutResultData{
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
			mockPayoutOrder: PayoutOrder{
				MerchantOrderID:                "134e4f443t651",
				MerchantOrderDesc:              "Test order description",
				OrderAmount:                    "500",
				OrderCurrency:                  "MYR", //"USD"
				CustomerEmail:                  "customer@email-address.com",
				CustomerFirstName:              "John",
				CustomerLastName:               "Doe",
				CustomerPhone:                  "+1 420-100-1000",
				CustomerIP:                     "127.0.0.1",
				CustomerBankCode:               "BBL",
				CallbackURL:                    "https://some.endpoint/callback",
				CheckoutURL:                    "https://some.endpoint/checkout",
				CustomerBankAccountNumber:      "100200",
				CustomerBankAccountName:        "John Doe",
				CustomerBankBranch:             "Bank Branch",
				CustomerBankAddress:            "Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
				CustomerBankZipCode:            "84280",
				CustomerBankProvince:           "Bank Province",
				CustomerBankArea:               "Bank Area / City",
				CustomerBankRoutingNumber:      "000",
				CustomerCountryCode:            "MY",
				CustomerPersonalID:             "12345678",
				CustomerBankAccountNumberDigit: "02",
				CustomerBankAccountType:        "03",
				CustomerBankSwiftCode:          "123456789",
				CustomerBankBranchDigit:        "04",
			},
		},
		{
			// --------------------Test Error Api response----------------------
			name: "Payout Request API Error",
			expectedPayoutResult: PayoutResult{
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
			mockPayoutOrder: PayoutOrder{
				MerchantOrderID:                "134e4f443t651",
				MerchantOrderDesc:              "Test order description",
				OrderAmount:                    "500",
				OrderCurrency:                  "MYR", //"USD"
				CustomerEmail:                  "customer@email-address.com",
				CustomerFirstName:              "John",
				CustomerLastName:               "Doe",
				CustomerPhone:                  "+1 420-100-1000",
				CustomerIP:                     "127.0.0.1",
				CustomerBankCode:               "BBL",
				CallbackURL:                    "https://some.endpoint/callback",
				CheckoutURL:                    "https://some.endpoint/checkout",
				CustomerBankAccountNumber:      "100200",
				CustomerBankAccountName:        "John Doe",
				CustomerBankBranch:             "Bank Branch",
				CustomerBankAddress:            "Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
				CustomerBankZipCode:            "84280",
				CustomerBankProvince:           "Bank Province",
				CustomerBankArea:               "Bank Area / City",
				CustomerBankRoutingNumber:      "000",
				CustomerCountryCode:            "MY",
				CustomerPersonalID:             "12345678",
				CustomerBankAccountNumberDigit: "02",
				CustomerBankAccountType:        "03",
				CustomerBankSwiftCode:          "123456789",
				CustomerBankBranchDigit:        "04",
			},
		},
		{
			// --------------------Test Http Do Error Api response----------------------
			name: "Payout Request API Error",
			expectedPayoutResult: PayoutResult{
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
			mockPayoutOrder: PayoutOrder{
				MerchantOrderID:                "134e4f443t651",
				MerchantOrderDesc:              "Test order description",
				OrderAmount:                    "500",
				OrderCurrency:                  "MYR", //"USD"
				CustomerEmail:                  "customer@email-address.com",
				CustomerFirstName:              "John",
				CustomerLastName:               "Doe",
				CustomerPhone:                  "+1 420-100-1000",
				CustomerIP:                     "127.0.0.1",
				CustomerBankCode:               "BBL",
				CallbackURL:                    "https://some.endpoint/callback",
				CheckoutURL:                    "https://some.endpoint/checkout",
				CustomerBankAccountNumber:      "100200",
				CustomerBankAccountName:        "John Doe",
				CustomerBankBranch:             "Bank Branch",
				CustomerBankAddress:            "Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
				CustomerBankZipCode:            "84280",
				CustomerBankProvince:           "Bank Province",
				CustomerBankArea:               "Bank Area / City",
				CustomerBankRoutingNumber:      "000",
				CustomerCountryCode:            "MY",
				CustomerPersonalID:             "12345678",
				CustomerBankAccountNumberDigit: "02",
				CustomerBankAccountType:        "03",
				CustomerBankSwiftCode:          "123456789",
				CustomerBankBranchDigit:        "04",
			},
		},
		{
			// --------------------Test Unmarshal Err----------------------
			name:                 "Payout Request Unexpected Resp",
			expectedPayoutResult: PayoutResult{},
			expectedError:        fmt.Errorf("json Unmarshal err:invalid character 'U' looking for beginning of value"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIUnexpectedResp,
			},
			mockPayoutOrder: PayoutOrder{
				MerchantOrderID:                "134e4f443t651",
				MerchantOrderDesc:              "Test order description",
				OrderAmount:                    "500",
				OrderCurrency:                  "MYR", //"USD"
				CustomerEmail:                  "customer@email-address.com",
				CustomerFirstName:              "John",
				CustomerLastName:               "Doe",
				CustomerPhone:                  "+1 420-100-1000",
				CustomerIP:                     "127.0.0.1",
				CustomerBankCode:               "BBL",
				CallbackURL:                    "https://some.endpoint/callback",
				CheckoutURL:                    "https://some.endpoint/checkout",
				CustomerBankAccountNumber:      "100200",
				CustomerBankAccountName:        "John Doe",
				CustomerBankBranch:             "Bank Branch",
				CustomerBankAddress:            "Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
				CustomerBankZipCode:            "84280",
				CustomerBankProvince:           "Bank Province",
				CustomerBankArea:               "Bank Area / City",
				CustomerBankRoutingNumber:      "000",
				CustomerCountryCode:            "MY",
				CustomerPersonalID:             "12345678",
				CustomerBankAccountNumberDigit: "02",
				CustomerBankAccountType:        "03",
				CustomerBankSwiftCode:          "123456789",
				CustomerBankBranchDigit:        "04",
			},
		}, {
			// --------------------Test Unmarshal Err----------------------
			name:                 "Payout Struct Validate",
			expectedPayoutResult: PayoutResult{},
			expectedError:        fmt.Errorf("MerchantOrderID is required"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIUnexpectedResp,
			},
			mockPayoutOrder: PayoutOrder{},
		},
		{
			// --------------------Test Validations----------------------
			name:                 "Payout Request Validation Error",
			expectedPayoutResult: PayoutResult{},
			expectedError:        fmt.Errorf("MerchantID is required"),
			//empty struct will trigger validation
			mockSDK:         SDK{},
			mockPayoutOrder: PayoutOrder{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			//execute the Payout request
			res, err := test.mockSDK.Payout(test.mockPayoutOrder)

			assert.IsType(t, PayoutResult{}, res)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedPayoutResult, res)
		})
	}
}

// TestPayoutResult_SetMockResponse test mock logic
func TestPayoutResult_SetMockResponse(t *testing.T) {
	cMockSuccess := ClientMockSuccess{}

	//mock the response
	mock := &PayoutResult{
		Code: "200",
		Data: PayoutResultData{
			MerchantOrderID: "123",
			OrderID:         "1234",
		},
		Message: "SomeMockMsg",
	}
	mock.SetMockResponse()

	//init payout
	sdk := &SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        SANDBOX,
		HttpClient:        &cMockSuccess,
	}
	res, err := sdk.Payout(PayoutOrder{
		MerchantOrderID:                "134e4f443t651",
		MerchantOrderDesc:              "Test order description",
		OrderAmount:                    "500",
		OrderCurrency:                  "MYR", //"USD"
		CustomerEmail:                  "customer@email-address.com",
		CustomerFirstName:              "John",
		CustomerLastName:               "Doe",
		CustomerPhone:                  "+1 420-100-1000",
		CustomerIP:                     "127.0.0.1",
		CustomerBankCode:               "BBL",
		CallbackURL:                    "https://some.endpoint/callback",
		CheckoutURL:                    "https://some.endpoint/checkout",
		CustomerBankAccountNumber:      "100200",
		CustomerBankAccountName:        "John Doe",
		CustomerBankBranch:             "Bank Branch",
		CustomerBankAddress:            "Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
		CustomerBankZipCode:            "84280",
		CustomerBankProvince:           "Bank Province",
		CustomerBankArea:               "Bank Area / City",
		CustomerBankRoutingNumber:      "000",
		CustomerCountryCode:            "MY",
		CustomerPersonalID:             "12345678",
		CustomerBankAccountNumberDigit: "02",
		CustomerBankAccountType:        "03",
		CustomerBankSwiftCode:          "123456789",
		CustomerBankBranchDigit:        "04",
	})

	//assert that the response is same as the mocked struct
	assert.Equal(t, mock, &res)
	assert.Equal(t, nil, err)

	//do new new payout
	res, err = sdk.Payout(PayoutOrder{
		MerchantOrderID:                "134e4f443t651",
		MerchantOrderDesc:              "Test order description",
		OrderAmount:                    "500",
		OrderCurrency:                  "MYR", //"USD"
		CustomerEmail:                  "customer@email-address.com",
		CustomerFirstName:              "John",
		CustomerLastName:               "Doe",
		CustomerPhone:                  "+1 420-100-1000",
		CustomerIP:                     "127.0.0.1",
		CustomerBankCode:               "BBL",
		CallbackURL:                    "https://some.endpoint/callback",
		CheckoutURL:                    "https://some.endpoint/checkout",
		CustomerBankAccountNumber:      "100200",
		CustomerBankAccountName:        "John Doe",
		CustomerBankBranch:             "Bank Branch",
		CustomerBankAddress:            "Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
		CustomerBankZipCode:            "84280",
		CustomerBankProvince:           "Bank Province",
		CustomerBankArea:               "Bank Area / City",
		CustomerBankRoutingNumber:      "000",
		CustomerCountryCode:            "MY",
		CustomerPersonalID:             "12345678",
		CustomerBankAccountNumberDigit: "02",
		CustomerBankAccountType:        "03",
		CustomerBankSwiftCode:          "123456789",
		CustomerBankBranchDigit:        "04",
	})

	//assert that the response is NOT same as the mocked struct
	assert.NotEqual(t, mock, &res)
	assert.Equal(t, nil, err)
}
