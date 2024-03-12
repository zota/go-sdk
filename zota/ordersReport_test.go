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
type testR struct {
	name                       string
	expectedOrdersReportResult OrdersReportResult
	mockSDK                    SDK
	mockOrdersReport           OrdersReport
	expectedError              error
}

// ClientMockOrdersRepSuccess mock successful api response
// implement httpClient interface
type ClientMockOrdersRepSuccess struct{}

func (c *ClientMockOrdersRepSuccess) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`id,parent_id,order_type,status,merchant_id,batch_id,endpoint_id,endpoint_group_id,order_currency,order_amount,original_amount,amount_changed,merchant_order_id,payment_method_id,external_transaction_id,client_error_message,status_changed,created_at,ended_at,ended_with_status,last_update_at,is_refunded,is_fully_refunded,refunded_amount,refunded_at,request_customer_ip,entered_deposit_url,entered_bank_selection_page,bank_selected,selected_bank_code,selected_bank_name,callback_sent_to_merchant,callback_received_by_merchant,customer_email,customer_first_name,customer_last_name,customer_address,customer_country_code,customer_city,customer_state,customer_zip_code,customer_phone,customer_bank_code,customer_bank_account_number,customer_bank_account_name,customer_bank_branch,customer_bank_address,customer_bank_zip_code,customer_bank_routing_number,customer_bank_province,customer_bank_area
24050211,,PAYOUT,DECLINED,API_MERCHANT_ID,,503368,,MYR,500.00000000,500.00000000,false,TbbQzewLWwDW6goc,PAYOUT,,insufficient funds,false,2020-08-05 13:03:30 +0000 UTC,2020-08-05 13:03:30 +0000 UTC,DECLINED,2020-08-05 13:03:30 +0000 UTC,false,false,0.00000000,0001-01-01 00:00:00 +0000 UTC,103.106.8.104,false,false,false,,Bangkok Bank,true,false,customer@email-address.com,John,,,,,,,66-77999110,BBL,100200,John Doe,Bank Branch,"Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",84280,000,Bank Province,Bank Area / City`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

// ClientMockOrdersRepSuccess mock api error response
// implement httpClient interface
type ClientMockOrdersRepErr struct{}

func (c *ClientMockOrdersRepErr) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
   "code": "400",
   "message": "bad request"
}`)))
	return &http.Response{
		StatusCode: 400,
		Body:       r,
	}, nil
}

// ClientMockOrdersRepErrDo mock client do error response
// implement httpClient interface
type ClientMockOrdersRepErrDo struct{}

func (c *ClientMockOrdersRepErrDo) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`{
   "code": "400",
   "message": "bad request"
}`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, fmt.Errorf("do error")
}

// ClientMockOrdersRepUnexpectedResp mock unexpected api response
// implement httpClient interface
type ClientMockOrdersRepUnexpectedResp struct{}

func (c *ClientMockOrdersRepUnexpectedResp) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(`Unexpected data`)))
	return &http.Response{
		StatusCode: 400,
		Body:       r,
	}, nil
}

func Test_OrdersReport(t *testing.T) {
	// init http mocks
	cMockSuccess := ClientMockOrdersRepSuccess{}
	cMockAPIErr := ClientMockOrdersRepErr{}
	cMockAPIErrDo := ClientMockOrdersRepErrDo{}
	cMockAPIUnexpectedResp := ClientMockOrdersRepUnexpectedResp{}

	// init all test cases
	tests := []testR{
		{
			// --------------------Test Successful Api response----------------------
			name: "OrdersReport Request Success",
			expectedOrdersReportResult: OrdersReportResult{
				Code: "200",
				OrdersReport: `id,parent_id,order_type,status,merchant_id,batch_id,endpoint_id,endpoint_group_id,order_currency,order_amount,original_amount,amount_changed,merchant_order_id,payment_method_id,external_transaction_id,client_error_message,status_changed,created_at,ended_at,ended_with_status,last_update_at,is_refunded,is_fully_refunded,refunded_amount,refunded_at,request_customer_ip,entered_deposit_url,entered_bank_selection_page,bank_selected,selected_bank_code,selected_bank_name,callback_sent_to_merchant,callback_received_by_merchant,customer_email,customer_first_name,customer_last_name,customer_address,customer_country_code,customer_city,customer_state,customer_zip_code,customer_phone,customer_bank_code,customer_bank_account_number,customer_bank_account_name,customer_bank_branch,customer_bank_address,customer_bank_zip_code,customer_bank_routing_number,customer_bank_province,customer_bank_area
24050211,,PAYOUT,DECLINED,API_MERCHANT_ID,,503368,,MYR,500.00000000,500.00000000,false,TbbQzewLWwDW6goc,PAYOUT,,insufficient funds,false,2020-08-05 13:03:30 +0000 UTC,2020-08-05 13:03:30 +0000 UTC,DECLINED,2020-08-05 13:03:30 +0000 UTC,false,false,0.00000000,0001-01-01 00:00:00 +0000 UTC,103.106.8.104,false,false,false,,Bangkok Bank,true,false,customer@email-address.com,John,,,,,,,66-77999110,BBL,100200,John Doe,Bank Branch,"Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",84280,000,Bank Province,Bank Area / City`,
			},
			expectedError: nil,
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockSuccess,
			},
			mockOrdersReport: OrdersReport{
				DateType:    "created",
				EndpointIds: "503368,503365",
				FromDate:    "2020-08-01",
				ToDate:      "2020-09-01",
				Statuses:    "APPROVED,DECLINED",
				Types:       "SALE,PAYOUT",
			},
		},
		{
			// --------------------Test Http Error Api response----------------------
			name: "OrdersReport Http Request API Error",
			expectedOrdersReportResult: OrdersReportResult{
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
			mockOrdersReport: OrdersReport{
				DateType:    "created",
				EndpointIds: "503368,503365",
				FromDate:    "2020-08-01",
				ToDate:      "2020-09-01",
				Statuses:    "APPROVED,DECLINED",
				Types:       "SALE,PAYOUT",
			},
		},
		{
			// --------------------Test Http Do Error Api response----------------------
			name: "OrdersReport Request API Error",
			expectedOrdersReportResult: OrdersReportResult{
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
			mockOrdersReport: OrdersReport{
				DateType:    "created",
				EndpointIds: "503368,503365",
				FromDate:    "2020-08-01",
				ToDate:      "2020-09-01",
				Statuses:    "APPROVED,DECLINED",
				Types:       "SALE,PAYOUT",
			},
		},
		{
			// --------------------Test Unmarshal Err----------------------
			name:                       "OrdersReport Request Unexpected Resp",
			expectedOrdersReportResult: OrdersReportResult{},
			expectedError:              fmt.Errorf("json Unmarshal err:invalid character 'U' looking for beginning of value"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIUnexpectedResp,
			},
			mockOrdersReport: OrdersReport{
				DateType:    "created",
				EndpointIds: "503368,503365",
				FromDate:    "2020-08-01",
				ToDate:      "2020-09-01",
				Statuses:    "APPROVED,DECLINED",
				Types:       "SALE,PAYOUT",
			},
		}, {
			// --------------------Test Validate OrdersReport----------------------
			name:                       "OrdersReport Struct Validate",
			expectedOrdersReportResult: OrdersReportResult{},
			expectedError:              fmt.Errorf("DateType is required"),
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
				HttpClient:        &cMockAPIUnexpectedResp,
			},
			mockOrdersReport: OrdersReport{},
		},
		{
			// --------------------Test Validations----------------------
			name:                       "OrdersReport Request Validation Error",
			expectedOrdersReportResult: OrdersReportResult{},
			expectedError:              fmt.Errorf("MerchantID is required"),
			//empty struct will trigger validation
			mockSDK:          SDK{},
			mockOrdersReport: OrdersReport{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			//execute the OrdersReport request
			res, err := test.mockSDK.OrdersReport(test.mockOrdersReport)

			assert.IsType(t, OrdersReportResult{}, res)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedOrdersReportResult, res)
		})
	}
}

// TestOrdersReportResult_SetMockResponse test mock logic
func TestOrdersReportResult_SetMockResponse(t *testing.T) {
	cMockSuccess := ClientMockOrdersRepSuccess{}

	//mock the response
	mock := &OrdersReportResult{
		Code: "200",
		OrdersReport: `id,parent_id,order_type,status,merchant_id,batch_id,endpoint_id,endpoint_group_id,order_currency,order_amount,original_amount,amount_changed,merchant_order_id,payment_method_id,external_transaction_id,client_error_message,status_changed,created_at,ended_at,ended_with_status,last_update_at,is_refunded,is_fully_refunded,refunded_amount,refunded_at,request_customer_ip,entered_deposit_url,entered_bank_selection_page,bank_selected,selected_bank_code,selected_bank_name,callback_sent_to_merchant,callback_received_by_merchant,customer_email,customer_first_name,customer_last_name,customer_address,customer_country_code,customer_city,customer_state,customer_zip_code,customer_phone,customer_bank_code,customer_bank_account_number,customer_bank_account_name,customer_bank_branch,customer_bank_address,customer_bank_zip_code,customer_bank_routing_number,customer_bank_province,customer_bank_area
24050211,,PAYOUT,DECLINED,API_MERCHANT_ID,,503368,,MYR,500.00000000,500.00000000,false,TbbQzewLWwDW6goc,PAYOUT,,insufficient funds,false,2020-08-05 13:03:30 +0000 UTC,2020-08-05 13:03:30 +0000 UTC,DECLINED,2020-08-05 13:03:30 +0000 UTC,false,false,0.00000000,0001-01-01 00:00:00 +0000 UTC,103.106.8.104,false,false,false,,Bangkok Bank,true,false,customer@email-address.com,John,,,,,,,66-77999110,BBL,100200,John Doe,Bank Branch,"Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",84280,000,Bank Province,Bank Area / City`,
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
	res, err := sdk.OrdersReport(OrdersReport{
		DateType:    "created",
		EndpointIds: "503368,503365",
		FromDate:    "2020-08-01",
		ToDate:      "2020-09-01",
		Statuses:    "APPROVED,DECLINED",
		Types:       "SALE,PAYOUT",
	})

	//assert that the response is same as the mocked struct
	assert.Equal(t, mock, &res)
	assert.Equal(t, nil, err)

	//do new new deposit
	res, err = sdk.OrdersReport(OrdersReport{
		DateType:    "created",
		EndpointIds: "503368,503365",
		FromDate:    "2020-08-01",
		ToDate:      "2020-09-01",
		Statuses:    "APPROVED,DECLINED",
		Types:       "SALE,PAYOUT",
	})

	//assert that the response is NOT same as the mocked struct
	assert.NotEqual(t, mock, &res)
	assert.Equal(t, nil, err)
}

func Test_OrdersReport_Validate(t *testing.T) {

	tests := []test{
		{
			name: "success",
			mock: &OrdersReport{
				DateType:    "created",
				EndpointIds: "503368,503365",
				FromDate:    "2020-08-01",
				ToDate:      "2020-09-01",
				Statuses:    "APPROVED,DECLINED",
				Types:       "SALE,PAYOUT",
			},
			expectedError: nil,
		}, {
			name:          "missing field",
			mock:          &OrdersReport{},
			expectedError: fmt.Errorf("DateType is required"),
		},
	}

	for _, test := range tests {
		err := test.mock.(*OrdersReport).validate()
		assert.Equal(t, test.expectedError, err)
	}
}
