package zota

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

// test base structure
type testO struct {
	name                      string
	expectedOrderStatusResult OrderStatusResult
	mockSDK                   SDK
	mockOrderStatus           OrderStatus
	expectedError             error
}

// ClientMockOrderStatusSuccess mock successful api response
// implement httpClient interface
type ClientMockOrderStatusSuccess struct{}

func (c *ClientMockOrderStatusSuccess) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{"code":"200","data":{"type":"SALE","status":"CREATED","errorMessage":"","endpointID":"503368","externalTransactionID":"","orderID":"24057644","merchantOrderID":"134e4f44t65111","amount":"500.00","currency":"MYR","customerEmail":"customer@email-address.com","customParam":"","extraData":{"dcc":false,"selectedBankCode":"BBL","selectedBankName":""},"request":{"merchantID":"API_MERCHANT_ID","orderID":"24057644","merchantOrderID":"134e4f44t65111","timestamp":"1598872227"}}}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

// ClientMockOrderStatusSuccess mock api error response
// implement httpClient interface
type ClientMockOrderStatusErr struct{}

func (c *ClientMockOrderStatusErr) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
   "code": "400",
   "message": "bad request"
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

// ClientMockOrderStatusErrDo mock client do error response
// implement httpClient interface
type ClientMockOrderStatusErrDo struct{}

func (c *ClientMockOrderStatusErrDo) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
   "code": "400",
   "message": "bad request"
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, fmt.Errorf("do error")
}

// ClientMockOrderStatusUnexpectedResp mock unexpected api response
// implement httpClient interface
type ClientMockOrderStatusUnexpectedResp struct{}

func (c *ClientMockOrderStatusUnexpectedResp) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`Unexpected data`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

func Test_OrderStatus(t *testing.T) {
	// init http mocks
	cMockSuccess := ClientMockOrderStatusSuccess{}
	cMockAPIErr := ClientMockOrderStatusErr{}
	cMockAPIErrDo := ClientMockOrderStatusErrDo{}
	cMockAPIUnexpectedResp := ClientMockOrderStatusUnexpectedResp{}

	// init all test cases
	tests := []testO{
		{
			// --------------------Test Successful Api response----------------------
			name: "OrderStatus Request Success",
			expectedOrderStatusResult: OrderStatusResult{
				Code: "200",
				OrderStatusResultData: OrderStatusResultData{Type: "SALE", Status: "CREATED", ErrorMessage: "", EndpointID: "503368", ExternalTransactionID: "", OrderID: "24057644", MerchantOrderID: "134e4f44t65111", Amount: "500.00", Currency: "MYR", CustomerEmail: "customer@email-address.com", CustomParam: "", ExtraData: struct {
					Dcc              bool   "json:\"dcc\""
					SelectedBankCode string "json:\"selectedBankCode\""
					SelectedBankName string "json:\"selectedBankName\""
				}{Dcc: false, SelectedBankCode: "BBL", SelectedBankName: ""}, Request: struct {
					MerchantID      string "json:\"merchantID\""
					OrderID         string "json:\"orderID\""
					MerchantOrderID string "json:\"merchantOrderID\""
					Timestamp       string "json:\"timestamp\""
				}{MerchantID: "API_MERCHANT_ID", OrderID: "24057644", MerchantOrderID: "134e4f44t65111", Timestamp: "1598872227"}},
			},
			expectedError: nil,
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockSuccess,
			},
			mockOrderStatus: OrderStatus{
				MerchantOrderID: "134",
				OrderID:         "135",
			},
		},
		{
			// --------------------Test Http Error Api response----------------------
			name: "OrderStatus Http Request API Error",
			expectedOrderStatusResult: OrderStatusResult{
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
			mockOrderStatus: OrderStatus{
				MerchantOrderID: "134",
				OrderID:         "135",
			},
		},
		{
			// --------------------Test Http Do Error Api response----------------------
			name: "OrderStatus Request API Error",
			expectedOrderStatusResult: OrderStatusResult{
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
			mockOrderStatus: OrderStatus{
				MerchantOrderID: "134",
				OrderID:         "135",
			},
		},
		{
			// --------------------Test Unmarshal Err----------------------
			name:                      "OrderStatus Request Unexpected Resp",
			expectedOrderStatusResult: OrderStatusResult{},
			expectedError:             fmt.Errorf("json Unmarshal err:invalid character 'U' looking for beginning of value"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIUnexpectedResp,
			},
			mockOrderStatus: OrderStatus{
				MerchantOrderID: "134",
				OrderID:         "135",
			},
		}, {
			// --------------------Test Validate OrderStatus----------------------
			name:                      "OrderStatus Struct Validate",
			expectedOrderStatusResult: OrderStatusResult{},
			expectedError:             fmt.Errorf("MerchantOrderID is required"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIUnexpectedResp,
			},
			mockOrderStatus: OrderStatus{},
		},
		{
			// --------------------Test Validations----------------------
			name:                      "OrderStatus Request Validation Error",
			expectedOrderStatusResult: OrderStatusResult{},
			expectedError:             fmt.Errorf("MerchantID is required."),
			//empty struct will trigger validation
			mockSDK:         SDK{},
			mockOrderStatus: OrderStatus{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			//execute the OrderStatus request
			res, err := test.mockSDK.OrderStatus(test.mockOrderStatus)

			assert.IsType(t, OrderStatusResult{}, res)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedOrderStatusResult, res)
		})
	}
}

// TestOrderStatusResult_SetMockResponse test mock logic
func TestOrderStatusResult_SetMockResponse(t *testing.T) {
	cMockSuccess := ClientMockOrderStatusSuccess{}

	//mock the response
	mock := &OrderStatusResult{
		Code: "200",
		OrderStatusResultData: OrderStatusResultData{
			Status:          "http://some.mock",
			MerchantOrderID: "123",
			OrderID:         "1234",
		},
		Message: "SomeMockMsg",
	}
	mock.SetMockResponse()

	//init order status check
	sdk := &SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        SANDBOX,
		HttpClient:        &cMockSuccess,
	}
	res, err := sdk.OrderStatus(OrderStatus{
		MerchantOrderID: "134e4f44t651",
		OrderID:         "135",
	})

	//assert that the response is same as the mocked struct
	assert.Equal(t, mock, &res)
	assert.Equal(t, nil, err)

	//do new new deposit
	res, err = sdk.OrderStatus(OrderStatus{
		MerchantOrderID: "134e4f44t651",
		OrderID:         "135",
	})

	//assert that the response is NOT same as the mocked struct
	assert.NotEqual(t, mock, &res)
	assert.Equal(t, nil, err)
}

func Test_OrderStatus_Validate(t *testing.T) {

	tests := []test{
		{
			name: "success",
			mock: &OrderStatus{
				MerchantOrderID: "134e4f44t651",
				OrderID:         "135",
			},
			expectedError: nil,
		}, {
			name:          "missing field",
			mock:          &OrderStatus{},
			expectedError: fmt.Errorf("MerchantOrderID is required"),
		},
	}

	for _, test := range tests {
		err := test.mock.(*OrderStatus).validate()
		assert.Equal(t, test.expectedError, err)
	}
}
