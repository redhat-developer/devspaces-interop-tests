package deploy

import (
	"context"

	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/internal/hlog"
	testContext "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy/context"

	orgv2 "github.com/eclipse-che/che-operator/api/v2"
	gherr "github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// GetCustomResource makes an models request to K8s API to get Che Cluster
func (c *TestHarnessController) GetCustomResource() (*orgv2.CheCluster, error) {
	CheCluster := orgv2.CheCluster{}

	err := c.kubeClient.KubeRest().Get(context.TODO(), types.NamespacedName{Namespace: testContext.Config.DevSpacesNamespace, Name: crName}, &CheCluster)
	if err != nil {
		hlog.Log.Error("Error to get custom resource", zap.Error(err))

		return nil, err
	}

	return &CheCluster, nil
}

// CreateCustomResource makes an models request to K8s API to create Che Cluster
func (c *TestHarnessController) CreateCustomResource() (err error) {
	if err := c.kubeClient.KubeRest().Create(context.TODO(), GetCustomResourceSpec()); err != nil {
		if errors.IsAlreadyExists(err) {
			return gherr.Wrapf(err, "Failed to create devspaces custom resource in cluster. '%s' Already exists ", crName)
		} else {
			return gherr.Wrapf(err, "Failed to create devspaces custom resource in cluster %v", err)
		}
	}

	return err
}

// DeleteCustomResource makes an models request to K8s API to delete Che Cluster
func (c *TestHarnessController) DeleteCustomResource() (err error) {
	if err := c.kubeClient.KubeRest().Delete(context.TODO(), GetCustomResourceSpec()); err != nil {
		return gherr.Wrapf(err, "Failed to delete devspaces custom resource in cluster %v", err)
	}

	return err
}

// GetCustomResourceSpec returns CR che cluster k8s object
func GetCustomResourceSpec() *orgv2.CheCluster {
	return &orgv2.CheCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: testContext.Config.DevSpacesNamespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       DevSpacesKind,
			APIVersion: DevSpacesAPIVersion,
		},
		Spec: orgv2.CheClusterSpec{
			DevEnvironments: orgv2.CheClusterDevEnvironments{
				DefaultNamespace: orgv2.DefaultNamespace{
					Template: "<username>-devspaces",
				},
			},
			Components: orgv2.CheClusterComponents{
				CheServer: orgv2.CheServer{
					Debug:    NewBoolPointer(false),
					LogLevel: "INFO",
				},
			},
		},
	}
}

// NewBoolPointer returns `bool` pointer to value in the memory.
// Unfortunately golang hasn't got syntax to create `bool` pointer.
func NewBoolPointer(value bool) *bool {
	variable := value
	return &variable
}
