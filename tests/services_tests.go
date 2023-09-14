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

var _ = ginkgo.Describe("[Services]", func() {
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create k8s client go", zap.Error(err))
	}

	ginkgo.It("Check if devspaces services exist in cluster", func() {
		hlog.Log.Info("Checking if all services for Dev Spaces")
		services, err := k8sClient.Kube().CoreV1().Services(testContext.Config.DevSpacesNamespace).List(context.TODO(), metav1.ListOptions{})

		Expect(services).NotTo(BeNil())

		confmap := map[string]string{}
		for _, v := range services.Items {
			confmap[v.Name] = v.Name
		}

		Expect(confmap["che-gateway"]).NotTo(BeEmpty())
		Expect(confmap["che-host"]).NotTo(BeEmpty())
		Expect(confmap["devspaces-dashboard"]).NotTo(BeEmpty())
		Expect(confmap["plugin-registry"]).NotTo(BeEmpty())
		Expect(confmap["postgres"]).NotTo(BeEmpty())
		Expect(confmap["devfile-registry"]).NotTo(BeEmpty())

		Expect(err).NotTo(HaveOccurred())
	})
})
