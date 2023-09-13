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

var _ = KubeDescribe("[ConfigMaps]", func() {
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create k8s client go", zap.Error(err))
	}

	ginkgo.It("Config map `che` should exist", func() {
		hlog.Log.Info("Checking `che` config map integrity")

		che, err := k8sClient.Kube().CoreV1().ConfigMaps(testContext.Config.DevSpacesNamespace).Get(context.TODO(), DevSpacesConfigMap, metav1.GetOptions{})

		Expect(che).NotTo(BeNil())

		if err != nil {
			hlog.Log.Error("Error to verify `che` config map")
		}

		Expect(err).NotTo(HaveOccurred())
	})
})
