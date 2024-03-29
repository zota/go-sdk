package main

import (
	"fmt"
	"net/url"

	"github.com/zota/go-sdk/zota"
)

var mockRedirect = ``

// redirect example for parsing redirect
func redirect() {
	var sdk = zota.SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "MERCHANT-SECRET-KEY",
		EndpointID:        "503368",
		ApiBaseURL:        zota.SANDBOX,
	}

	res, err := sdk.Redirect(url.URL{RawQuery: "orderID=12345678&merchantOrderID=1&errorMessage=&billingDescriptor=sandbox-payment&signature=6a4f1ad55ee636e65b8aece10b1025f28566c2896b23d623a42e101b905d043c&status=APPROVED"})

	if err != nil {
		fmt.Printf("sdk error:%v \n", err)
		return
	}

	if res.ErrorMessage != "" {
		fmt.Printf("Zota api return an error, error message:%v \n", res.ErrorMessage)
		return
	}

	fmt.Printf("successful redirect received from Zota order ID:%v, merchant order ID:%v, order status:%v\n",
		res.OrderID, res.MerchantOrderID, res.Status)

}
