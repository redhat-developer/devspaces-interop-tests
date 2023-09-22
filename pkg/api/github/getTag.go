package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/redhat-developer/devspaces-interop-tests/internal/hlog"
)

type GitHubTagResponse struct {
	TagName string `json:"tag_name,omitempty"`
}

func (c *API) GetLatestCodeReadyWorkspacesTag() (tag string, err error) {
	gh := GitHubTagResponse{}

	response, err := c.Get(context.Background(), "aplication/json", nil)
	if err != nil {
		return "", err
	}
	err = json.NewDecoder(response.Body).Decode(&gh)
	if err != nil {
		hlog.Log.Fatal(err)
	}
	fmt.Println("TagName:", gh.TagName)
	version := strings.Split(gh.TagName, "-GA")[0]
	fmt.Println("Version Dev Spaces:", version)

	if _, err = c.Get(context.Background(), "aplication/json", nil); err != nil {
		return "", err
	}
	return version, nil
}
