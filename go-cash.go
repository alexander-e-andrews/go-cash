package gocash

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

//Just because we use only two decimals usually, gas is traded to the 3 decimals
//Money is a dollar amount and information on how to store it
type Money struct {
	Dollar   uint64  //The point before your decimal
	Decimal  Decimal //the point after decimal These both need new names
	Currency Currency
	Negative bool //Unused
}

type Decimal struct {
	D string //I was hoping to be able to just use "string" and then have Decimal just act as an identifier, but that didnt give access to the string
}

type JSONParseError struct {
	AttemptedParse string
}

func (e *JSONParseError) Error() string {
	return fmt.Sprintf("Something went wrong trying to parse: %s", e.AttemptedParse)
}

type OverflowError struct{}

func (e *OverflowError) Error() string {
	return "Overflow Error"
}

type WrongCurrencyError struct{}

func (e *WrongCurrencyError) Error() string {
	return "Wrong Currency"
}

type StringParseError struct {
	parsedString string
}

func (e *StringParseError) Error() string {
	return fmt.Sprintf("Could not parse %s into uint64", e.parsedString)
}

func (m Money) String() string {
	//fmt.Println(bits.TrailingZeros64(m.Fractional))
	if m.Negative {
		return fmt.Sprintf("-%s%d%s%s", m.Currency.Symbol, m.Dollar, ".", m.Decimal.D)
	}
	return fmt.Sprintf("%s%d%s%s", m.Currency.Symbol, m.Dollar, ".", m.Decimal.D)
}

//The assumed type that unmarshalled moneys will be. Set this to "" to have the unmarshaller try to determine monetary type
//If there are multiple currencies using the same monetary format, there is no guarantee the correct one will be used
//Unmarshaller will first look for all the monetary Codes before moving on to format match
var UnmarshalledType string = "USD"

//Add this money to another money, will return an error if they are not the same currency
//Will throw an error if whole amounts overflow
func (m Money) Add(second Money) (value Money, err error) {
	return hiddenAddTwoMonies(m, second)
}

//These add and subtract functions all need to be cleaned up to properly handle negative numbers
func (m Money) Sub(second Money) (value Money, err error) {
	second.Negative = !second.Negative
	return hiddenAddTwoMonies(m, second)
}

//When a user says subtract, actually just make the number !negative
func hiddenAddTwoMonies(a Money, b Money) (c Money, err error) {
	//we could short circuit if one of these numbers is 0, but going to skip for now
	var overflow bool
	var underflow bool
	_ = underflow
	if !a.Negative && !b.Negative {
		c.Decimal, overflow, err = a.Decimal.Add(b.Decimal)
		//c.Fractional, overflow = uint64OverflowAdd(a.Fractional, b.Fractional)
		//fmt.Println(c.Fractional >> bits.TrailingZeros64(c.Fractional))
		if overflow {
			a.Dollar, overflow = uint64OverflowAdd(a.Dollar, 1)
		}
		if overflow {
			err = &OverflowError{}
		}
		c.Dollar, overflow = uint64OverflowAdd(a.Dollar, b.Dollar)
		if overflow {
			err = &OverflowError{}
		}
		return c, err
	} else if a.Negative && b.Negative {
		//Since both number are negative, its the same as adding them, just need to set c as a negative value
		c.Negative = true
		c.Decimal, overflow, err = a.Decimal.Add(b.Decimal)
		if overflow {
			a.Dollar, overflow = uint64OverflowAdd(a.Dollar, 1)
		}
		if overflow {
			err = &OverflowError{}
		}
		c.Dollar, overflow = uint64OverflowAdd(a.Dollar, b.Dollar)
		if overflow {
			err = &OverflowError{}
		}
		return c, err
	} else if !a.Negative && b.Negative {
		//a - b
		var aLessThanb bool
		if a.Dollar == b.Dollar {
			aLessThanb, err = b.Decimal.GreaterThanEq(a.Decimal)
		} else {
			aLessThanb = a.Dollar < b.Dollar
		}
		//If a is less than b, do the number swap style of subtraction, and set c.Negative to true
		if aLessThanb {
			c, underflow = subtractAMoney(b, a)
			c.Negative = true
			return c, nil
		} else {
			c, underflow = subtractAMoney(a, b)
			return c, nil
		}
	} else if a.Negative && !b.Negative {
		//-a + b
		//do b - a
		var bLessThana bool
		if b.Dollar == a.Dollar {
			bLessThana, err = a.Decimal.GreaterThanEq(b.Decimal)
		} else {
			bLessThana = b.Dollar < a.Dollar
		}
		//If a is less than b, do the number swap style of subtraction, and set c.Negative to true
		if bLessThana {
			c, underflow = subtractAMoney(a, b)
			c.Negative = true
			return c, nil
		} else {
			c, underflow = subtractAMoney(b, a)
			return c, nil
		}
	}

	return
}

//a is positive, b is negative. The operation does  +a-b. Broke it out for more control
//Requires a to be greater than b
func subtractAMoney(a Money, b Money) (c Money, underflow bool) {
	var err error
	c.Decimal, underflow, err = a.Decimal.Sub(b.Decimal)
	if err != nil {
		fmt.Println(err)
	}
	if underflow {
		a.Dollar, underflow = uint64UnderflowSub(a.Dollar, 1)
	}
	c.Dollar, underflow = uint64UnderflowSub(a.Dollar, b.Dollar)

	return
}

func uint64OverflowAdd(a uint64, b uint64) (c uint64, overflow bool) {
	c = a + b
	if c < a || c < b {
		return c, true
	}
	return c, false
}

func uint64UnderflowSub(a uint64, b uint64) (c uint64, underflow bool) {
	var isSecondCentsBigger bool //If the second cents is bigger than our first cents
	if b > a {
		isSecondCentsBigger = true
	}

	if isSecondCentsBigger {
		c = b - a
		underflow = true
		//Create 10^x to subtract c from to act as a carry
		word := strconv.FormatUint(c, 10)
		bigSub := uint64(math.Pow(10, float64(len(word))))
		c = bigSub - c
	} else {
		c = a - b
	}

	return
}

//UnmarshalJSON will have to be able to determine the type of the value, unless we set it not to
func (m *Money) UnmarshalJSON(bytes []byte)(err error) {
	if len(bytes) < 1 {
		return &JSONParseError{AttemptedParse: string(bytes)}
	}
	//fmt.Println(string(bytes))

	*m, err = ParseString(string(bytes))
	//fmt.Println(m.String())
	//if unicode.IsDigit(s)
	return err
}

//ParseString takes in a monetary string in the form of $4.20 and turns it into a money. If there is any issues, it will return an err
func ParseString(str string)(m Money, err error){
	s := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(str), `"`), `"`)
	//just a temporary patch to get the money working
	s = strings.Replace(s, "$", "", -1)
	moneyArray := strings.Split(s, ".")
	//TODO Change make error reporting a real thing
	if len(moneyArray) < 2 {
		m, _ = MakeAMoney(moneyArray[0], "00", "USD")
		//return &JSONParseError{AttemptedParse: string(bytes)}
	} else {
		m, _ = MakeAMoney(moneyArray[0], moneyArray[1], "USD")
	}

	return
}

func (m Money) MarshalJSON() ([]byte, error) {
	//IN a weird spot. Since we made an unmarshal to take a simple input, we now need to do a complete reparcing for when we unmarshal from a marshal
	//Could create a two part branch. If we fail to unmarshal how to struct is normaly, then we use our custom unmarshal. And maybe the user
	//can set which type of marshalling they would like
	return []byte("\"" + m.String() + "\""), nil
}

func MakeAMoney(dollar string, fractional string, code string) (m Money, err error) {
	dollar = strings.ReplaceAll(dollar, ",", "")
	m.Dollar, err = strconv.ParseUint(dollar, 10, 64)

	m.Decimal.D = fractional
	m.Currency = ParseCurrencyType(code)

	return
}

func (d Decimal) Add(b Decimal) (c Decimal, overflow bool, err error) {
	maxLength := 0
	if len(d.D) > len(b.D) {
		maxLength = len(d.D)
	} else {
		maxLength = len(b.D)
	}
	d.D = d.D + strings.Repeat("0", maxLength-len(d.D))
	b.D = b.D + strings.Repeat("0", maxLength-len(b.D))

	di, err := strconv.ParseUint(d.D, 10, 64)
	if err != nil {
		return c, overflow, &StringParseError{d.D}
	}
	bi, err := strconv.ParseUint(b.D, 10, 64)
	if err != nil {
		return c, overflow, &StringParseError{b.D}
	}
	//Need to check for normal unuint64 overflow here
	f, overflow := uint64OverflowAdd(di, bi)
	if overflow {
		return c, overflow, &OverflowError{}
	}

	c.D = strconv.FormatUint(f, 10)
	if len(c.D) > maxLength {
		c.D = strings.TrimPrefix(c.D, "1")
		overflow = true
	}
	c = c.Trim()
	return
}

//Underflow means we have gone negative
func (d Decimal) Sub(b Decimal) (c Decimal, underflow bool, err error) {
	maxLength := 0
	if len(d.D) > len(b.D) {
		maxLength = len(d.D)
	} else {
		maxLength = len(b.D)
	}
	d.D = d.D + strings.Repeat("0", maxLength-len(d.D))
	b.D = b.D + strings.Repeat("0", maxLength-len(b.D))

	di, err := strconv.ParseUint(d.D, 10, 64)
	if err != nil {
		return c, underflow, &StringParseError{d.D}
	}
	bi, err := strconv.ParseUint(b.D, 10, 64)
	if err != nil {
		return c, underflow, &StringParseError{b.D}
	}
	//Need to check for normal unuint64 overflow here
	f, underflow := uint64UnderflowSub(di, bi)
	/* if underflow{
		return c, underflow, &OverflowError{}
	} */

	c.D = strconv.FormatUint(f, 10)
	if underflow {
		//c.D = "0" + c.D
		underflow = true
	}
	c = c.Trim()
	return
}

func (d Decimal) GreaterThanEq(b Decimal) (greater bool, err error) {
	maxLength := 0
	if len(d.D) > len(b.D) {
		maxLength = len(d.D)
	} else {
		maxLength = len(b.D)
	}
	d.D = d.D + strings.Repeat("0", maxLength-len(d.D))
	b.D = b.D + strings.Repeat("0", maxLength-len(b.D))

	di, err := strconv.ParseUint(d.D, 10, 64)
	if err != nil {
		return false, &StringParseError{d.D}
	}
	bi, err := strconv.ParseUint(b.D, 10, 64)
	if err != nil {
		return false, &StringParseError{b.D}
	}

	return di >= bi, nil
}

func (d Decimal) Trim() Decimal {
	return Decimal{D: strings.TrimRight(d.D, "0")}
}

func (d Money) IsZero() bool {
	if d.Dollar == 0 && d.Decimal.D == "" {
		return true
	}
	return false
}
