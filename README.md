# Bank Muscat Payment Gateway Implemention for Go
This Go library provides a simple and efficient way to integrate with the Bank Muscat Payment Gateway for online payment processing. With this library, you can easily handle payment requests and responses, making it suitable for e-commerce websites and other applications that require secure payment processing.

## Installation
To use this library in your Go project, you can install it using `go get`:
```
go get github.com/Mimurious/go-bankmuscat-pg
```

## Usage
Here's an example of how to use this library to initiate a payment request:

```
import (
	"fmt"
	bankmuscatpg "github.com/Mimurious/go-bankmuscat-pg"
)

func main() {
    // Set to true for test environment, false for production
    TestEnv := true

    // Create a new BMPG (Bank Muscat Payment Gateway) object
    bmpg := bankmuscatpg.New(bankmuscatpg.BankMuscatPG{
        MerchantId:  ***REMOVED***,
        AccessCode:  "my16digitIvKey12",
        WorkingKey:  "my32digitkey12345678901234567890",
        TestEnv:     &TestEnv,
        CallbackUrl: "example.com",
    })

    // Initiate a payment request
    response, err := bmpg.Request(bankmuscatpg.RequestInfo{
        TId:     111111,
        OrderId: 1111111,
        Amount:  22.11,
    }, bankmuscatpg.BillingReqInfo{}, bankmuscatpg.ShippingReqInfo{})

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(response)
}
```

### Creating a BMPG Object
To use this library, you need to create a `BankMuscatPG` object using the `bankmuscatpg.New` function. This object contains the following required parameters:

- MerchantId: Your Bank Muscat merchant ID.
- AccessCode: Your 16-digit access code.
- WorkingKey: Your 32-digit working key.
- CallbackUrl: The URL to which Bank Muscat will send payment responses.

You can also specify optional parameters such as Currency, TestEnv, and Language within the BankMuscatPG object.

### Initiating a Payment Request
Once you have created a `BankMuscatPG` object, you can use it to initiate a payment request. Provide the necessary payment information in the `RequestInfo` struct and any additional billing and shipping information as needed.

```
type BillingReqInfo struct {
	Name, Address, City, State, Country, Mail *string
	ZIP, Phone                                *int
}

type ShippingReqInfo struct {
	Name, Address, City, State, Country, Mail *string
	ZIP, Phone                                *int
}
```

### Handling the Response
After initiating a payment request, you'll receive a response from the Bank Muscat Payment Gateway. To extract meaningful information from the response, follow these steps:

1. Decrypt the Response: You must use the bankmuscatpg.DecryptAES256GCM function to decrypt the response data. This function will ensure that the data is securely decrypted for further processing.

2. Interpret the Response: The decrypted response will contain various fields that provide information about the transaction status, payment details, and more. Refer to the [Official API documents](https://mti.bankmuscat.com:7443/kitLibrary/kits/download/SmartPay_Integration_Guide.pdf) for a comprehensive guide on interpreting the response data. This document provides detailed explanations of the response fields and their meanings.

3. Handle Errors: In case of errors during the payment process, you should consult the [SmartPay Error Codes document](https://mti.bankmuscat.com:7443/kitLibrary/kits/download/SmartPay_Error_Codes.pdf) to understand the error codes and their corresponding descriptions. This will help you diagnose and resolve any issues that may arise during payment processing.
