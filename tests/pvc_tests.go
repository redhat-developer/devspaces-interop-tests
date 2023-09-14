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

var _ = KubeDescribe("[PersistentVolumeClaims]", func() {
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create k8s client go", zap.Error(err))
	}

	ginkgo.It("PVC `postgres-data` should be created", func() {
		hlog.Log.Info("Check if PVC for postgres was created")
		secret, err := k8sClient.Kube().CoreV1().PersistentVolumeClaims(testContext.Config.DevSpacesNamespace).Get(context.TODO(), PostgresPVCName, metav1.GetOptions{})

		if err != nil {
			hlog.Log.Error("Error on getting info about pvc status")
		}

		Expect(secret).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred(), "failed to get pvc %v\n", PostgresPVCName)
	})
})
