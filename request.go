package bankmuscatpg

import (
	"fmt"
	"net/url"
	"reflect"
)

type merchantData struct {
	TrackingId         int     `mapkey:"tid"`
	MerchantId         int     `mapkey:"merchant_id"`
	OrderId            int     `mapkey:"order_id"`
	Currency           string  `mapkey:"currency"`
	RedirectUrl        string  `mapkey:"redirect_url"`
	CancelUrl          string  `mapkey:"cancel_url"`
	Language           string  `mapkey:"language"`
	Amount             float32 `mapkey:"amount"`
	BName              *string `mapkey:"billing_name"`
	BAddress           *string `mapkey:"billing_address"`
	BCity              *string `mapkey:"billing_city"`
	BState             *string `mapkey:"billing_state"`
	BCountry           *string `mapkey:"billing_country"`
	BMail              *string `mapkey:"billing_email"`
	BZip               *int    `mapkey:"billing_zip"`
	BPhone             *int    `mapkey:"billing_tel"`
	SName              *string `mapkey:"delivery_name"`
	SAddress           *string `mapkey:"delivery_address"`
	SCity              *string `mapkey:"delivery_city"`
	SState             *string `mapkey:"delivery_state"`
	SCountry           *string `mapkey:"delivery_country"`
	SZip               *int    `mapkey:"delivery_zip"`
	SPhone             *int    `mapkey:"delivery_tel"`
	MerchantParam1     *string `mapkey:"merchant_param1"`
	MerchantParam2     *string `mapkey:"merchant_param2"`
	MerchantParam3     *string `mapkey:"merchant_param3"`
	MerchantParam4     *string `mapkey:"merchant_param4"`
	MerchantParam5     *string `mapkey:"merchant_param5"`
	PromoCode          *string `mapkey:"promo_code"`
	CustomerIdentifier *string `mapkey:"customer_identifier"`
}

type RequestInfo struct {
	TId, OrderId                                                                                                  int
	Amount                                                                                                        float32
	MerchantParam1, MerchantParam2, MerchantParam3, MerchantParam4, MerchantParam5, PromoCode, CustomerIdentifier *string
}

type BillingReqInfo struct {
	Name, Address, City, State, Country, Mail *string
	ZIP, Phone                                *int
}

type ShippingReqInfo struct {
	Name, Address, City, State, Country, Mail *string
	ZIP, Phone                                *int
}

// Creating a request using provided data, encrypt it and then pass it to createRequest function
func (Bmpg *BankMuscatPG) Request(rI RequestInfo, bRI BillingReqInfo, sRI ShippingReqInfo) (string, error) {
	requestMap := make(map[string]interface{})

	request := merchantData{
		TrackingId:         rI.TId,
		MerchantId:         Bmpg.MerchantId,
		OrderId:            rI.OrderId,
		Currency:           *Bmpg.Currency,
		RedirectUrl:        Bmpg.CallbackUrl,
		CancelUrl:          Bmpg.CallbackUrl,
		Language:           *Bmpg.Language,
		Amount:             rI.Amount,
		BName:              bRI.Name,
		BAddress:           bRI.Address,
		BCity:              bRI.City,
		BState:             bRI.State,
		BCountry:           bRI.Country,
		BMail:              bRI.Mail,
		BZip:               bRI.ZIP,
		BPhone:             bRI.Phone,
		SName:              sRI.Name,
		SAddress:           sRI.Address,
		SCity:              sRI.Address,
		SState:             sRI.State,
		SCountry:           sRI.Country,
		SZip:               sRI.ZIP,
		SPhone:             sRI.Phone,
		MerchantParam1:     rI.MerchantParam1,
		MerchantParam2:     rI.MerchantParam2,
		MerchantParam3:     rI.MerchantParam3,
		MerchantParam4:     rI.MerchantParam4,
		MerchantParam5:     rI.MerchantParam5,
		PromoCode:          rI.PromoCode,
		CustomerIdentifier: rI.CustomerIdentifier,
	}

	v := reflect.ValueOf(request)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("mapkey")
		if tag != "" {
			requestMap[tag] = field.Interface()
		}
	}

	requestData := mapToString(requestMap)

	encRequest, err := getAES256GCMEncrypted(requestData, Bmpg.WorkingKey)
	if err != nil {
		return "", err
	}

	response, err := createRequest(encRequest, Bmpg.AccessCode, *Bmpg.TestEnv)
	if err != nil {
		return "", err
	}

	return response, nil
}

func mapToString(requestMap map[string]interface{}) string {
	var requestData string
	for key, value := range requestMap {
		requestData += key + "=" + url.QueryEscape(fmt.Sprint(value)) + "&"
	}

	if len(requestData) > 0 {
		requestData = requestData[:len(requestData)-1]
	}
	return requestData
}

// Make the request page HTML using the encrypted data, this page must be rendered directly to the user
func createRequest(encRequest string, accessCode string, testEnv bool) (response string, err error) {
	redirectUrl := "https://smartpaytrns.bankmuscat.com/transaction.do?command=initiateTransaction"
	if testEnv {
		redirectUrl = "https://mti.bankmuscat.com:6443/transaction.do?command=initiateTransaction"
	}

	render := fmt.Sprintf(
		`<form method="post" name="redirect" action="%s">
			<input type=hidden name=encRequest value=%s>
			<input type=hidden name=access_code value=%s>
		</form>
		<script language='javascript'>document.redirect.submit();</script>`,
		redirectUrl, encRequest, accessCode)

	return render, nil
}
