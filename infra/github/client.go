package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/rerost/vgithub-api/infra/github/api"
)

// Client is interface for github
type Client interface {
	GetReviewRequests(org string) ([]*ReviewRequest, error)
}

const (
	maxGetRepository   = 50
	maxGetPullRequests = 50
	graphqlURL         = "https://api.github.com/graphql"
)

type clientImp struct {
	apiClient api.Client
}

// NewClient is create github api client
func NewClient(httpClient *http.Client, token string) Client {
	url, err := url.Parse(graphqlURL)
	if err != nil {
		// Beacause graphqlURL is constant, but parse error occured
		panic(err)
	}

	apiClient, err := api.NewClient(httpClient, url, token)
	if err != nil {
		panic(err)
	}

	return &clientImp{apiClient: apiClient}
}

func (c *clientImp) GetReviewRequests(org string) ([]*ReviewRequest, error) {
	if c == nil {
		return nil, nil
	}

	// Get repositories count, pull_requests count
	repositoriesCountResponse, err := c.getRepositoriesCount(org)
	if err != nil {
		log.Printf("Failed to get repositories count: %s", err)
		return nil, err
	}

	repositoriesTotalCount := repositoriesCountResponse.Data.Organization.Repositories.TotalCount

	var repositoryNames []string
	pullRequestsCountMap := map[string]int64{}

	var i int64
	var cursor string
	for i = 0; i < repositoriesTotalCount/maxGetPullRequests+1; i++ {
		response, err := c.getPullRequestsCount(org, maxGetPullRequests, cursor)
		if err != nil {
			log.Printf("getPullRequestsCount error: %s", err)
			return nil, err
		}

		edges := response.Data.Organization.Repositories.Edges

		// edges == [] => nodes == [] が成立するので
		if len(edges) == 0 {
			break
		}

		// これは、edgesの一番最後の`cursor`を指定することで、今手に入れたレスポンスと重ならないことが前提
		cursor = edges[len(edges)-1].Cursor

		for _, repository := range response.Data.Organization.Repositories.Nodes {
			repositoryNames = append(repositoryNames, repository.Name)
			pullRequestsCountMap[repository.Name] = repository.PullRequests.TotalCount
		}
	}

	// Request query with paging
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	var ReviewRequests []*ReviewRequest
	for _, repositoryName := range repositoryNames {
		wg.Add(1)
		go func(repositoryName string) {
			defer wg.Done()

			var reviewRequests []*ReviewRequest
			var prCursor string
			for i = 0; i < pullRequestsCountMap[repositoryName]/maxGetPullRequests+1; i++ {
				response, err := c.getReviewRequests(org, repositoryName, int64(maxGetPullRequests), prCursor)
				if err != nil {
					log.Printf("Occured error in getReviewRequests: %s", err)
					return
				}

				pullReqeusts := response.Data.Organization.Repository.PullRequests
				if len(pullReqeusts.Edges) == 0 {
					continue
				}
				prCursor = pullReqeusts.Edges[len(pullReqeusts.Edges)-1].Cursor
				for _, pullRequest := range pullReqeusts.Nodes {
					for _, reviewRequest := range pullRequest.ReviewRequests.Nodes {
						reviewRequest := &ReviewRequest{
							Repository: repositoryName,
							Reviewer:   reviewRequest.RequestedReviewer,
							Reviewee:   pullRequest.Author,
						}

						reviewRequests = append(reviewRequests, reviewRequest)
					}
				}
			}

			mutex.Lock()
			ReviewRequests = append(ReviewRequests, reviewRequests...)
			mutex.Unlock()

		}(repositoryName)
	}
	wg.Wait()

	// TODO Convert to ReviewRequests

	return ReviewRequests, nil
}

func (c *clientImp) getReviewRequests(org string, repositoryName string, pullRequestsSize int64, pullRequestsCursor string) (*reviewRequestsResponse, error) {
	query := reviewRequestQuery(org, repositoryName, pullRequestsSize, pullRequestsCursor)
	res, err := c.apiClient.Request(query)
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		log.Printf("Failed to github request err: %v", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		log.Printf("Failed to github request. status code %v", res.StatusCode)
		body, _ := ioutil.ReadAll(res.Body)

		return nil, fmt.Errorf("StatusCode is %v, Response: %v", res.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to iouti.ReadAll err: %v", err)
		return nil, err
	}

	queryResponse := reviewRequestsResponse{}
	err = json.Unmarshal(body, &queryResponse)

	if err != nil {
		log.Printf("Failed to json.Unmarshal err: %v", err)
		return nil, err
	}

	return &queryResponse, nil
}

func (c *clientImp) getRepositoriesCount(org string) (*repositoriesCountResponse, error) {
	query := repositoriesCountQuery(org)
	res, err := c.apiClient.Request(query)

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		log.Printf("Failed to github request err: %v", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		log.Printf("Failed to github request. status code %v", res.StatusCode)
		body, _ := ioutil.ReadAll(res.Body)

		return nil, fmt.Errorf("StatusCode is %v, Response: %v", res.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to iouti.ReadAll err: %v", err)
		return nil, err
	}

	queryResponse := repositoriesCountResponse{}
	err = json.Unmarshal(body, &queryResponse)

	if err != nil {
		log.Printf("Failed to json.Unmarshal err: %v", err)
		return nil, err
	}

	return &queryResponse, nil
}

func (c *clientImp) getPullRequestsCount(org string, repositroeiesSize int64, cursor string) (*pullRequestsCountResponse, error) {
	query := pullRequestsCountQuery(org, repositroeiesSize, cursor)
	res, err := c.apiClient.Request(query)

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		log.Printf("Failed to github request err: %v", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		log.Printf("Failed to github request. status code %v", res.StatusCode)
		return nil, fmt.Errorf("StatusCode is %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to iouti.ReadAll err: %v", err)
		return nil, err
	}

	queryResponse := pullRequestsCountResponse{}
	err = json.Unmarshal(body, &queryResponse)

	if err != nil {
		log.Printf("Failed to json.Unmarshal err: %v", err)
		return nil, err
	}

	return &queryResponse, nil
}
