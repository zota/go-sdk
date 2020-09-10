package main

import (
	"fmt"
	"github.com/zotapay/go-sdk/zotapay"
)

// ordersReport example init orderStatus function
func ordersReport() {
	var sdk = zotapay.SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        zotapay.SANDBOX,
	}

	res, err := sdk.OrdersReport(zotapay.OrdersReport{
		DateType:    "created",
		EndpointIds: "503368,503365",
		FromDate:    "2020-08-01",
		ToDate:      "2020-09-01",
		Statuses:    "APPROVED,DECLINED",
		Types:       "SALE,PAYOUT",
	})

	if err != nil {
		fmt.Printf("sdk error:%v \n", err)
		return
	}

	if res.Code != "200" {
		fmt.Printf("non-successful response from Zotapay server code:%v, error message:%v \n", res.Code, res.Message)
		return
	}

	fmt.Printf("successful response from Zotapay server code:%v, order Report:%v \n",
		res.Code, res.OrdersReport)
}
