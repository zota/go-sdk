package zota

import (
	"fmt"
	"net/url"
)

// redirect struct - handle the redirect query params
type redirect struct {
	Signature         string
	BillingDescriptor string
	ErrorMessage      string
	MerchantOrderID   string
	OrderID           string
	Status            string
}

// Redirect parse the redirect query params into redirect struct
// and validate the Signature
func (s *SDK) Redirect(url url.URL) (r redirect, err error) {

	r.OrderID = url.Query().Get("orderID")
	r.MerchantOrderID = url.Query().Get("merchantOrderID")
	r.ErrorMessage = url.Query().Get("errorMessage")
	r.BillingDescriptor = url.Query().Get("billingDescriptor")
	r.Signature = url.Query().Get("signature")
	r.Status = url.Query().Get("status")

	isValid := s.validateredirectSign(&r)
	if isValid != true {
		return redirect{}, fmt.Errorf("wrong signature")
	}

	return
}

// validateredirectSign generate and validate
// redirect sign
func (s *SDK) validateredirectSign(c *redirect) bool {
	signature := s.sign(c.Status, c.OrderID, c.MerchantOrderID)
	return signature == c.Signature
}
