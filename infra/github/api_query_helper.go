package github

import (
	"fmt"

	"github.com/rerost/vgithub-api/infra/github/api"
)

// ReviewRequestQuery is create query for PullRequests reviewer
func reviewRequestQuery(org string, repository string, pullRequestsSize int64, pullRequestsCursor string) api.Query {
	if pullRequestsCursor == "" {
		return api.Query(fmt.Sprintf(`
			{
				organization(login: "%s") {
					repository(name: "%s") {
						pullRequests(first: %d) {
							edges {
								cursor
							},
							nodes {
								author {
									login,
									avatarUrl
								}
								reviewRequests(first: 5) {
									nodes {
										requestedReviewer {
											__typename
											... on User {
												login,
												avatarUrl
											}
										}
									}
								}
							}
						}
					}
				}
			}
		`, org, repository, pullRequestsSize))
	}
	return api.Query(fmt.Sprintf(`
		{
			organization(login: "%s") {
				repository(name: "%s") {
					pullRequests(first: %d, after: "%s") {
						edges {
							cursor
						},
						nodes {
							author {
								login,
								avatarUrl
							}
							reviewRequests(first: 5) {
								nodes {
									requestedReviewer {
										__typename
										... on User {
											login,
											avatarUrl
										}
									}
								}
							}
						}
					}
				}
			}
		}
	`, org, repository, pullRequestsSize, pullRequestsCursor))
}

type reviewRequestsResponse struct {
	Data struct {
		Organization struct {
			Repository struct {
				PullRequests struct {
					Edges []struct {
						Cursor string `json:"cursor"`
					} `json:"edges"`
					Nodes []struct {
						Author         User `json:"author"`
						ReviewRequests struct {
							Nodes []struct {
								RequestedReviewer User `json:"requestedReviewer"`
							}
						}
					}
				} `json:"pullRequests"`
			} `json:"repository"`
		} `json:"organization"`
	} `json:"data"`
}

func repositoriesCountQuery(org string) api.Query {
	return api.Query(fmt.Sprintf(`
		{
			organization(login: "%s") {
				repositories(first: 1) {
					totalCount
				}
			}
		}
	`, org))
}

type repositoriesCountResponse struct {
	Data struct {
		Organization struct {
			Repositories struct {
				TotalCount int64 `json:"totalCount"`
			} `json:"repositories"`
		} `json:"organization"`
	} `json:"data"`
}

func pullRequestsCountQuery(org string, repositoriesSize int64, cursor string) api.Query {
	if cursor == "" {
		return api.Query(fmt.Sprintf(`
{
	organization(login: "%s") {
		repositories(first: %d) {
			edges {
				cursor
			},
			nodes {
				name,
				pullRequests(last: 1) {
					totalCount
				}
			}
		}
	}
}
`, org, repositoriesSize))

	}
	return api.Query(fmt.Sprintf(`
		  {
		  	organization(login: "%s") {
		  		repositories(first: %d, after: "%s") {
						edges {
							cursor
						},
		  			nodes {
		  				name,
		  				pullRequests(last: 1) {
		  					totalCount
		  				}
		  			}
		  		}
		  	}
		  }
		`, org, repositoriesSize, cursor))
}

type pullRequestsCountResponse struct {
	Data struct {
		Organization struct {
			Repositories struct {
				Edges []struct {
					Cursor string `json:"cursor"`
				} `json:"edges"`
				Nodes []struct {
					Name         string `json:"name"`
					PullRequests struct {
						TotalCount int64 `json:"totalCount"`
					}
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"organization"`
	} `json:"data"`
}
