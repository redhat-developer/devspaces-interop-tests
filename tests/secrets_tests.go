package tests

import (
	 "context"

	 "github.com/onsi/ginkgo"
	 . "github.com/onsi/gomega"
	 "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/internal/hlog"
	 "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/client"
	 testContext "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy/context"
	 "go.uber.org/zap"
	 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = KubeDescribe("[Secrets]", func() {
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create k8s client go", zap.Error(err))
	}

	ginkgo.It("Secret `self-signed-certificate` should exist", func() {
		hlog.Log.Info("Checking secrets created for code ready workspaces")
		secret, err := k8sClient.Kube().CoreV1().Secrets(testContext.Config.DevSpacesNamespace).Get(context.TODO(), secretSelfSignedCrt, metav1.GetOptions{})

		if err != nil {
			hlog.Log.Error("Error on get info about secrets")
		}

		Expect(secret).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred(), "failed to get secretName %v\n", secretSelfSignedCrt)
	})
})
