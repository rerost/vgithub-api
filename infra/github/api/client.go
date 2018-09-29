package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Client is client for github v4. See https://developer.github.com/v4/
type Client interface {
	Request(query Query) (*http.Response, error)
}

// Query is GraphQL query
type Query string

const queryFmt = `{"query": "query %s"}`

func (query *Query) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(queryFmt, query)), nil
}

func (query *Query) UnmarshalJSON(b []byte) error {
	var q string
	fmt.Scanf(string(*query), queryFmt, &q)
	*query = Query(q)
	return nil
}

type clientImp struct {
	httpClient *http.Client
	graphqlURL *url.URL
	token      string
}

// NewClient is create client for github v4
func NewClient(httpClient *http.Client, graphqlURL *url.URL, token string) (Client, error) {
	if token == "" {
		return nil, errors.New("Need github token")
	}

	return &clientImp{
		httpClient: httpClient,
		graphqlURL: graphqlURL,
		token:      token,
	}, nil
}

func (api *clientImp) Request(query Query) (*http.Response, error) {
	requestBody, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("Query Marshaling error: %s", err)
	}
	request, err := http.NewRequest(http.MethodGet, api.graphqlURL.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("bearer %s", api.token))

	return api.httpClient.Do(request)
}
