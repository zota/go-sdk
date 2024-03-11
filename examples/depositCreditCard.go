package main

import (
	"fmt"

	"github.com/zota/go-sdk/zota"
)

// credit card deposit example init deposit function
func depositCC() {
	var sdk = zota.SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        zota.SANDBOX,
	}

	res, err := sdk.DepositCC(zota.DepositCCOrder{
		MerchantOrderID:     "134e4f44t651112121",
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

	if err != nil {
		fmt.Printf("sdk error:%v \n", err)
		return
	}

	if res.Code != "200" {
		fmt.Printf("non-successful response from zota server code:%v, error message:%v \n", res.Code, res.Message)
		return
	}

	fmt.Printf("successful response from zota server code:%v, order ID:%v, merchant order ID:%v, deposit Status:%v \n",
		res.Code, res.Data.OrderID, res.Data.MerchantOrderID, res.Data.Status)
}
