package main

import (
	"fmt"
	"strconv"
)

type BalanceAmount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type Balance struct {
	BalanceAmount BalanceAmount `json:"balanceAmount"`
	BalanceType   string        `json:"balanceType"`
	ReferenceDate string        `json:"referenceDate"`
}

type BalanceData struct {
	Balances []Balance `json:"balances"`
}

type AccountDetails struct {
	ResourceID      string `json:"resourceId"`
	IBAN            string `json:"iban"`
	Currency        string `json:"currency"`
	OwnerName       string `json:"ownerName"`
	Name            string `json:"name"`
	Product         string `json:"product"`
	CashAccountType string `json:"cashAccountType"`
}

type AccountData struct {
	Account AccountDetails `json:"account"`
}

var currencySymbols = map[string]string{
	"USD": "$",
	"EUR": "€",
	"GBP": "£",
	"JPY": "¥",
	"AUD": "A$",
	"CAD": "C$",
	"CHF": "CHF",
	"CNY": "¥",
	"SEK": "kr",
	"NZD": "NZ$",
}

func (gc *GoCardlessClient) GetAccountBalance(accountID string, balanceType string) (float64, string, string, error) {
	result := BalanceData{}
	balance := 0.0
	currency := ""
	currencySymbol := ""

	resp, err := gc.R().
		SetResult(&result).
		Get(fmt.Sprintf("/accounts/%s/balances/", accountID))

	if err != nil {
		return 0, "", "", err
	}

	if resp.IsError() {
		return 0, "", "", fmt.Errorf("API error:\n%s\n", resp.String())
	}

	for _, b := range result.Balances {
		// https://gist.github.com/amilos/1ce55dbdfa336eee1de74d3e800496c1#file-bg-psd2-yaml-L3699
		if b.BalanceType == balanceType {
			balance, _ = strconv.ParseFloat(b.BalanceAmount.Amount, 64)
			currencySymbol = getCurrencySymbol(b.BalanceAmount.Currency)
			currency = b.BalanceAmount.Currency
			break
		}
	}

	/*
		jsonPretty, err := json.MarshalIndent(result, "", "  ")
		fmt.Println("Account ID: ", accountID)
		fmt.Println("Balance Details: ", string(jsonPretty))

	*/

	return balance, currency, currencySymbol, nil
}

func getCurrencySymbol(code string) string {
	if symbol, found := currencySymbols[code]; found {
		return symbol
	}
	return code // Return the code itself if symbol not found
}

/*
49.42
Account ID:  1193a86d-4f8a-4214-9873-143bb2dfa03d
Balance Details:  {
  "balances": [
    {
      "balanceAmount": {
        "amount": "73.74",
        "currency": "GBP"
      },
      "balanceType": "interimBooked",
      "referenceDate": "2025-02-24"
    },
    {
      "balanceAmount": {
        "amount": "49.42",
        "currency": "GBP"
      },
      "balanceType": "interimAvailable",
      "referenceDate": "2025-02-24"
    }
  ]
}
Account Details:  {
  "account": {
    "resourceId": "588f9b67-7ad7-38d6-ab94-278875302c5d",
    "iban": "",
    "currency": "GBP",
    "ownerName": "Mr Benjamin Lavery-Griffiths",
    "name": "",
    "product": "",
    "cashAccountType": "CACC"
  }
}
: £49.42
2255.64
Account ID:  6522408d-4546-4886-bcb5-f562f302e3ea
Balance Details:  {
  "balances": [
    {
      "balanceAmount": {
        "amount": "2255.64",
        "currency": "GBP"
      },
      "balanceType": "interimAvailable",
      "referenceDate": "2025-02-24"
    },
    {
      "balanceAmount": {
        "amount": "-144.36",
        "currency": "GBP"
      },
      "balanceType": "interimBooked",
      "referenceDate": "2025-02-24"
    }
  ]
}
Account Details:  {
  "account": {
    "resourceId": "ad07f018-8bc0-3584-9d81-544221dbe436",
    "iban": "",
    "currency": "GBP",
    "ownerName": "Mr Benjamin Lavery-Griffiths",
    "name": "",
    "product": "",
    "cashAccountType": "CARD"
  }
}
: £2255.64
204.20
Account ID:  f7ac8973-0499-4421-bc46-a49827103892
Balance Details:  {
  "balances": [
    {
      "balanceAmount": {
        "amount": "204.20",
        "currency": "GBP"
      },
      "balanceType": "interimBooked",
      "referenceDate": "2025-02-24"
    },
    {
      "balanceAmount": {
        "amount": "204.20",
        "currency": "GBP"
      },
      "balanceType": "interimAvailable",
      "referenceDate": "2025-02-24"
    }
  ]
}
Account Details:  {
  "account": {
    "resourceId": "6b64f1e9-761a-3fdb-b38d-834ad8c6c90e",
    "iban": "",
    "currency": "GBP",
    "ownerName": "Mr Benjamin Lavery-Griffiths",
    "name": "Joint Current Account",
    "product": "",
    "cashAccountType": "CACC"
  }
}
Joint Current Account: £204.20
*/
