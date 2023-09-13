package workspaces

import (
	"context"

	v1alpha2 "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	v2 "github.com/eclipse-che/che-operator/api/v2"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/internal/hlog"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/client"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy"
	"go.uber.org/zap"
)

const (
	WorkspaceRunningStatus = "RUNNING"
)

// RunWorkspace create a new performance from a given devFile and call an method to get measure time for performance after a performance pod is up and ready
func (w *Controller) TestWorkspaceStartAndDelete(devWorkspaceDefenition *v1alpha2.DevWorkspace) (workspaceID *Workspace, err error) {
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create kubernetes client.", zap.Error(err))
	}

	ctrl := deploy.NewTestHarnessController(k8sClient)
	resource, err := ctrl.GetCustomResource()
	if err != nil {
		hlog.Log.Panic("Failed to get Custom Resource.", zap.Error(err))
	}

	workspace, err := w.CreateAndRunWorkspace(devWorkspaceDefenition)

	if err != nil {
		hlog.Log.Panic("Error on create performance.", zap.Error(err))
	}

	hlog.Log.Info("Successfully started workspace", zap.String("workspaceID", workspace.ID))

	w.DeleteWorkspace(getCheUrl(resource), workspace)

	return workspace, err
}

// CreateAndRunWorkspace creates and runs an workspace using DevWorkspace yaml
func (w *Controller) CreateAndRunWorkspace(devWorkspaceDefenition *v1alpha2.DevWorkspace) (workspace *Workspace, err error) {

	/*request, _ := http.NewRequest("POST", cheURL+"/api/workspace/devfile", bytes.NewBuffer(workspaceDefinition))
	request.Header.Add("Content-Type", "text/yaml")

	response, err := w.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.Status == "409" {
		return nil, errors.New("workspace already exist")
	}

	err = json.NewDecoder(response.Body).Decode(&workspace)
	startWorkspace, err := w.startWorkspace(cheURL, workspace.ID)
	if err != nil && !startWorkspace {
		return nil, errors.New("error sending request to start workspace")
	}*/

	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create kubernetes client.", zap.Error(err))
	}

	if err := k8sClient.KubeRest().Create(context.TODO(), devWorkspaceDefenition); err != nil {
		hlog.Log.Panic("Failed to create devworkspace template and devworkspace components.", zap.Error(err))
	}

	statusWorkspace, err := w.statusWorkspace(workspace, WorkspaceRunningStatus)
	if !statusWorkspace {
		return nil, err
	}

	return workspace, err
}

/*func (w *Controller) startWorkspace(cheUrl string, workspaceID string) (boolean bool, err error) {
	request, err := http.NewRequest("POST", cheUrl+"/api/workspace/"+workspaceID+"/runtime", nil)
	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")

	_, err = w.httpClient.Do(request)

	if err != nil {
		return false, err
	}

	return true, err
}*/

func getCheUrl(che *v2.CheCluster) string {
	return che.Status.CheURL
}
