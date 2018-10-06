package github

import (
	"net/url"
	"strings"
)

// ReviewRequest is type of Github repsonse for reviews
type ReviewRequest struct {
	Repository string `json:"repository"`
	Reviewee   User   `json:"reviewer"`
	Reviewer   User   `json:"reviewee"`
}

// User is github user info
type User struct {
	Login     string `json:"login"`
	AvatarURL Link   `json:"avatarUrl"`
}

type Link struct {
	URL *url.URL
}

func (l *Link) MarshalJSON() ([]byte, error) {
	return []byte(l.URL.String()), nil
}
func (l *Link) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	l.URL = u
	return nil
}
