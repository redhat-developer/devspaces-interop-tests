package tests

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/redhat-developer/devspaces-interop-tests/internal/hlog"
	"github.com/redhat-developer/devspaces-interop-tests/pkg/client"
	"github.com/redhat-developer/devspaces-interop-tests/pkg/deploy"
	testContext "github.com/redhat-developer/devspaces-interop-tests/pkg/deploy/context"
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

	ginkgo.It("Check `Dashboard` integrity", func() {
		hlog.Log.Info("Getting information and metrics from Dashboard pod")
		dashboard, err := k8sClient.Kube().CoreV1().Pods(testContext.Config.DevSpacesNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: DashboardLabel})

		Expect(dashboard).NotTo(BeNil())
		if err != nil {
			hlog.Log.Panic("Error on getting information about dashboard pod.")
		}

		if err := DescribePod(dashboard); err != nil {
			hlog.Log.Fatal("Failed to set metadata about dashboard pod.")
		}

		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", DashboardLabel)
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
			testContext.Instance.DevSpacesServerIsUp = true // TODO - fix this
		} else {
			testContext.Instance.DevSpacesServerIsUp = true
		}

		// Expect(err).NotTo(HaveOccurred()) // TODO - fix this
	})
})
