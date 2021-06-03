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

func TestParseMicroAmount(t *testing.T){
	m, err := ParseString("$1.008")
	if err != nil{
		t.Error(err)
	}

	if m.String() != "$1.008"{
		t.Log(m.String())
		t.Error("Did not parse right")
	}
}

func TestParseZero(t *testing.T){
	m, err := ParseString("$0.00")
	if err != nil{
		t.Error(err)
	}

	if m.String() != "$0.00"{
		t.Log(m.String())
		t.Error("Did not parse right")
	}
}

func TestMarshall(t *testing.T){
	m, _ := ParseString("$12.34")
	b, err := json.Marshal(m)
	if err != nil{
		t.Error(err)
	}

	t.Log(string(b))
}
