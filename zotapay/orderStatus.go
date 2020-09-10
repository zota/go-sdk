package zotapay

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// OrderStatus represents check order status
type OrderStatus struct {
	MerchantID      string `url:"merchantID"`
	MerchantOrderID string `url:"merchantOrderID"`
	OrderID         string `url:"orderID"`
	Timestamp       string `url:"timestamp"`
	Signature       string `url:"signature"`
}

// OrderStatusResult represents the response from Zotapay API
type OrderStatusResult struct {
	Code                  string `json:"code"`
	OrderStatusResultData `json:"data"`
	Message               string `json:"message"`
}

type OrderStatusResultData struct {
	Type                  string `json:"type"`
	Status                string `json:"status"`
	ErrorMessage          string `json:"errorMessage"`
	EndpointID            string `json:"endpointID"`
	ExternalTransactionID string `json:"externalTransactionID"`
	OrderID               string `json:"orderID"`
	MerchantOrderID       string `json:"merchantOrderID"`
	Amount                string `json:"amount"`
	Currency              string `json:"currency"`
	CustomerEmail         string `json:"customerEmail"`
	CustomParam           string `json:"customParam"`
	ExtraData             struct {
		Dcc              bool   `json:"dcc"`
		SelectedBankCode string `json:"selectedBankCode"`
		SelectedBankName string `json:"selectedBankName"`
	} `json:"extraData"`
	Request struct {
		MerchantID      string `json:"merchantID"`
		OrderID         string `json:"orderID"`
		MerchantOrderID string `json:"merchantOrderID"`
		Timestamp       string `json:"timestamp"`
	} `json:"request"`
}

var mockedOrderStatusResult *OrderStatusResult

// OrderStatus init validation of the SDK struct and the OrderStatus
// generate sign and
// init order status request to Zotapay API
func (s *SDK) OrderStatus(d OrderStatus) (res OrderStatusResult, err error) {

	//validate that SDK is properly initialized
	err = s.validate()
	if err != nil {
		return
	}

	//validate that OrderStatus is properly initialized
	err = d.validate()
	if err != nil {
		return
	}

	d.MerchantID = s.MerchantID

	//set current timestamp
	d.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)

	//if mockedOrderStatusResult is set return it as response
	//only for testing
	if mockedOrderStatusResult != nil {
		res = *mockedOrderStatusResult
		mockedOrderStatusResult = nil
		return
	}

	//generate signature
	d.Signature = s.sign(d.MerchantID, d.MerchantOrderID, d.OrderID, d.Timestamp)

	v, err := query.Values(d)
	if err != nil {
		return
	}

	_, body, err := s.httpDo(http.MethodGet, fmt.Sprintf("%v/api/v1/query/order-status/?%v", s.ApiBaseURL, v.Encode()), []byte(""))
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("json Unmarshal err:%v", err)
		return
	}

	return
}

// SetMockResponse will set *OrderStatusResult as mocked
// response on the next order status request
// To be used only for test purposes
func (mock *OrderStatusResult) SetMockResponse() {
	mockedOrderStatusResult = mock
}

// validate the instance of OrderStatus
// if not valid returns an error
func (d *OrderStatus) validate() error {
	required := []string{"MerchantOrderID", "OrderID"}
	for _, fieldName := range required {
		r := reflect.ValueOf(d)
		f := reflect.Indirect(r).FieldByName(fieldName)
		value := f.String()
		if value == "" {
			return fmt.Errorf("%v is required", fieldName)
		}
	}
	return nil
}
