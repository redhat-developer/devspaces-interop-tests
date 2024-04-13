package tests

import (
	"context"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/redhat-developer/devspaces-interop-tests/internal/hlog"
	"github.com/redhat-developer/devspaces-interop-tests/pkg/client"
	testContext "github.com/redhat-developer/devspaces-interop-tests/pkg/deploy/context"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = KubeDescribe("[Secrets]", func() {
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create k8s client go", zap.Error(err))
	}

	ginkgo.It("Secret `self-signed-certificate` should exist", func() {
		ginkgo.Skip("This test is skipped.")
		hlog.Log.Info("Checking secrets created for Dev Spaces")
		secret, err := k8sClient.Kube().CoreV1().Secrets(testContext.Config.DevSpacesNamespace).Get(context.TODO(), secretSelfSignedCrt, metav1.GetOptions{})

		if err != nil {
			hlog.Log.Error("Error on get info about secrets")
		}

		Expect(secret).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred(), "failed to get secretName %v\n", secretSelfSignedCrt)
	})
})
