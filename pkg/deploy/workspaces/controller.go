package workspaces

import (
	"net/http"
)

// WorkspacesController useful to add logger and http client.
type Controller struct {
	httpClient *http.Client
}

// NewWorkspaceController creates a new WorkspacesController from a given client.
func NewWorkspaceController(c *http.Client) *Controller {
	return &Controller{
		httpClient: c,
	}
}
