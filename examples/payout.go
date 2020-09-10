package main

import (
	"fmt"
	"github.com/zotapay/go-sdk/zotapay"
)

// payout example init payout function
func payout() {
	var sdk = zotapay.SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        zotapay.SANDBOX,
	}

	res, err := sdk.Payout(zotapay.PayoutOrder{
		MerchantOrderID:           "134e4f443t651",
		MerchantOrderDesc:         "Test order description",
		OrderAmount:               "500",
		OrderCurrency:             "MYR",
		CustomerEmail:             "customer@email-address.com",
		CustomerFirstName:         "John",
		CustomerLastName:          "Doe",
		CustomerPhone:             "+1 420-100-1000",
		CustomerIP:                "127.0.0.1",
		CustomerBankCode:          "BBL",
		CallbackURL:               "https://some.endpoint/callback",
		CheckoutURL:               "https://some.endpoint/checkout",
		CustomerBankAccountNumber: "100200",
		CustomerBankAccountName:   "John Doe",
		CustomerBankBranch:        "Bank Branch",
		CustomerBankAddress:       "Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
		CustomerBankZipCode:       "84280",
		CustomerBankProvince:      "Bank Province",
		CustomerBankArea:          "Bank Area / City",
		CustomerBankRoutingNumber: "000",
	})

	if err != nil {
		fmt.Printf("sdk error:%v \n", err)
		return
	}

	if res.Code != "200" {
		fmt.Printf("non-successful payout response from Zotapay server code:%v, error message:%v \n", res.Code, res.Message)
		return
	}

	fmt.Printf("successful payout response from Zotapay server code:%v, order ID:%v, merchant order ID:%v\n",
		res.Code, res.Data.OrderID, res.Data.MerchantOrderID)
}

// payoutMocked example mocking payout result
// only for test purposes
func payoutMocked() {
	var sdk = zotapay.SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        zotapay.SANDBOX,
	}

	// ---------------only for test purposes-----------------
	//init and set the struct that will be mocked as response
	//only for next execution of sdk.Payout()
	mock := &zotapay.PayoutResult{
		Code: "200",
		Data: zotapay.PayoutResultData{
			MerchantOrderID: "123",
			OrderID:         "1234",
		},
		Message: "SomeMockMsg",
	}
	mock.SetMockResponse()
	// ------------------------------------------------------

	res, err := sdk.Payout(zotapay.PayoutOrder{
		MerchantOrderID:           "134e4f443t651",
		MerchantOrderDesc:         "Test order description",
		OrderAmount:               "500",
		OrderCurrency:             "MYR",
		CustomerEmail:             "customer@email-address.com",
		CustomerFirstName:         "John",
		CustomerLastName:          "Doe",
		CustomerPhone:             "+1 420-100-1000",
		CustomerIP:                "127.0.0.1",
		CustomerBankCode:          "BBL",
		CallbackURL:               "https://some.endpoint/callback",
		CheckoutURL:               "https://some.endpoint/checkout",
		CustomerBankAccountNumber: "100200",
		CustomerBankAccountName:   "John Doe",
		CustomerBankBranch:        "Bank Branch",
		CustomerBankAddress:       "Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
		CustomerBankZipCode:       "84280",
		CustomerBankProvince:      "Bank Province",
		CustomerBankArea:          "Bank Area / City",
		CustomerBankRoutingNumber: "000",
	})

	if err != nil {
		fmt.Printf("sdk error:%v \n", err)
		return
	}

	if res != *mock {
		fmt.Printf("the response is not successfully mocked\n")
		return
	}

	fmt.Printf("successful mocked response from Zotapay server code:%v, order ID:%v, merchant order ID:%v \n",
		res.Code, res.Data.OrderID, res.Data.MerchantOrderID)
}
