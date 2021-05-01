package gocash

import (
	"testing"
)

func TestAddTwoPositives(t *testing.T) {
	a, _ := MakeAMoney("5", "75", "USD")
	b, _ := MakeAMoney("4", "30", "USD")
	c, _ := a.Add(b)
	t.Logf("%s + %s = %s", a.String(), b.String(), c.String())
}

func TestSubTwoPositivesOutPositive(t *testing.T) {
	a, _ := MakeAMoney("5", "75", "USD")
	b, _ := MakeAMoney("4", "30", "USD")
	c, _ := a.Sub(b)
	t.Logf("%s + %s = %s", a.String(), b.String(), c.String())
}
func TestAddPositivesNegativeOutNegative(t *testing.T) {
	b, _ := MakeAMoney("5", "75", "USD")
	a, _ := MakeAMoney("4", "30", "USD")
	b.Negative = true
	c, _ := a.Add(b)
	t.Logf("%s + %s = %s", a.String(), b.String(), c.String())
}

func TestAddPositivesNegativeOutPositive(t *testing.T) {
	a, _ := MakeAMoney("5", "75", "USD")
	b, _ := MakeAMoney("4", "30", "USD")
	b.Negative = true
	c, _ := a.Add(b)
	t.Logf("%s + %s = %s", a.String(), b.String(), c.String())
}

func TestSubTwoPositivesOutNegative(t *testing.T) {
	b, _ := MakeAMoney("5", "75", "USD")
	a, _ := MakeAMoney("4", "30", "USD")
	c, _ := a.Sub(b)
	t.Logf("%s + %s = %s", a.String(), b.String(), c.String())
}

func TestSubPositiveNegative(t *testing.T) {
	b, _ := MakeAMoney("5", "75", "USD")
	a, _ := MakeAMoney("4", "30", "USD")
	b.Negative = true
	c, _ := a.Sub(b)
	t.Logf("%s + %s = %s", a.String(), b.String(), c.String())
}


func TestStringDecimalAdd(t *testing.T){
	a := Decimal{"50"}
	b := Decimal{"25"}
	c, overflow, err := a.Add(b)
	if err != nil{
		t.Log("Overflowed the binaires")
	}
	t.Logf("%s + %s = %s overflow: %t", a, b, c, overflow)

	a = Decimal{"50"}
	b = Decimal{"75"}
	c, overflow, err = a.Add(b)
	if err != nil{
		t.Log("Overflowed the binaires")
	}
	t.Logf("%s + %s = %s overflow: %t", a, b, c, overflow)

	a = Decimal{"75"}
	b = Decimal{"75"}
	c, overflow, err = a.Add(b)
	if err != nil{
		t.Log("Overflowed the binaires")
	}
	t.Logf("%s + %s = %s overflow: %t", a, b, c, overflow)

	a = Decimal{"001"}
	b = Decimal{"75"}
	c, overflow, err = a.Add(b)
	if err != nil{
		t.Log("Overflowed the binaires")
	}
	t.Logf("%s + %s = %s overflow: %t", a, b, c, overflow)
}

func TestStringDecimalSub(t *testing.T){
	a := Decimal{"50"}
	b := Decimal{"25"}
	c, underflow, err := a.Sub(b)
	if err != nil{
		t.Log("underflowed the binaires")
	}
	t.Logf("%s - %s = %s underflow: %t", a, b, c, underflow)

	a = Decimal{"50"}
	b = Decimal{"75"}
	c, underflow, err = a.Sub(b)
	if err != nil{
		t.Log("underflowed the binaires")
	}
	t.Logf("%s - %s = %s underflow: %t", a, b, c, underflow)

	a = Decimal{"75"}
	b = Decimal{"75"}
	c, underflow, err = a.Sub(b)
	if err != nil{
		t.Log("underflowed the binaires")
	}
	t.Logf("%s - %s = %s underflow: %t", a, b, c, underflow)

	a = Decimal{"001"}
	b = Decimal{"75"}
	c, underflow, err = a.Sub(b)
	if err != nil{
		t.Log("underflowed the binaires")
	}
	t.Logf("%s - %s = %s underflow: %t", a, b, c, underflow)
}