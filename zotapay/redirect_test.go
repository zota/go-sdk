package zotapay

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

// testRedirect test base structure
type testRedirect struct {
	name            string
	expected        redirect
	mockSDK         SDK
	mockRedirectUrl url.URL
	expectedError   error
}

func TestSDK_redirect(t *testing.T) {
	tests := []testRedirect{
		{
			name:            `Success`,
			mockRedirectUrl: url.URL{RawQuery: "orderID=12345678&merchantOrderID=1&errorMessage=&billingDescriptor=sandbox-payment&signature=6a4f1ad55ee636e65b8aece10b1025f28566c2896b23d623a42e101b905d043c&status=APPROVED"},
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "MERCHANT-SECRET-KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
			},
			expectedError: nil,
			expected:      redirect{OrderID: "12345678", MerchantOrderID: "1", Status: "APPROVED", BillingDescriptor: "sandbox-payment", Signature: "6a4f1ad55ee636e65b8aece10b1025f28566c2896b23d623a42e101b905d043c"},
		}, {
			name:            `WrongSign`,
			mockRedirectUrl: url.URL{RawQuery: "orderID=12345678&merchantOrderID=1&errorMessage=&billingDescriptor=sandbox-payment&signature=16a4f1ad55ee636e65b8aece10b1025f28566c2896b23d623a42e101b905d043c&status=APPROVED"},
			mockSDK: SDK{
				MerchantID:        "API_MERCHANT_ID",
				MerchantSecretKey: "MERCHANT-SECRET-KEY",
				EndpointID:        "503368",
				ApiBaseURL:        SANDBOX,
			},
			expectedError: fmt.Errorf("wrong signature"),
			expected:      redirect{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cb, err := test.mockSDK.Redirect(test.mockRedirectUrl)
			assert.IsType(t, redirect{}, cb)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expected, cb)
		})
	}
}
