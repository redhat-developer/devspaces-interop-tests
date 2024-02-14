package tests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"sigs.k8s.io/yaml"

	v1alpha2 "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	orgv2 "github.com/eclipse-che/che-operator/api/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/redhat-developer/devspaces-interop-tests/internal/hlog"
	"github.com/redhat-developer/devspaces-interop-tests/pkg/client"
	"github.com/redhat-developer/devspaces-interop-tests/pkg/deploy"
	"github.com/redhat-developer/devspaces-interop-tests/pkg/deploy/workspaces"
	"github.com/sirupsen/logrus"
)

var workspaceDefinition []byte
var devWorkspacePatchedYaml []byte

var _ = Describe("[WORKSPACES]", func() {
	Context("Create workspace from devfile registry", func() {
		k8sClient, err := client.NewK8sClient()
		httpClient, err := client.NewHttpClient()
		cheCluster := &orgv2.CheCluster{}
		testHarnessController := deploy.NewTestHarnessController(k8sClient)

		It("Obtain CheCluster object", func() {
			cheCluster, err = testHarnessController.GetCustomResource()
			Expect(err).NotTo(HaveOccurred())
			Expect(cheCluster.Status.CheURL).ToNot(BeNil())
			Expect(cheCluster.Status.DevfileRegistryURL).ToNot(BeNil())
		})

		It("Obtain and patch Devfile from DevFile Registry", func() {
			request, err := http.NewRequest("GET", ObtainJavaDevFileUrl(cheCluster), nil)
			Expect(err).NotTo(HaveOccurred())

			response, err := httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).Should(Equal(200))
			if response.StatusCode == http.StatusOK {
				workspaceDefinition, err = io.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				hlog.Log.Infof("Java Workspace Definition")
				patchWorkspaceDefenition := [][]byte{workspaceDefinition, []byte("routingClass: che")}
				devWorkspacePatchedYaml = bytes.Join(patchWorkspaceDefenition, []byte("  "))
				fmt.Println(string(devWorkspacePatchedYaml))
			}
		})

		It("Create and start Workspace", func() {
			Skip("Run a workspace start in Dev Spaces requires a new implementation.The test is temporary skipped, while investigating.")
			hlog.Log.Info("Starting a new workspace")
			ctrl := workspaces.NewWorkspaceController(httpClient)

			_, err = ctrl.TestWorkspaceStartAndDelete(GetDevWorkspaceYaml(devWorkspacePatchedYaml))

			Expect(err).NotTo(HaveOccurred())
		})
	})
})

func ObtainJavaDevFileUrl(cheCluster *orgv2.CheCluster) string {
	return cheCluster.Status.DevfileRegistryURL + "/devfiles/java-maven-lombok__lombok-project-sample/devworkspace-che-code-latest.yaml"
}

func GetDevWorkspaceYaml(dwYamlFile []byte) *v1alpha2.DevWorkspace {
	devWorkspace := &v1alpha2.DevWorkspace{}
	if err := ReadObjectInto(dwYamlFile, devWorkspace); err != nil {
		logrus.Fatalf("Failed to read devworkspace yaml from '%s', cause: %v", dwYamlFile, err)
	}

	fmt.Println(devWorkspace)
	return devWorkspace
}

func ReadObjectInto(dwYamlFile []byte, obj interface{}) error {
	if err := yaml.Unmarshal(dwYamlFile, obj); err != nil {
		return err
	}

	return nil
}
