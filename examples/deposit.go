package main

import (
	"fmt"
	"github.com/zota/go-sdk/zota"
)

// deposit example init deposit function
func deposit() {
	var sdk = zota.SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        zota.SANDBOX,
	}

	res, err := sdk.Deposit(zota.DepositOrder{
		MerchantOrderID:     "134e4f44t65111",
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
	})

	if err != nil {
		fmt.Printf("sdk error:%v \n", err)
		return
	}

	if res.Code != "200" {
		fmt.Printf("non-successful response from zota server code:%v, error message:%v \n", res.Code, res.Message)
		return
	}

	fmt.Printf("successful response from zota server code:%v, order ID:%v, merchant order ID:%v, deposit URL:%v \n",
		res.Code, res.Data.OrderID, res.Data.MerchantOrderID, res.Data.DepositURL)
}

// depositMocked example mocking deposit result
// only for test purposes
func depositMocked() {
	var sdk = zota.SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        zota.SANDBOX,
	}

	// ---------------only for test purposes-----------------
	//init and set the struct that will be mocked as response
	//only for next execution of sdk.Deposit()
	mock := &zota.DepositResult{
		Code: "200",
		Data: zota.DepositResultData{
			DepositURL:      "http://some.mock",
			MerchantOrderID: "123",
			OrderID:         "1234",
		},
		Message: "SomeMockMsg",
	}
	mock.SetMockResponse()
	// ------------------------------------------------------

	res, err := sdk.Deposit(zota.DepositOrder{
		MerchantOrderID:     "134e4f44t651",
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
	})

	if err != nil {
		fmt.Printf("sdk error:%v \n", err)
		return
	}

	if res != *mock {
		fmt.Printf("the response is not successfully mocked\n")
		return
	}

	fmt.Printf("successful mocked response from zota server code:%v, order ID:%v, merchant order ID:%v, deposit URL:%v \n",
		res.Code, res.Data.OrderID, res.Data.MerchantOrderID, res.Data.DepositURL)
}
