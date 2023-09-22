package olm

import (
	"os"

	operatorv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	olmversioned "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"
	"github.com/redhat-developer/devspaces-interop-tests/internal/hlog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	Scheme               = runtime.NewScheme()
	v1AlphaSchemeBuilder = runtime.NewSchemeBuilder(addV1AlphaKnownTypes)
	AddV1AlphaToScheme   = v1AlphaSchemeBuilder.AddToScheme
	SchemeV1AlphaVersion = schema.GroupVersion{Group: operatorsv1alpha1.SchemeGroupVersion.Group, Version: operatorsv1alpha1.SchemeGroupVersion.Version}

	v1SchemeBuilder = runtime.NewSchemeBuilder(addV1KnownTypes)
	AddV1ToScheme   = v1SchemeBuilder.AddToScheme
	SchemeV1Version = schema.GroupVersion{Group: operatorv1.SchemeGroupVersion.Group, Version: operatorv1.SchemeGroupVersion.Version}
)

// NewK8sClient creates kubernetes client wrapper with helper functions and direct access to k8s go client
func NewOLMK8sClient() (*Client, error) {
	if err := AddV1ToScheme(scheme.Scheme); err != nil {
		hlog.Log.Fatalf("Failed to add v1 api operator scheme")
	}

	if err := AddV1AlphaToScheme(scheme.Scheme); err != nil {
		hlog.Log.Fatalf("Failed to add v1alpha api operator scheme")
	}

	if err := api.AddToScheme(Scheme); err != nil {
		hlog.Log.Fatalf("Failed to add CRD to scheme")
	}

	cfg, err := config.GetConfig()
	if err != nil {
		hlog.Log.Error(err, "Failed to create client config")
		os.Exit(1)
	}

	client, err := client.New(cfg, client.Options{})

	if err != nil {
		hlog.Log.Error(err, "Failed to create client")
		os.Exit(1)
	}

	clientsOLM, err := olmversioned.NewForConfig(cfg)

	if err != nil {
		hlog.Log.Error(err, "Failed to create olm client")
		os.Exit(1)
	}

	h := &Client{client, clientsOLM}

	return h, nil
}

func addV1AlphaKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeV1AlphaVersion,
		&operatorsv1alpha1.CatalogSource{},
		&operatorsv1alpha1.Subscription{},
	)
	metav1.AddToGroupVersion(scheme, SchemeV1AlphaVersion)
	return nil
}

func addV1KnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeV1Version,
		&operatorv1.OperatorGroup{},
	)
	metav1.AddToGroupVersion(scheme, SchemeV1Version)
	return nil
}
