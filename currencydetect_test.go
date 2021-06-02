package gocash

import (
	"encoding/json"
	"testing"
)

func TestCurrencyDetectHasCode(t *testing.T) {
	cur := ParseCurrencyType("1,000.00 USD")
	if cur.Code != "USD" {
		t.Logf("Got wrong type. Expecting USD, got : %s", cur.Code)
	}

	cur = ParseCurrencyType("1.000,00 CHF")
	if cur.Code != "CHF" {
		t.Logf("Got wrong type. Expecting USD, got : %s", cur.Code)
	}
}

func TestParseUSDEasy(t *testing.T){
	mon := []byte(`{"Dollar": "$5.34",
					"Hello": "Hello"}`)
	type TmpStruct struct{
		Mon Money `json:"Dollar"`
		H string `json:"Hello"`
	}
	var m TmpStruct
	err := json.Unmarshal(mon, &m)
	if err != nil{
		t.Error(err)
	}
	t.Log(m.Mon.String())
	/* if m.Currency != "USD"{

	} */
	t.Log(m.Mon.String())
	if m.Mon.Dollar != 5{
		t.Error("Dollar incorrect")
	}
	if m.Mon.Decimal.D != "34"{
		t.Error("Cents wrong")
	}
}
