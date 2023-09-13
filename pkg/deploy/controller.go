package deploy

import (
	"errors"
	"sync"
	"time"

	orgv2 "github.com/eclipse-che/che-operator/api/v2"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/internal/hlog"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/client"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy/context"
	testContext "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy/context"
	"go.uber.org/zap"
)

// TestHarnessController useful to add all kubernetes objects to cluster.
type TestHarnessController struct {
	sync.Mutex
	kubeClient *client.K8sClient
}

// NewTestHarnessController creates a new TestHarnessController from a given client.
func NewTestHarnessController(c *client.K8sClient) *TestHarnessController {
	return &TestHarnessController{
		kubeClient: c,
	}
}

func (c *TestHarnessController) DeployDevSpaces() bool {
	hlog.Log.Infof("Apply '%s' custom resource in namespace '%s'", crName, testContext.Config.DevSpacesNamespace)

	//Create a new Dev Spaces Custom resources into a giving namespace.
	if err := c.CreateCustomResource(); err != nil {
		hlog.Log.Panic("Failed to create custom resources in cluster", zap.Error(err))
	}

	hlog.Log.Infof("Successfully created '%s' custom resource in namespace '%s'", crName, testContext.Config.DevSpacesNamespace)

	// Check If all kubernetes objects for dev spaces performance are created in cluster
	// !Timeout is 10 minutes
	hlog.Log.Info("Waiting for Dev Spaces to be deployed in cluster")
	deploy, _ := c.WaitDevSpacesToBeUp(orgv2.ClusterPhaseActive)

	return deploy
}

// WatchCustomResource wait to deploy all performance/crw pods
func (c *TestHarnessController) WaitDevSpacesToBeUp(status orgv2.CheClusterPhase) (deployed bool, err error) {
	timeout := time.After(10 * time.Minute)
	tick := time.Tick(1 * time.Second)
	var clusterStarted = time.Now()

	stopCh := make(chan struct{})
	defer close(stopCh)

	for {
		select {
		case <-timeout:
			return false, errors.New("Error. Dev Spaces didn't deploy in 10 mins")
		case <-tick:
			customResource, _ := c.GetCustomResource()
			if customResource.Status.ChePhase == status {
				context.Instance.ClusterTimeUp = time.Since(clusterStarted).Seconds()
				hlog.Log.Info("Successfully deployed Dev Spaces in ", context.Instance.ClusterTimeUp)

				return true, nil
			}
		}
	}
}
