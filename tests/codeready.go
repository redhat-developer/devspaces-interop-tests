package tests

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"github.com/onsi/ginkgo"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/internal/hlog"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/client"
	testContext "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy/context"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
)

const (
	// Config Maps constants
	DevSpacesConfigMap = "che"

	// Pod Names used to get info
	DevSpacesOperatorLabel = "olm.owner.kind=ClusterServiceVersion"
	PostgresLabel          = "component=postgres"
	DevFileLabel           = "component=devfile-registry"
	PluginRegistryLabel    = "component=plugin-registry"
	DevSpacesServerLabel   = "component=devspaces"

	//Custom Resource name to get info
	CRDName = "checlusters.org.eclipse.che"

	// PVC name used for postgres-data
	PostgresPVCName = "postgres-data"

	// Secret name used for add ca.crt
	secretSelfSignedCrt = "postgres-credentials"
)

type PodInfo struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
}

type CHE struct {
	CheStatus string `json:"status"`
}

// KubeDescribe is wrapper function for ginkgo describe. .
func KubeDescribe(text string, body func()) bool {
	return ginkgo.Describe("[Dev Spaces Test Harness] "+text, body)
}

// DescribePod set metadata and metrics about a specific pod
func DescribePod(pod *v1.PodList) (err error) {
	var podInfo testContext.CodeReadyPods

	for _, v := range pod.Items {
		podInfo.Name = v.Name
		podInfo.Status = v.Status.Phase
		podInfo.Labels = v.Labels
		DescribePodLogs(v.Name)

		for _, val := range v.Spec.Containers {
			podInfo.DockerImage = val.Image
		}
		a := append(testContext.Instance.CodeReadyPodsInfo, podInfo)

		testContext.Instance.CodeReadyPodsInfo = a
	}
	return err
}

// DescribePodLogs get all logs from a specific pod and write to a file
func DescribePodLogs(podName string) {
	podLogOpts := v1.PodLogOptions{}
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		hlog.Log.Panic("Failed to create k8s client go", zap.Error(err))
	}

	req := k8sClient.Kube().CoreV1().Pods(testContext.Config.DevSpacesNamespace).GetLogs(podName, &podLogOpts)
	podLogs, _ := req.Stream(context.TODO())

	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, podLogs)

	str := buf.Bytes()

	err = ioutil.WriteFile("/test-run-results/devspaces_"+podName+".log", str, 0644)
	if err != nil {
		hlog.Log.Error("error writing logs to file")
	}
}
