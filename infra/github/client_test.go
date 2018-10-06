package github_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/rerost/vgithub-api/infra/github"
)

func TestGetReviewRequests(t *testing.T) {
	// TOOD(@rerost) Use dummy server
	client := github.NewClient(new(http.Client), "dummytoken")
	_, err := client.GetReviewRequests("github")
	if err != nil {
		t.Error(err)
	}
}

func TestGetReviewRequestsWithRealServer(t *testing.T) {
	client := github.NewClient(new(http.Client), os.Getenv("TEST_GITHUB_TOKEN"))
	res, err := client.GetReviewRequests("rerost-test")
	if err != nil {
		t.Errorf("Failed to GetReviewRequests: %v", err)
	}

	if len(res) == 0 {
		t.Error("Review is not 0")
	}
}
