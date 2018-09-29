package api_test

import (
	"net/http"
	"net/url"
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
	dummyURL, _ := url.Parse("http://momokuri-rock.com")

	inOutPairs := []struct {
		test string
		in   string // Token
		out  string // Error message
	}{
		{
			test: "When passed token",
			in:   "asd",
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
