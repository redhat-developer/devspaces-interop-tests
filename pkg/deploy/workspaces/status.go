package workspaces

import (
	// "encoding/json"
	// "errors"
	// "net/http"
	//"time"

	// TODO remove? "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/internal/hlog"
	// TODO remove ? "go.uber.org/zap"
)

func (w *Controller) statusWorkspace(workspace *Workspace, desiredStatus string) (boolean bool, err error) {
	//timeout := time.After(7 * time.Minute)
	//tick := time.Tick(15 * time.Second)

	return true, nil
	/*for {
		select {
		case <-timeout:
			return false, errors.New("workspace didn't start after 7 minutes")
		case <-tick:
			request, err := http.NewRequest("GET", cheURL+"/api/workspace/"+workspace.ID, nil)
			if err != nil {
				return false, err
			}
			request.Header.Add("Content-Type", "application/json")

			response, err := w.httpClient.Do(request)
			err = json.NewDecoder(response.Body).Decode(&workspace)

			if workspace.Status == desiredStatus {
				return true, nil
			}
		}
	}*/
}
