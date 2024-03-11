package zota

import (
	"encoding/json"
	"fmt"
)

// callback struct - handle the whole callback data
type callback struct {
	Type                   string      `json:"type"`
	Amount                 string      `json:"amount"`
	Status                 string      `json:"status"`
	OrderID                string      `json:"orderID"`
	Currency               string      `json:"currency"`
	ExtraData              interface{} `json:"extraData"`
	Signature              string      `json:"signature"`
	EndpointID             string      `json:"endpointID"`
	CustomParam            string      `json:"customParam"`
	ErrorMessage           string      `json:"errorMessage"`
	CustomerEmail          string      `json:"customerEmail"`
	MerchantOrderID        string      `json:"merchantOrderID"`
	OriginalRequest        interface{} `json:"originalRequest"`
	ProcessorTransactionID string      `json:"processorTransactionID"`
}

// Callback parse the callback data into callback struct
// and validate the Signature
func (s *SDK) Callback(b []byte) (callback, error) {

	var c callback
	err := json.Unmarshal(b, &c)
	if err != nil {
		return callback{}, fmt.Errorf("unexpected callback json:%v", err)
	}

	isValid := s.validateCallbackSign(&c)
	if isValid != true {
		return callback{}, fmt.Errorf("wrong signature")
	}

	return c, nil
}

// validateCallbackSign generate and validate
// callback sign
func (s *SDK) validateCallbackSign(c *callback) bool {
	signature := s.sign(c.EndpointID + c.OrderID + c.MerchantOrderID + c.Status + c.Amount + c.CustomerEmail)
	return signature == c.Signature
}
