package workspaces

import (
	"net/http"

	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/internal/hlog"
	"go.uber.org/zap"
)

const (
	WorkspaceStoppedStatus = "STOPPED"
)

// DeleteWorkspace delete a workspace from a given workspace_id
func (w *Controller) DeleteWorkspace(cheURL string, workspace *Workspace) (err error) {
	// var keycloakAuth *KeycloakAuth

	hlog.Log.Info("Cleaning workspace from cluster...", zap.String("workspaceID", workspace.ID))

	/*if keycloakAuth, err = w.KeycloakToken(keycloakUrl); err != nil {
		hlog.Log.Error("Failed to get user token ", zap.Error(err))
	}*/

	if err = w.stopWorkspace(cheURL, workspace); err != nil {
		hlog.Log.Error("Failed to get user token ", zap.Error(err))
	}
	statusWorkspace, err := w.statusWorkspace(workspace, WorkspaceStoppedStatus)

	if !statusWorkspace {
		return err
	}

	request, err := http.NewRequest("DELETE", cheURL+"/api/workspace/"+workspace.ID, nil)

	if err != nil {
		hlog.Log.Error("Failed to delete workspace", zap.Error(err))
	}

	/*if keycloakAuth, err = w.KeycloakToken(keycloakUrl); err != nil {
		hlog.Log.Error("Failed to get user token ", zap.Error(err))
	}*/

	// request.Header.Add("Authorization", "Bearer "+keycloakAuth.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	_, err = w.httpClient.Do(request)

	return err
}

// StopWorkspace stop a workspace from a given workspace_id
func (w *Controller) stopWorkspace(cheURL string, workspace *Workspace) (err error) {
	// var keycloakAuth *KeycloakAuth

	request, err := http.NewRequest("DELETE", cheURL+"/api/workspace/"+workspace.ID+"/runtime", nil)

	if err != nil {
		hlog.Log.Error("Failed to stop workspace", zap.Error(err))
	}

	/*if keycloakAuth, err = w.KeycloakToken(keycloakUrl); err != nil {
		hlog.Log.Error("Failed to get user token ", zap.Error(err))
	}*/

	// request.Header.Add("Authorization", "Bearer "+keycloakAuth.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	_, err = w.httpClient.Do(request)

	return err
}
