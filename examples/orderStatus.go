package main

import (
	"fmt"
	"github.com/zota/go-sdk/zota"
)

// orderStatus example init orderStatus function
func orderStatus() {
	var sdk = zota.SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        zota.SANDBOX,
	}

	res, err := sdk.OrderStatus(zota.OrderStatus{
		MerchantOrderID: "134e4f44t65111",
		OrderID:         "24057644",
	})

	if err != nil {
		fmt.Printf("sdk error:%v \n", err)
		return
	}

	if res.Code != "200" {
		fmt.Printf("non-successful response from zota server code:%v, error message:%v \n", res.Code, res.Message)
		return
	}

	fmt.Printf("successful response from zota server code:%v, order ID:%v, merchant order ID:%v, Status:%v \n",
		res.Code, res.OrderID, res.MerchantOrderID, res.Status)
}
