package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"sync"
	"time"
)

type GoCardlessClient struct {
	client                 *resty.Client
	secretID               string
	secretKey              string
	accessToken            string
	refreshToken           string
	accessTokenExpiryTime  time.Time
	refreshTokenExpiryTime time.Time
	tokenMutex             sync.Mutex
}

// NewGoCardlessClient initializes the API client with authentication
func NewGoCardlessClient(secretID, secretKey string) (*GoCardlessClient, error) {
	client := resty.New().
		SetBaseURL("https://bankaccountdata.gocardless.com/api/v2").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetTimeout(10 * time.Second)

	gc := &GoCardlessClient{
		client:    client,
		secretID:  secretID,
		secretKey: secretKey,
	}

	err := gc.getToken()
	if err != nil {
		return nil, err
	}

	return gc, err
}

func (gc *GoCardlessClient) R() *resty.Request {
	if err := gc.ensureValidAccessToken(); err != nil {
		log.Fatalf("failed to refresh access token: %v", err)
	}
	return gc.client.R().SetAuthToken(gc.accessToken)
}

func (gc *GoCardlessClient) ensureValidAccessToken() error {
	gc.tokenMutex.Lock()
	defer gc.tokenMutex.Unlock()

	if time.Now().Before(gc.accessTokenExpiryTime) {
		// Token is still valid
		return nil
	}

	if time.Now().After(gc.refreshTokenExpiryTime) {
		// Refresh token expired, get a new access token
		gc.accessToken = ""
	}

	if gc.accessToken == "" {
		err := gc.getToken()
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		return nil
	}

	// If we get here, the access token is invalid, but the refresh token is valid
	// refresh the access token
	err := gc.refresh()
	if err != nil {
		return err
	}

	return nil
}

func (gc *GoCardlessClient) getToken() error {
	if gc.accessToken != "" {
		return nil
	}

	requestBody := map[string]string{
		"secret_id":  gc.secretID,
		"secret_key": gc.secretKey,
	}

	var response struct {
		Access         string `json:"access"`
		Expires        int    `json:"access_expires"`
		Refresh        string `json:"refresh"`
		RefreshExpires int    `json:"refresh_expires"`
	}

	resp, err := gc.client.R().
		SetBody(requestBody).
		SetResult(&response).
		Post("/token/new/")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("token request failed: %s", resp.String())
	}

	gc.accessToken = response.Access
	gc.refreshToken = response.Refresh
	gc.accessTokenExpiryTime = time.Now().Add(time.Duration(response.Expires) * time.Second)
	gc.refreshTokenExpiryTime = time.Now().Add(time.Duration(response.RefreshExpires) * time.Second)

	return nil
}

func (gc *GoCardlessClient) refresh() error {
	if gc.refreshToken == "" {
		return fmt.Errorf("refresh token not set")
	}

	requestBody := map[string]string{
		"refresh": gc.refreshToken,
	}

	var response struct {
		Access  string `json:"access"`
		Expires int    `json:"access_expires"`
	}

	resp, err := gc.client.R().
		SetBody(requestBody).
		SetResult(&response).
		Post("/token/refresh")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("token refresh failed: %s", resp.String())
	}

	gc.accessToken = response.Access
	gc.accessTokenExpiryTime = time.Now().Add(time.Duration(response.Expires) * time.Second)

	return nil
}
