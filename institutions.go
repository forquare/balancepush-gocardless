package main

import "fmt"

type Institution struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	BIC                   string   `json:"bic"`
	TransactionTotalDays  string   `json:"transaction_total_days"`
	Countries             []string `json:"countries"`
	Logo                  string   `json:"logo"`
	IdentificationCodes   []string `json:"identification_codes"`
	MaxAccessValidForDays string   `json:"max_access_valid_for_days"`
}

func (gc *GoCardlessClient) GetInstitutions() ([]Institution, error) {
	var institutions []Institution

	resp, err := gc.R().
		SetResult(&institutions).
		Get("/institutions/?country=gb")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	return institutions, nil
}
