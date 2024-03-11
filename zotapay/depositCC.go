package zota

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// DepositCCOrder represents credit card deposit order
type DepositCCOrder struct {
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

	//credit card data
	CardNumber          string `json:"cardNumber"`
	CardHolderName      string `json:"cardHolderName"`
	CardExpirationMonth string `json:"cardExpirationMonth"`
	CardExpirationYear  string `json:"cardExpirationYear"`
	CardCvv             string `json:"cardCvv"`
}

// DepositCCResult represents credit card deposit response from Zota API
type DepositCCResult struct {
	Code    string              `json:"code"`
	Data    DepositCCResultData `json:"data"`
	Message string              `json:"message"`
}

type DepositCCResultData struct {
	Status          string `json:"status"`
	MerchantOrderID string `json:"merchantOrderID"`
	OrderID         string `json:"orderID"`
}

var mockedDepositCCResult *DepositCCResult

// DepositCC init validation of the SDK struct and the DepositCCOrder
// generate sign and
// init a credit card deposit request to Zota API
func (s *SDK) DepositCC(d DepositCCOrder) (res DepositCCResult, err error) {

	//validate that SDK is properly initialized
	err = s.validate()
	if err != nil {
		return
	}

	//validate that DepositCCOrder is properly initialized
	err = d.validate()
	if err != nil {
		return
	}

	//if mockedDepositCCResult is set return it as response
	//only for testing
	if mockedDepositCCResult != nil {
		res = *mockedDepositCCResult
		mockedDepositCCResult = nil
		return
	}

	//generate signature
	d.Signature = s.sign(s.EndpointID, d.MerchantOrderID, d.OrderAmount, d.CustomerEmail)

	deposit, err := json.Marshal(d)
	if err != nil {
		return
	}

	_, body, err := s.httpDo(http.MethodPost, fmt.Sprintf("%v/api/v1/deposit/request/%v/", s.ApiBaseURL, s.EndpointID), deposit)
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

// SetMockResponse will set *DepositCCResult as mocked
// response on the next deposit request
// To be used only for test purposes
func (mock *DepositCCResult) SetMockResponse() {
	mockedDepositCCResult = mock
}

// validate the instance of DepositCCOrder
// if not valid returns an error
func (d *DepositCCOrder) validate() error {
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
