package github

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
)

type API struct {
	httpClient   *http.Client
	githubAPIURL string
	organization string
	repository   string
}

func NewGitubClient(organization string, repo string) *API {
	api := API{
		githubAPIURL: "https://api.github.com/repos/",
		organization: organization,
		repository:   repo,
	}
	api.httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return &api
}

func (c *API) Do(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	return res, err
}

func (c *API) Get(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.githubAPIURL+c.organization+"/"+c.repository+"/releases/latest", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}
