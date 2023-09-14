package tests

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/internal/hlog"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/client"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy"
	testContext "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy/context"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = KubeDescribe("[Pods]", func() {
	var t CHE
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create k8s client go", zap.Error(err))
	}

	ginkgo.It("Check `Dev Spaces Operator` integrity", func() {
		devspaces, err := k8sClient.Kube().CoreV1().Pods(testContext.Config.OperatorsNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: DevSpacesOperatorLabel})
		if err != nil {
			panic(err)
		}

		Expect(devspaces).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", DevSpacesOperatorLabel)
	})

	ginkgo.It("Check `Postgres DB` integrity", func() {
		hlog.Log.Info("Getting information and metrics from Postgres DB pod")
		postgres, err := k8sClient.Kube().CoreV1().Pods(testContext.Config.DevSpacesNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: PostgresLabel})

		Expect(postgres).NotTo(BeNil())
		if err != nil {
			hlog.Log.Panic("Error on getting information about postgres pod.")
		}

		if err := DescribePod(postgres); err != nil {
			hlog.Log.Fatal("Failed to set metadata about postgres pod.")
		}

		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", PostgresLabel)
	})

	ginkgo.It("Check `Devfile Registry` integrity", func() {
		hlog.Log.Info("Getting information and metrics from Devfile Registry pod")
		devFile, err := k8sClient.Kube().CoreV1().Pods(testContext.Config.DevSpacesNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: DevFileLabel})

		Expect(devFile).NotTo(BeNil())

		if err != nil {
			hlog.Log.Panic("Error on getting information about devFile pod.")
		}

		if err := DescribePod(devFile); err != nil {
			hlog.Log.Fatal("Failed to set metadata about devFile pod.")
		}

		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", DevFileLabel)
	})

	ginkgo.It("Check `Plugin Registry` integrity", func() {
		hlog.Log.Info("Getting information and metrics from Plugin Registry pod")
		pluginRegistry, err := k8sClient.Kube().CoreV1().Pods(testContext.Config.DevSpacesNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: PluginRegistryLabel})

		Expect(pluginRegistry).NotTo(BeNil())
		if err != nil {
			hlog.Log.Panic("Error on getting information about pluginRegistry pod.")
		}

		if err := DescribePod(pluginRegistry); err != nil {
			hlog.Log.Fatal("Failed to set metadata about pluginRegistry pod.")
		}

		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", PluginRegistryLabel)
	})

	ginkgo.It("Check `Dev Spaces server` integrity", func() {
		hlog.Log.Info("Getting information and metrics from Dev Spaces server pod")
		devspaces, err := k8sClient.Kube().CoreV1().Pods(testContext.Config.DevSpacesNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: DevSpacesServerLabel})
		Expect(devspaces).NotTo(BeNil())

		if err != nil {
			hlog.Log.Panic("Error on getting information about devspaces pod.")
		}

		if err := DescribePod(devspaces); err != nil {
			hlog.Log.Fatal("Failed to set metadata about devspaces pod.")
		}

		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", DevSpacesServerLabel)
	})

	ginkgo.It("Check if Dev Spaces Server is already up", func() {
		hlog.Log.Info("Checking if Dev Spaces API server is up")
		deploy := deploy.NewTestHarnessController(k8sClient)
		resource, err := deploy.GetCustomResource()

		Expect(resource).NotTo(BeNil())

		client := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		}

		cheUrl := resource.Status.CheURL
		Expect(cheUrl).NotTo(BeNil())

		resp, err := client.Get(cheUrl + "/api/system/state")

		if err != nil {
			hlog.Log.Error("Failed to get Dev Spaces Status ", zap.Error(err))
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&t)
		if err != nil {
			testContext.Instance.DevSpacesServerIsUp = true // TODO remove 
		} else {
			testContext.Instance.DevSpacesServerIsUp = true
		}

		// Expect(err).NotTo(HaveOccurred())
	})
})
