package api_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/rerost/vgithub-api/infra/github/api"

	"net/http/httptest"
)

func newDummyServer(callback func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	dummyServerHandler := http.HandlerFunc(callback)
	dummyServer := httptest.NewServer(dummyServerHandler)
	return dummyServer
}

func TestNewClient(t *testing.T) {
	inOutPairs := []struct {
		test string
		in   string // Token
		out  string // Error message
	}{
		{
			test: "When passed token",
			in:   "TESTTOKEN",
			out:  "",
		},
		{
			test: "When not passed token",
			in:   "",
			out:  "Need github token",
		},
	}

	for _, p := range inOutPairs {
		t.Run(p.test, func(t *testing.T) {
			dummyServer := newDummyServer(func(w http.ResponseWriter, r *http.Request) {})
			dummyURL, _ := url.Parse(dummyServer.URL)
			_, err := api.NewClient(dummyServer.Client(), dummyURL, p.in)
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			} else if p.out != errMsg {
				t.Errorf(`Want: %s\nHave: %s\n`, p.out, err.Error())
			}
		})
	}
}

func TestRequest(t *testing.T) {
	dummyToken := "TESTTOKEN"

	inOutPairs := []struct {
		test  string
		query string // DummyResponse (== Query)
	}{
		{
			test: "Empty query",
		},
	}
	var wg sync.WaitGroup
	for _, p := range inOutPairs {
		wg.Add(1)
		t.Run(p.test, func(t *testing.T) {
			defer wg.Done()
			dummyServer := newDummyServer(func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				b, _ := ioutil.ReadAll(r.Body)
				w.Write([]byte(b))
			})
			dummyURL, _ := url.Parse(dummyServer.URL)
			client, _ := api.NewClient(dummyServer.Client(), dummyURL, dummyToken)
			res, err := client.Request(p.query)
			if err != nil {
				t.Errorf("Unexpected error is occured: %s", err)
				return
			}
			body, _ := ioutil.ReadAll(res.Body)
			if string(body) != p.query {
				t.Errorf("\nWant: %s\nHave: %s\n", p.query, string(body))
			}
		})
	}
	wg.Wait()
}
