package zotapay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// DepositOrder represents deposit order
type DepositOrder struct {
	MerchantOrderID     string `json:"merchantOrderID"`
	MerchantOrderDesc   string `json:"merchantOrderDesc"`
	OrderAmount         string `json:"orderAmount"`
	OrderCurrency       string `json:"orderCurrency"`
	CustomerEmail       string `json:"customerEmail"`
	CustomerFirstName   string `json:"customerFirstName"`
	CustomerLastName    string `json:"customerLastName"`
	CustomerAddress     string `json:"customerAddress"`
	CustomerCountryCode string `json:"customerCountryCode"`
	CustomerCity        string `json:"customerCity"`
	CustomerState       string `json:"customerState"`
	CustomerZipCode     string `json:"customerZipCode"`
	CustomerPhone       string `json:"customerPhone"`
	CustomerIP          string `json:"customerIP"`
	CustomerBankCode    string `json:"customerBankCode"`
	RedirectURL         string `json:"redirectUrl"`
	CallbackURL         string `json:"callbackUrl"`
	CustomParam         string `json:"customParam"`
	CheckoutURL         string `json:"checkoutUrl"`
	Language            string `json:"language"`
	Signature           string `json:"signature"`
}

// DepositResult represents the response from Zotapay API
type DepositResult struct {
	Code    string            `json:"code"`
	Data    DepositResultData `json:"data"`
	Message string            `json:"message"`
}

type DepositResultData struct {
	DepositURL      string `json:"depositUrl"`
	MerchantOrderID string `json:"merchantOrderID"`
	OrderID         string `json:"orderID"`
}

var mockedDepositResult *DepositResult

// Deposit init validation of the SDK struct and the DepositOrder
// generate sign and
// init a deposit request to Zotapay API
func (s *SDK) Deposit(d DepositOrder) (res DepositResult, err error) {

	//validate that SDK is properly initialized
	err = s.validate()
	if err != nil {
		return
	}

	//validate that DepositOrder is properly initialized
	err = d.validate()
	if err != nil {
		return
	}

	//if mockedDepositResult is set return it as response
	//only for testing
	if mockedDepositResult != nil {
		res = *mockedDepositResult
		mockedDepositResult = nil
		return
	}

	//generate signature
	d.Signature = s.sign(s.EndpointID, d.MerchantOrderID, d.OrderAmount, d.CustomerEmail)

	deposit, err := json.Marshal(d)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%v/api/v1/deposit/request/%v/", s.ApiBaseURL, s.EndpointID)
	fmt.Printf("before deposit req to:%v", url)
	code, body, err := s.httpDo(http.MethodPost, url, deposit)
	fmt.Printf("after deposit req to:%v, code:%v, body:%v, err:%v", url, code, body, err)
	if err != nil {
		return
	}

	if code != 200 && code != 201 {
		fmt.Printf("deposit request: expected http code 200 or 201, received: httpCode:%v, url:%v \n", code, url)
		err = fmt.Errorf("deposit request: payout request:expected http code 200 or 201, received: httpCode:%v", code)
		return
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("json Unmarshal err:%v, httpCode:%v", err, code)
		return
	}

	return
}

// SetMockResponse will set *DepositResult as mocked
// response on the next deposit request
// To be used only for test purposes
func (mock *DepositResult) SetMockResponse() {
	mockedDepositResult = mock
}

// validate the instance of DepositOrder
// if not valid returns an error
func (d *DepositOrder) validate() error {
	required := []string{"MerchantOrderID", "MerchantOrderDesc", "OrderAmount", "OrderCurrency", "MerchantOrderDesc", "CustomerEmail", "CustomerLastName", "CustomerAddress", "CustomerCountryCode", "CustomerCity", "CustomerZipCode", "CustomerPhone", "CustomerIP", "RedirectURL", "CheckoutURL"}
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
