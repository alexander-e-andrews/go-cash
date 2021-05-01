package gocash

import (
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
