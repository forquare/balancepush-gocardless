package main

import (
	"fmt"
	"time"
)

type AgreementData struct {
	ID                 string    `json:"id"`
	Created            time.Time `json:"created"`
	MaxHistoricalDays  int       `json:"max_historical_days"`
	AccessValidForDays int       `json:"access_valid_for_days"`
	AccessScope        []string  `json:"access_scope"`
	Accepted           string    `json:"accepted"`
	InstitutionID      string    `json:"institution_id"`
}

func (gc *GoCardlessClient) GetAgreement(institutionID string) (string, error) {
	var agreement AgreementData

	requestBody := map[string]string{
		"institution_id":        institutionID,
		"max_historical_days":   "180",
		"access_valid_for_days": "90",
	}

	resp, err := gc.R().
		SetResult(&agreement).
		SetBody(requestBody).
		Post("/agreements/enduser/")

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("API error: %s", resp.String())
	}

	return agreement.ID, err
}
