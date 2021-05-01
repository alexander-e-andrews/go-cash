package gocash

import (
	"encoding/json"
	"os"

	"github.com/alexander-e-andrews/GoRegexTree"
	"github.com/dlclark/regexp2"
)

type Currency struct {
	FullName string
	Code     string
	Symbol   string
	Format   Format //May change this to a string, undecided
	//Could also include the country code and name from the iso, but ehhhh
}

type Format struct {
	SymbolLeft bool   //If the symbol is on the left or right side
	Thousands  string //The symbol to separate thousands USD: ","
	Decimal    string //Symbol for decimal split USD: "."
}

func ParseCurrencyType(value string) (cur Currency) {
	match, err := currencyCatchRegex.FindStringMatch(value)
	pError(err)

	if match != nil && len(match.Groups()) > 1 {
		currencyString := match.Groups()[1].String()
		return *currencyListByCode[currencyString]
	}
	//Else we need to check the symbol and format to determine what the type is. Going to ignore for now, and just return USD so I can get back to my main project

	return *currencyListByCode["USD"]
}

var currencyList []Currency
var currencyListByCode map[string]*Currency

//var currencyCodeCatch *goregextree.Node  My package doesn't implement the trie, so instead I will use the regex
var currencyCatchRegex *regexp2.Regexp

func init() {
	currencyList = make([]Currency, 0)
	currencyListJson, err := os.Open("currencies.json")
	pError(err)
	decoder := json.NewDecoder(currencyListJson)
	err = decoder.Decode(&currencyList)
	pError(err)

	currencyCodeCatch := goregextree.CreateSearchTree()
	currencyListByCode = make(map[string]*Currency)
	for x := range currencyList {
		currencyCodeCatch.AddWordString(currencyList[x].Code)
		currencyListByCode[currencyList[x].Code] = &currencyList[x]
	}

	currencyCatchRegex = currencyCodeCatch.BuildRegex(true, []rune{}, []rune{}, false)
}

//Example money formats
//$1,234.89 USD
//fr. 1.234,56  Swiss Franc

//Going to build the struct that determines how we parse a value
func parseTree() {

}

func pError(err error) {
	if err != nil {
		panic(err)
	}
}
