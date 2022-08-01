package zotapay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// PayoutOrder represents payout order
type PayoutOrder struct {
	MerchantOrderID                string `json:"merchantOrderID"`
	MerchantOrderDesc              string `json:"merchantOrderDesc"`
	OrderAmount                    string `json:"orderAmount"`
	OrderCurrency                  string `json:"orderCurrency"`
	CustomerEmail                  string `json:"customerEmail"`
	CustomerFirstName              string `json:"customerFirstName"`
	CustomerLastName               string `json:"customerLastName"`
	CustomerPhone                  string `json:"customerPhone"`
	CustomerIP                     string `json:"customerIP"`
	CallbackURL                    string `json:"callbackUrl"`
	CustomerBankCode               string `json:"customerBankCode"`
	CustomerBankAccountNumber      string `json:"customerBankAccountNumber"`
	CustomerBankAccountName        string `json:"customerBankAccountName"`
	CustomerBankBranch             string `json:"customerBankBranch"`
	CustomerBankAddress            string `json:"customerBankAddress"`
	CustomerBankZipCode            string `json:"customerBankZipCode"`
	CustomerBankProvince           string `json:"customerBankProvince"`
	CustomerBankArea               string `json:"customerBankArea"`
	CustomerBankRoutingNumber      string `json:"customerBankRoutingNumber"`
	CustomParam                    string `json:"customParam"`
	CheckoutURL                    string `json:"checkoutUrl"`
	RedirectUrl                    string `json:"redirectUrl"`
	CustomerCountryCode            string `json:"customerCountryCode"`
	CustomerPersonalID             string `json:"customerPersonalID"`
	CustomerBankAccountNumberDigit string `json:"customerBankAccountNumberDigit"`
	CustomerBankAccountType        string `json:"customerBankAccountType"`
	CustomerBankSwiftCode          string `json:"customerBankSwiftCode"`
	CustomerBankBranchDigit        string `json:"customerBankBranchDigit"`
	Signature                      string `json:"signature"`
}

// PayoutResult represents the payout response from Zotapay API
type PayoutResult struct {
	Code    string           `json:"code"`
	Data    PayoutResultData `json:"data"`
	Message string           `json:"message"`
}

type PayoutResultData struct {
	MerchantOrderID string `json:"merchantOrderID"`
	OrderID         string `json:"orderID"`
}

var mockedPayoutResult *PayoutResult

// Payout init validation of the SDK struct and PayoutOrder
// generate sign and
// init a payout request to Zotapay API
func (s *SDK) Payout(p PayoutOrder) (res PayoutResult, err error) {

	//validate that SDK is properly initialized
	err = s.validate()
	if err != nil {
		return
	}

	//validate that PayoutOrder is properly initialized
	err = p.validate()
	if err != nil {
		return
	}

	//if mockedPayoutResult is set return it as response
	//only for testing
	if mockedPayoutResult != nil {
		res = *mockedPayoutResult
		mockedPayoutResult = nil
		return
	}

	//generate signature
	p.Signature = s.sign(s.EndpointID, p.MerchantOrderID, p.OrderAmount, p.CustomerEmail, p.CustomerBankAccountNumber)

	payout, err := json.Marshal(p)
	if err != nil {
		return
	}

	_, body, err := s.httpDo(http.MethodPost, fmt.Sprintf("%v/api/v1/payout/request/%v/", s.ApiBaseURL, s.EndpointID), payout)
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

// SetMockResponse will set *PayoutResult as mocked
// response on the next payout request
// To be used only for test purposes
func (mock *PayoutResult) SetMockResponse() {
	mockedPayoutResult = mock
}

// validate the instance of PayoutOrder
// if not valid returns an error
func (d *PayoutOrder) validate() error {
	required := []string{"MerchantOrderID", "MerchantOrderDesc", "OrderAmount", "OrderCurrency", "CustomerBankAccountNumber", "CustomerBankAccountName"}
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
