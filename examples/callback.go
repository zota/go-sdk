package main

import (
	"fmt"
	"github.com/zota/go-sdk/zota"
)

var mockCallback = `{
    "type":"SALE",
    "status":"AUTHORIZED",
    "errorMessage":"",
    "endpointID":"503364",
    "processorTransactionID":"ee7c9d1b-5e68-4408-9794-5d40aedead6c",
    "orderID":"24043630",
    "merchantOrderID":"172",
    "amount":"100.00",
    "currency":"USD",
    "customerEmail":"testing@api-requests.com",
    "customParam":"",
    "extraData":{
        "billingDescriptor":"sandbox-payment",
        "card":{
            "cvv":"***",
            "expiration":"03/2023",
            "holder":"A A",
            "number":"000000***1111"
        },
        "cardData":{
            "bank":{},
            "brand":"",
            "country":{}
        },
        "dcc":false,
        "paymentMethod":"CREDITCARD"
    },
    "originalRequest":{
        "merchantOrderID":"172",
        "merchantOrderDesc":"Test order description",
        "orderAmount":"100.00","orderCurrency":"USD",
        "customerEmail":"testing@api-requests.com",
        "customerFirstName":"John",
        "customerLastName":"Lock",
        "customerAddress":"The Swan, Jungle St. 108",
        "customerCountryCode":"US","customerCity":"Los Angeles",
        "customerState":"CA",
        "customerZipCode":"90015",
        "customerPhone":"1-420-100-1000",
        "customerIP":"134.201.250.130",
        "redirectUrl":"https://example.com/redirect.php",
        "callbackUrl":"https://example.com/callback.php",
        "checkoutUrl":"https://example.com/checkout.php",
        "signature":"4b92c1b81807a302e4db98028b6fe9bfb94d802df0d0582798ae416119184e5a",
        "requestedAt":"0001-01-01T00:00:00Z"
    },
    "signature":"90d2168c785573b654102d007f0fcabf44f2be31cfbd4b48febdd779d145a2fd"
}`

// callback example for parsing callback
func callback() {
	var sdk = zota.SDK{
		MerchantID:        "API_MERCHANT_ID",
		MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
		EndpointID:        "503368",
		ApiBaseURL:        zota.SANDBOX,
	}

	res, err := sdk.Callback([]byte(mockCallback))

	if err != nil {
		fmt.Printf("sdk error:%v \n", err)
		return
	}

	if res.ErrorMessage != "" {
		fmt.Printf("the transaction is declined or return an error, error message:%v \n", res.ErrorMessage)
		return
	}

	fmt.Printf("successful callback received from zota order ID:%v, merchant order ID:%v, order status:%v\n",
		res.OrderID, res.MerchantOrderID, res.Status)

}
