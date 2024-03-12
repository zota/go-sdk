[![Codecov.io](https://codecov.io/gh/zota/go-sdk/graph/badge.svg?token=WKOX8Lm3My)](https://codecov.io/gh/zota/go-sdk)
[![Go Matrix Test](https://github.com/zota/go-sdk/actions/workflows/test-pipeline.yml/badge.svg)](https://github.com/zota/go-sdk/actions/workflows/test-pipeline.yml)
[![Golang Pipeline](https://github.com/zota/go-sdk/actions/workflows/go-quality.yml/badge.svg)](https://github.com/zota/go-sdk/actions/workflows/go-quality.yml)
![golang-github](https://user-images.githubusercontent.com/174284/106497798-2cee0d00-64c7-11eb-9014-9e0d8c4231cf.jpg)

# Official Zota Go SDK

This is the official page of the [Zota](https://www.zota.com) Go SDK. It is intended to be used by developers who run modern Go applications and would like to integrate our next generation payments platform.

## Introduction

Go SDK provides all the necessary methods for integrating the zota Merchant API.

## Requirements

- A functioning Zota Sandbox or Production account and related credentials (`MerchantID`, `MerchantSecretKey`, `EndpointID`)
- Go 1.20 or greater

## Usage

### Main configuration

After all the files are loaded configuration is needed. This can be done initiation of `zota.SDK{}`. Configuration includes:

- Credentials
- Endpoint API url - test or production environment 

### API requests

After everything is setup all requests to the API are made with the corresponding methods:

- Deposit
- Payout
- Callback
- Order Status
- Orders Report

All the methods belongs to `zota.SDK` struct which we discussed in configuration section.

### Making the request

First the data object has to be created with all required fields (ex. `zota.DepositOrder{...}`).

After that the request method is called with the data object as parameter. (ex. `sdk.Deposit(zota.DepositOrder{...})`)

### Retrieving the response

Every request method returns response and error objects. The error needs to be handled properly. If error is equal to _nil_ you can access the Http code and the Data object in the response.

### Callback

Method for callback handling is available:

```golang
sdk.Callback(CallbackRequestBody)
```

## Examples

Examples are available in `examples` folder.

Requests:

- `deposit.go` - Deposit request
- `payout.go` - Payout request
- `orderStatus.go` - Order status request
- `ordersReport.go` - Orders report request

Order Handlers:

- `callback.go` - API Callback

## Code Test Coverage

[![codecov](https://codecov.io/gh/zota/go-sdk/graphs/tree.svg?width=650&height=150&src=pr&token=WKOX8Lm3My)](https://codecov.io/gh/zota/go-sdk/)
> Codecov.io visualisation of code blocks and their test coverage.

## Resources

The Zota API guide can be found on the official API Documentation pages for [deposit](https://doc.zota.com/deposit/1.0/) and [payout](https://doc.zota.com/payout/1.0/) operations.

## Support

This SDK is developed and supported by Zota. For sign-up and sales inquiries, please contact us through our official website. For technical inquiries, please visit our Github page and open an issue.
