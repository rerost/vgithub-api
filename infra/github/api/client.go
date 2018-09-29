package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Client is client for github v4. See https://developer.github.com/v4/
type Client interface {
	Request(query string) (*http.Response, error)
}

type clientImp struct {
	httpClient *http.Client
	graphqlURL *url.URL
	token      string
}

// NewClient is create client for github v4
func NewClient(httpClient *http.Client, graphqlURL *url.URL, token string) (Client, error) {
	if token == "" {
		return nil, errors.New("You neeed github token")
	}

	return &clientImp{
		httpClient: httpClient,
		graphqlURL: graphqlURL,
		token:      token,
	}, nil
}

func (api *clientImp) Request(query string) (*http.Response, error) {
	requestBody := []byte(query)
	request, err := http.NewRequest(http.MethodGet, api.graphqlURL.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("bearer %s", api.token))

	return api.httpClient.Do(request)
}
