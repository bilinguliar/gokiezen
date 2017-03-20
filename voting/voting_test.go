package voting

import (
	"errors"
	"testing"
)

func TestRegisterVote(t *testing.T) {
	stats := make(map[string]int)
	countries := make(map[string]bool)
	messages := make(map[string]string)

	svc := New(
		&MessengerMock{
			RequestSMSFunc: func(originator, recipient, text string) {
				messages[recipient] = text
			},
		},
		&EnquirerMock{},
		&SkoreKprMock{
			AddPointFunc: func(key string) error {
				stats[key]++
				return nil
			},
			AddCountryFunc: func(code string) error {
				countries[code] = true
				return nil
			},
			GetAllCandidatesFunc: func() ([]string, error) {
				return nil, nil
			},
			GetAllCountriesFunc: func() ([]string, error) {
				return nil, nil
			},
			GetFunc: func(key string) (int, error) {
				return stats[key], nil
			},
		},
		"EuroVision",
	)

	msisdn := "380661234567"
	abba := "ABBA"

	svc.RegisterVote(msisdn, abba)

	if stats[abba] != 1 {
		t.Error("Score for ABBA is %d, expected %d", stats[abba], 1)
	}

	if _, ok := countries["UKR"]; !ok {
		t.Error("UKR was not resolved by MSISDN")
	}

	msisdn2 := "310213243546"
	gc := "Gigliola Cinquetti"

	svc.RegisterVote(msisdn2, gc)

	if stats[gc] != 1 {
		t.Error("Score for Gigliola Cinquetti is %d, expected %d", stats[gc], 1)
	}

	if _, ok := countries["NLD"]; !ok {
		t.Error("NDL was not resolved by MSISDN")
	}
}

type SkoreKprMock struct {
	AddPointFunc         func(key string) error
	AddCountryFunc       func(name string) error
	GetAllCandidatesFunc func() ([]string, error)
	GetAllCountriesFunc  func() ([]string, error)
	GetFunc              func(key string) (int, error)
}

func (sk *SkoreKprMock) GetAllCandidates() ([]string, error) {
	return sk.GetAllCandidatesFunc()
}

func (sk *SkoreKprMock) GetAllCountries() ([]string, error) {
	return sk.GetAllCountriesFunc()
}

func (sk *SkoreKprMock) Get(key string) (int, error) {
	return sk.GetFunc(key)
}

func (sk *SkoreKprMock) AddPoint(key string) error {
	return sk.AddPointFunc(key)
}

func (sk *SkoreKprMock) AddCountry(code string) error {
	return sk.AddCountryFunc(code)
}

type EnquirerMock struct{}

func (e *EnquirerMock) Lookup(msisdn string) (string, error) {
	var country string

	switch msisdn[:3] {
	case "380":
		country = "UKR"
	case "310":
		country = "NLD"
	default:
		return "", errors.New("lookup failed")
	}

	return country, nil
}

type MessengerMock struct {
	RequestSMSFunc func(originator, recipient, text string)
}

func (mm *MessengerMock) RequestSMS(originator, recipient, text string) {
	mm.RequestSMSFunc(originator, recipient, text)
}
