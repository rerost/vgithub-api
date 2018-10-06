package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Client is client for github v4. See https://developer.github.com/v4/
type Client interface {
	Request(query Query) (*http.Response, error)
}

// Query is GraphQL query
type Query string

const queryFmt = `{"query": "query %s"}`

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
	wrappedRequestBody := []byte(fmt.Sprintf(queryFmt, toJSONString(string(query))))
	request, err := http.NewRequest(http.MethodPost, api.graphqlURL.String(), bytes.NewBuffer(wrappedRequestBody))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("bearer %s", api.token))

	return api.httpClient.Do(request)
}

func toJSONString(s string) string {
	tmp := strings.Replace(s, "\n", "", -1)
	tmp = strings.Replace(tmp, "\t", "", -1)
	tmp = strings.Replace(tmp, "\"", "\\\"", -1)
	return tmp
}
