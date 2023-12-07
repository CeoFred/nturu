package provider

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nturu/microservice-template/sendgrid"
)

type Provider struct {
	BaseURL    string
	HTTPClient *http.Client
	MailClient *sendgrid.Client
}

func (p *Provider) init() {
	if p.HTTPClient == nil {
		p.HTTPClient = &http.Client{}
	}
}

func NewHttpProvider(baseURL string) *Provider {
	return &Provider{
		BaseURL: baseURL,
	}
}

func NewMailProvider(apiKey string) *Provider {
	return &Provider{
		MailClient: &sendgrid.Client{},
	}
}

func (p *Provider) Post(path string, payload interface{}, headers map[string]string) (*http.Response, error) {
	p.init()
	url := p.BaseURL + path

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create the request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *Provider) Get(path string) (*http.Response, error) {
	p.init()
	url := p.BaseURL + path
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
