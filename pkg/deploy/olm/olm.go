package olm

import (
	"context"
	"sync"

	olmversioned "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"
	v1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/internal/hlog"
	testContext "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy/context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	cl "sigs.k8s.io/controller-runtime/pkg/client"
	"go.uber.org/zap"
)

// OlmController useful to add all kubernetes objects to cluster.
type OlmController struct {
	sync.Mutex
	k8s *Client
}

type Client struct {
	cl.Client
	OLM olmversioned.Interface
}

// NewOLMController creates a new OlmController from a given client.
func NewOLMController(k8s *Client) *OlmController {
	return &OlmController{
		k8s: k8s,
	}
}

// InstallOLMOperator desc
func (o *OlmController) InstallOLMOperator() {
	hlog.Log.Infof("Dev Spaces Operator will be installed in namespace '%s'", testContext.Config.OperatorsNamespace)
	if ns := o.k8s.Get(context.TODO(), types.NamespacedName{Name: testContext.Config.OperatorsNamespace}, GetNamespaceSpec()); ns != nil {
		if errors.IsNotFound(ns) {
			hlog.Log.Infof("Namespace %s doesn't exist. Creating new one...", testContext.Config.OperatorsNamespace)
			if err := o.k8s.Create(context.TODO(), GetNamespaceSpec()); err != nil {
				hlog.Log.Fatalf("Failed to create namespace %s: %v", testContext.Config.OperatorsNamespace, err)
			}
		}
	}

	if err := o.InstallSubscription(); err != nil {
		hlog.Log.Fatal(err)
	}
}

// GetNamespaceSpec return namespace object
func GetNamespaceSpec() *v1.Namespace {
	return &v1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: testContext.Config.OperatorsNamespace,
		},
	}
}

// GetClusterServiceVersion makes an models request to K8s API to get CSV
func (o *OlmController) GetClusterServiceVersion() (*v1alpha1.ClusterServiceVersion, error) {
	CSV := &v1alpha1.ClusterServiceVersion{}

	CSV, err := o.k8s.OLM.OperatorsV1alpha1().
		ClusterServiceVersions(testContext.Config.OperatorsNamespace).
		Get(context.TODO(), testContext.Config.CSVName, metav1.GetOptions{})
	if err != nil {
		hlog.Log.Error("Error to get csv resource", zap.Error(err))

		return nil, err
	}

	return CSV, nil
}
