package client

import (
	"github.com/go-resty/resty/v2"
)

const momentum = "http://localhost:8000"

func Init() *resty.Client {
	return resty.New().
		SetBaseURL(momentum).
		SetHeader("Content-Type", "application/json")
}

func Get(endpoint string) (*resty.Response, error) {
	return Init().R().Get(endpoint)
}

func Post(endpoint string, body interface{}) (*resty.Response, error) {
	return Init().R().SetBody(body).Post(endpoint)
}
