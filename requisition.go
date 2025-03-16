package main

import (
	"fmt"
	"time"
)

type Requisition struct {
	ID                string    `json:"id"`
	Created           time.Time `json:"created"`
	Redirect          string    `json:"redirect"`
	Status            string    `json:"status"`
	InstitutionID     string    `json:"institution_id"`
	Agreement         string    `json:"agreement"`
	Reference         string    `json:"reference"`
	Accounts          []string  `json:"accounts"`
	UserLanguage      string    `json:"user_language"`
	Link              string    `json:"link"`
	SSN               string    `json:"ssn"`
	AccountSelection  bool      `json:"account_selection"`
	RedirectImmediate bool      `json:"redirect_immediate"`
}

func (gc *GoCardlessClient) CreateRequisition(institution string, agreementID string) (Requisition, error) {
	var requisition Requisition

	requestBody := map[string]string{
		"redirect":       "http://localhost:3000",
		"institution_id": institution,
		"user_language":  "EN",
		"agreement":      agreementID,
	}

	resp, err := gc.R().
		SetBody(requestBody).
		SetResult(&requisition).
		Post("requisitions/")

	if err != nil {
		return requisition, err
	}

	if resp.IsError() {
		return requisition, fmt.Errorf("API error: %s", resp.String())
	}

	return requisition, nil
}

func (gc *GoCardlessClient) GetRequisition(id string) (Requisition, error) {
	var requisition Requisition

	resp, err := gc.R().
		SetResult(&requisition).
		Get("requisitions/" + id + "/")

	if err != nil {
		return requisition, err
	}

	if resp.IsError() {
		return requisition, fmt.Errorf("API error: %s", resp.String())
	}

	return requisition, nil
}
