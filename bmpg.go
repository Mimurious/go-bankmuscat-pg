package bankmuscatpg

type BankMuscatPG struct {
	MerchantId                          int
	AccessCode, WorkingKey, CallbackUrl string
	Currency                            *string
	TestEnv                             *bool
	Language                            *string
}

func New(params BankMuscatPG) *BankMuscatPG {
	Bmpg := BankMuscatPG{
		MerchantId:  params.MerchantId,
		AccessCode:  params.AccessCode,
		CallbackUrl: params.CallbackUrl,
		WorkingKey:  params.WorkingKey,
		Currency:    params.Currency,
		TestEnv:     params.TestEnv,
		Language:    params.Language,
	}

	if Bmpg.TestEnv == nil {
		defaultEnv := false
		Bmpg.TestEnv = &defaultEnv
	}

	if Bmpg.Currency == nil {
		defaultCurrency := "OMR"
		Bmpg.Currency = &defaultCurrency
	}

	if Bmpg.Language == nil {
		defaultLanguage := "EN"
		Bmpg.Language = &defaultLanguage
	}

	return &Bmpg
}
