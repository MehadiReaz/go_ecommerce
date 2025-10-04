package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BkashClient struct {
	appKey    string
	appSecret string
	username  string
	password  string
	baseURL   string
}

// NewBkashClient creates a new bKash client
func NewBkashClient(appKey, appSecret, username, password string) *BkashClient {
	return &BkashClient{
		appKey:    appKey,
		appSecret: appSecret,
		username:  username,
		password:  password,
		baseURL:   "https://tokenized.sandbox.bka.sh/v1.2.0-beta", // Sandbox URL
	}
}

// GrantToken requests a grant token from bKash
func (c *BkashClient) GrantToken() (string, error) {
	url := fmt.Sprintf("%s/tokenized/checkout/token/grant", c.baseURL)

	payload := map[string]string{
		"app_key":    c.appKey,
		"app_secret": c.appSecret,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("username", c.username)
	req.Header.Set("password", c.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if token, ok := result["id_token"].(string); ok {
		return token, nil
	}

	return "", fmt.Errorf("failed to get token from bKash")
}

// CreatePayment creates a bKash payment
func (c *BkashClient) CreatePayment(amount float64, invoiceNumber string) (map[string]interface{}, error) {
	token, err := c.GrantToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/tokenized/checkout/create", c.baseURL)

	payload := map[string]interface{}{
		"mode":                  "0011",
		"payerReference":        " ",
		"callbackURL":           "https://example.com/callback",
		"amount":                amount,
		"currency":              "BDT",
		"intent":                "sale",
		"merchantInvoiceNumber": invoiceNumber,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Set("X-APP-Key", c.appKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ExecutePayment executes a bKash payment
func (c *BkashClient) ExecutePayment(paymentID string) (map[string]interface{}, error) {
	token, err := c.GrantToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/tokenized/checkout/execute", c.baseURL)

	payload := map[string]string{
		"paymentID": paymentID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Set("X-APP-Key", c.appKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
