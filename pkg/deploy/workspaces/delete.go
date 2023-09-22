package workspaces

import (
	"github.com/redhat-developer/devspaces-interop-tests/internal/hlog"
	"go.uber.org/zap"
)

// DeleteWorkspace deletes a workspace from Dev Spaces
func (w *Controller) DeleteWorkspace(workspaceName string, userNamespace string) (err error) {

	hlog.Log.Info("Cleaning workspace from cluster...", zap.String("workspaceName", workspaceName))

	// Get workspace status
	return nil
}
