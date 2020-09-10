package zotapay

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// testC test base structure
type testC struct {
	name          string
	expected      callback
	mockSDK       SDK
	mockCb        string
	expectedError error
}

func TestSDK_Callback(t *testing.T) {
	tests := []testC{
		{
			name: `Success`,
			mockCb: `{
    "type":"SALE",
    "status":"FILTERED",
    "errorMessage":"dummy sandbox filter",
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
    "signature":"cceacf126acf7bc77745c77f596835c3d4c0426ebb49e615559cc91589e7cff9"
}`,
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
			},
			expectedError: nil,
			expected:      callback{Type: "SALE", Amount: "100.00", Status: "FILTERED", OrderID: "24043630", Currency: "USD", ExtraData: map[string]interface{}{"billingDescriptor": "sandbox-payment", "card": map[string]interface{}{"cvv": "***", "expiration": "03/2023", "holder": "A A", "number": "000000***1111"}, "cardData": map[string]interface{}{"bank": map[string]interface{}{}, "brand": "", "country": map[string]interface{}{}}, "dcc": false, "paymentMethod": "CREDITCARD"}, Signature: "cceacf126acf7bc77745c77f596835c3d4c0426ebb49e615559cc91589e7cff9", EndpointID: "503364", CustomParam: "", ErrorMessage: "dummy sandbox filter", CustomerEmail: "testing@api-requests.com", MerchantOrderID: "172", OriginalRequest: map[string]interface{}{"callbackUrl": "https://example.com/callback.php", "checkoutUrl": "https://example.com/checkout.php", "customerAddress": "The Swan, Jungle St. 108", "customerCity": "Los Angeles", "customerCountryCode": "US", "customerEmail": "testing@api-requests.com", "customerFirstName": "John", "customerIP": "134.201.250.130", "customerLastName": "Lock", "customerPhone": "1-420-100-1000", "customerState": "CA", "customerZipCode": "90015", "merchantOrderDesc": "Test order description", "merchantOrderID": "172", "orderAmount": "100.00", "orderCurrency": "USD", "redirectUrl": "https://example.com/redirect.php", "requestedAt": "0001-01-01T00:00:00Z", "signature": "4b92c1b81807a302e4db98028b6fe9bfb94d802df0d0582798ae416119184e5a"}, ProcessorTransactionID: "ee7c9d1b-5e68-4408-9794-5d40aedead6c"},
		}, {
			name: `WrongSign`,
			mockCb: `{
    "type":"SALE",
    "status":"FILTERED",
    "errorMessage":"dummy sandbox filter",
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
    "signature":"e89454f30a99349084ab2581a526c3b7d4f5a5ab7c55015f819c74e24012670f1"
}`,
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
			},
			expectedError: fmt.Errorf("wrong signature"),
			expected:      callback{},
		}, {
			name:   `Unexpected`,
			mockCb: `some unexpected`,
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "API_MERCHANT_SECRET_KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
			},
			expectedError: fmt.Errorf("unexpected callback json:invalid character 's' looking for beginning of value"),
			expected:      callback{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cb, err := test.mockSDK.Callback([]byte(test.mockCb))
			assert.IsType(t, callback{}, cb)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expected, cb)
		})
	}
}
