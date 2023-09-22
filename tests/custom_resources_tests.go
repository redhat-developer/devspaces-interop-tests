package tests

import (
	"context"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/redhat-developer/devspaces-interop-tests/internal/hlog"
	testContext "github.com/redhat-developer/devspaces-interop-tests/pkg/deploy/context"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var _ = KubeDescribe("[Custom Resources]", func() {
	ginkgo.It("Check if CRD already exist in Cluster", func() {
		hlog.Log.Info("Checking if CRD for Dev Spaces exist in cluster")
		// Move this client
		cfg, err := config.GetConfig()
		apiextensions, err := clientset.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
		// Make sure the CRD exist in cluster
		_, err = apiextensions.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), CRDName, metav1.GetOptions{})
		if err != nil {
			testContext.Instance.FoundCRD = false
		} else {
			testContext.Instance.FoundCRD = true
		}

		Expect(err).NotTo(HaveOccurred())
	})
})
