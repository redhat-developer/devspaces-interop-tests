package client

import (
	orgv2 "github.com/eclipse-che/che-operator/api/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/pkg/scheduler/api"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	Scheme             = runtime.NewScheme()
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemeGroupVersion = schema.GroupVersion{Group: orgv2.GroupVersion.Group, Version: orgv2.GroupVersion.Version}
)

type K8sClient struct {
	kubeClient *kubernetes.Clientset
}

// NewK8sClient creates kubernetes client wrapper with helper functions and direct access to k8s go client
func NewK8sClient() (*K8sClient, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	h := &K8sClient{kubeClient: client}
	return h, nil
}

// KubeRest Add che schema to kubernetes client runtime to perform api rest actions agains k8s clusters
func (c *K8sClient) KubeRest() crclient.Client {
	if err := AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
	if err := api.AddToScheme(Scheme); err != nil {
		panic(err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	client, err := crclient.New(cfg, crclient.Options{})

	if err != nil {
		panic("Failed to create client")
	}

	return client
}

// Kube returns the clientset for Kubernetes upstream.
func (c *K8sClient) Kube() kubernetes.Interface {
	return c.kubeClient
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&orgv2.CheCluster{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
